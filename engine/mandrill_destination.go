package engine

import (
	"fmt"
	"github.com/keighl/mandrill"
	"regexp"
	"strings"
	"sync/atomic"
	"time"
)

var emailRecipientRegexp = regexp.MustCompile(`^\s*([\w\s]+)\s*<\s*(\w[-._\w]*\w@\w[-._\w]*\w\.\w{2,3})\s*>\s*$`)

type MandrillPrincipal struct {
	Name  string
	Email string
}

type MandrillDestination struct {
	Name       string
	APIKey     string `aql:"API_KEY"`
	Sender     *MandrillPrincipal
	Recipients []MandrillPrincipal
	SplitByRow bool   `aql:"SPLIT, optional"`
	Template   string `aql:"TEMPLATE"`
	Subject    string `aql:"SUBJECT, optional"`
	client     *mandrill.Client
	emailsSent int64
	cols       []string
}

func (d *MandrillDestination) Ping() error {
	d.client = mandrill.ClientWithKey(d.APIKey)
	_, err := d.client.Ping()
	return err
}

func ParseEmailRecipients(s string) ([]MandrillPrincipal, error) {
	var ret []MandrillPrincipal
	recipientStr := strings.Split(s, ",")

	for _, recip := range recipientStr {
		matches := emailRecipientRegexp.FindStringSubmatch(recip)
		if len(matches) != 3 {
			return nil, fmt.Errorf("invalid syntax or email for recipient %s. Expecting NAME <EMAIL>", recip)
		}
		ret = append(ret, MandrillPrincipal{Name: strings.TrimSpace(matches[1]), Email: matches[2]})
	}
	return ret, nil
}

func (d *MandrillDestination) Open(s Stream, l Logger, st Stopper) {
	c := s.Chan(d.Name)
	var (
		rows         []map[string]interface{}
		cols         []string
		firstMessage = true
	)

	for msg := range c {
		d.log(l, Info, "Mandrill destination opened")
		if st.Stopped() {
			d.log(l, Warning, "Mandrill destination aborted")
			return
		}
		if firstMessage {
			cols = s.Columns()
			firstMessage = false
		}

		if d.SplitByRow {
			m := d.prepareMsg()
			content := d.prepareContent(cols, msg.Data)
			_, err := d.client.MessagesSendTemplate(m, d.Template, content)
			atomic.AddInt64(&d.emailsSent, 1)
			if err != nil {
				d.fatalerr(err, l)
				return
			}
			d.log(l, Trace, "sent email to recipients with content %v", content)
		} else {
			rows = append(rows, d.prepareContent(cols, msg.Data))
		}
	}

	if !d.SplitByRow {
		m := d.prepareMsg()
		_, err := d.client.MessagesSendTemplate(m, d.Template, rows)
		atomic.AddInt64(&d.emailsSent, 1)
		if err != nil {
			d.fatalerr(err, l)
			return
		}
		d.log(l, Info, "sent email to recipients containing all received rows")
	}
	d.log(l, Info, "Sent a total of %v emails - finished", d.emailsSent)
}

func (d *MandrillDestination) prepareContent(cols []string, row []interface{}) map[string]interface{} {
	ret := make(map[string]interface{})
	for i, col := range cols {
		ret[col] = row[i]
	}
	return ret
}

func (d *MandrillDestination) prepareMsg() *mandrill.Message {
	var m mandrill.Message
	if d.Subject != "" {
		//Could be set as part of the template
		m.Subject = d.Subject
	}
	if d.Sender != nil {
		//Could be set as part of the template
		m.FromName = d.Sender.Name
		m.FromEmail = d.Sender.Email
	}
	for _, recipient := range d.Recipients {
		m.AddRecipient(recipient.Email, recipient.Name, "to")
	}
	return &m
}

func (d *MandrillDestination) log(l Logger, level LogLevel, msgFormat string, args ...interface{}) {
	l.Chan() <- Event{
		Source:  d.Name,
		Level:   level,
		Time:    time.Now(),
		Message: fmt.Sprintf(msgFormat, args...),
	}
}

func (d *MandrillDestination) fatalerr(err error, l Logger) {
	l.Chan() <- Event{
		Level:   Error,
		Source:  d.Name,
		Time:    time.Now(),
		Message: err.Error(),
	}
}
