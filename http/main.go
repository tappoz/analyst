package main

import (
	"fmt"

	"encoding/json"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/michaelbironneau/analyst"
	"golang.org/x/net/websocket"
)

const (
	MsgLog           = "LOG"
	MsgRunScript     = "RUN"
	MsgResult        = "RESULT"
	MsgCompileScript = "COMPILE"
	MsgOutput        = "OUTPUT"
)

type RunMessagePayload struct {
	Script string `json:"script"`
}

type RunResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

func receiveMessages(ws *websocket.Conn, c echo.Context) {
	fmt.Println("Starting to receive")
	for {
		var b []byte
		err := websocket.Message.Receive(ws, &b)
		if err != nil {
			c.Logger().Error(err)
			break
		}
		var m Message
		err = json.Unmarshal(b, &m)
		if err != nil {
			c.Logger().Error(err)
		}
		switch m.Type {
		case MsgRunScript:
			var payload RunMessagePayload
			if err := json.Unmarshal(m.Data, &payload); err != nil {
				c.Logger().Error(err)
				continue
			}
			opts := outputHooks(ws)
			err := analyst.ExecuteString(payload.Script, &opts)
			var response RunResponse
			response.Success = err == nil
			if err != nil {
				response.Error = err.Error()
			}

			send(ws, MsgRunScript, response)
		case MsgCompileScript:
			var payload RunMessagePayload
			if err := json.Unmarshal(m.Data, &payload); err != nil {
				c.Logger().Error(err)
				continue
			}
			err := analyst.ValidateString(payload.Script, &analyst.RuntimeOptions{})
			var resp RunResponse
			if err != nil {
				resp.Success = false
				resp.Error = err.Error()
			} else {
				resp.Success = true
			}
			send(ws, MsgCompileScript, resp)
		default:
			c.Logger().Error(fmt.Sprintf("unknown message type %s", m.Type))
		}
	}
}

//send message, ignoring errors
func send(ws *websocket.Conn, messageType string, payload interface{}) {
	var m Message
	m.Type = messageType
	b, _ := json.Marshal(payload)
	m.Data = json.RawMessage(b)
	b, _ = json.Marshal(m)
	err := websocket.Message.Send(ws, string(b))
	if err != nil {
		fmt.Println(err)
	}
}

func hello(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		for {

			// Write
			err := websocket.Message.Send(ws, "{\"type\": \"b\"}")
			if err != nil {
				c.Logger().Error(err)
			}

			// Read
			msg := ""
			err = websocket.Message.Receive(ws, &msg)
			if err != nil {
				c.Logger().Error(err)
			}
			fmt.Printf("%s\n", msg)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func receive(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		receiveMessages(ws, c)
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	//e.Use(middleware.Recover())
	//e.Static("/", "../public")
	e.GET("/", receive)
	e.Logger.Fatal(e.Start(":4040"))
}
