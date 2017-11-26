package aql

import (
	"encoding/json"
	"github.com/alecthomas/participle"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)


func getExpectedResult(scriptPath string) (*JobScript, error) {
	jsonPath := strings.Replace(scriptPath, ".txt", ".json", 1)
	var bb JobScript
	b, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &bb)
	return &bb, err
}

func saveExpectedResult(scriptPath string, b JobScript) error {
	bb, err := json.Marshal(b)
	if err != nil {
		return err
	}
	jsonPath := strings.Replace(scriptPath, ".txt", ".json", 1)
	err = ioutil.WriteFile(jsonPath, bb, 666)
	return err
}

func TestParseOptions(t *testing.T){
	Convey("Given some options as a string", t, func(){
		Convey("It should parse them correctly if they are valid", func(){
			s := "key1:1.2,key2:1,key3:\"asdf\""
			opts, err := StrToOpts(s)
			v1 := 1.2
			v2 := float64(1)
			v3 := "asdf"
			So(err, ShouldBeNil)
			So(opts, ShouldHaveLength, 3)
			So(opts[0].Key, ShouldEqual, "key1")
			So(opts[0].Value, ShouldResemble, &OptionValue{Number: &v1})
			So(opts[1].Value, ShouldResemble, &OptionValue{Number: &v2})
			So(opts[2].Value, ShouldResemble, &OptionValue{Str: &v3})
		})
		Convey("It should return an error if they are invalid", func(){
			s := "key1:null,key2:1,key3\"asdf\""
			_, err := StrToOpts(s)

			So(err, ShouldNotBeNil)
		})

	})

}

func TestParseExcelRange(t *testing.T){
	Convey("Given a valid Excel range", t, func(){
		s := "A2:*1"
		x1, x2, y1, y2, err := ParseExcelRange(s)
		So(err, ShouldBeNil)
		So(x1, ShouldEqual, 1)
		So(x2, ShouldBeNil)
		So(y1, ShouldEqual, 2)
		So(*y2, ShouldEqual, 1)
	})
	Convey("Given an invalid Excel range", t, func(){
		s := "A*:B]"
		_,_,_,_,err := ParseExcelRange(s)
		So(err, ShouldNotBeNil)
	})
}

func TestFindOption(t *testing.T){
	Convey("Given some options", t, func(){
		f := 1.0
		f2 := 2.0
		opts := []Option{
			Option{
				Key: "O1",
				Value: &OptionValue{
					Number: &f,
				},
			},
			Option{
				Key: "O2",
				Value: &OptionValue{
					Number : &f2,
				},
			},
		}
		opt, ok := FindOption(opts, "o2")
		So(ok, ShouldBeTrue)
		So(opt.Key, ShouldEqual, "O2")
		So(*opt.Value.Number, ShouldEqual, f2)
		_, ok = FindOption(opts, "o3")
		So(ok, ShouldBeFalse)
	})
}

func TestFindOverridableOption(t *testing.T){
	Convey("Given some options", t, func(){
		f := 1.0
		f2 := 2.0
		f3 := 3.0
		f4 := 4.0
		opts := []Option{
			Option{
				Key: "asdf_O1",
				Value: &OptionValue{
					Number: &f,
				},
			},
			Option{
				Key: "O2",
				Value: &OptionValue{
					Number : &f2,
				},
			},
		}
		opts2 := []Option{
			Option{
				Key: "O1",
				Value: &OptionValue{
					Number: &f3,
				},
			},
			Option{
				Key: "O3",
				Value: &OptionValue{
					Number: &f4,
				},
			},
		}
		opt, ok := FindOverridableOption("O3", "", opts, opts2)
		So(ok, ShouldBeTrue)
		So(*opt.Value.Number, ShouldEqual, f4)
		opt, ok = FindOverridableOption("O1", "ASDF", opts, opts2)
		So(ok, ShouldBeTrue)
		So(*opt.Value.Number, ShouldEqual, f)

	})
}

func TestTruthy(t *testing.T){
	Convey("Given some options that may or not be truthy", t, func(){
		v1 := float64(1)
		s1 := "true"
		s2 := "false"
		v2 := float64(0)
		o1 := Option{
			Value: &OptionValue{
				Number: &v1,
			},
		}
		o2 := Option{
			Value: &OptionValue{
				Number: &v2,
			},
		}
		o3 := Option{
			Value: &OptionValue{
				Str: &s1,
			},
		}
		o4 := Option{
			Value: &OptionValue{
				Str: &s2,
			},
		}
		Convey("It should return whether the value is truthy correctly", func(){
			So(o1.Truthy(), ShouldBeTrue)
			So(o2.Truthy(), ShouldBeFalse)
			So(o3.Truthy(), ShouldBeTrue)
			So(o4.Truthy(), ShouldBeFalse)
		})

	})
}

func TestString(t *testing.T){
	Convey("Given some options", t, func(){
		Convey("It should successfully return string value of a String option", func(){
			s := "string"
			o1 := Option{
				Value: &OptionValue{
					Str: &s,
				},
			}
			ss, ok := o1.String()
			So(ok, ShouldBeTrue)
			So(ss, ShouldEqual, s)
		})
		Convey("It should return false when passed Number option", func(){
			v := float64(1)
			o1 := Option{
				Value: &OptionValue{
					Number: &v,
				},
			}
			_, ok := o1.String()
			So(ok, ShouldBeFalse)

		})
	})
}

func TestCompareOutput(t *testing.T) {
	Convey("Given the scripts in the testing folder", t, func() {
		Convey("It should parse each without error and the output should be as expected", func() {
			s, err := filepath.Glob("./testing/*.txt")
			if err != nil {
				panic(err)
			}
			for _, ss := range s {
				bs, err := ParseFile(ss)
				//sss, err := json.Marshal(bs)
				//fmt.Println(string(sss))
				So(err, ShouldBeNil)
				js, err := getExpectedResult(ss)
				So(err, ShouldBeNil)
				So(bs, ShouldResemble, js)
			}
		})
	})
}

func TestQuery(t *testing.T) {
	parser, err := participle.Build(&Query{}, &definition{})
	if err != nil {
		panic(err)
	}
	Convey("It should parse query blocks successfully", t, func() {
		//1
		s1 := `QUERY 'name' FROM CONNECTION source (
			query_source()
		) INTO CONNECTION destination, GLOBAL
		AFTER dependency
		`
		b := &Query{}
		err = parser.ParseString(s1, b)
		So(err, ShouldBeNil)
		So(b.Name, ShouldEqual, "name")
		So(strings.TrimSpace(b.Content), ShouldEqual, "query_source()")
		So(b.Sources, ShouldHaveLength, 1)
		s := "source"
		So(b.Sources[0].Database, ShouldResemble, &s)
		So(b.Destinations[1].Global, ShouldBeTrue)
		So(*b.Destinations[0].Database, ShouldEqual, "destination")
		So(b.Dependencies, ShouldResemble, []string{"dependency"})

		//2
		s1 = `QUERY 'name' EXTERN 'sourcee'
		FROM GLOBAL, SCRIPT 'asdf.py' (
			thing''
		) INTO GLOBAL
		`
		b = &Query{}
		err = parser.ParseString(s1, b)
		So(err, ShouldBeNil)
		So(b.Name, ShouldEqual, "name")
		ss := "sourcee"
		So(b.Extern, ShouldResemble, &ss)
		So(strings.TrimSpace(b.Content), ShouldEqual, "thing''")
		So(b.Sources, ShouldHaveLength, 2)
		So(b.Sources[0].Global, ShouldBeTrue)
		ss = "asdf.py"
		So(b.Sources[1].Script, ShouldResemble, &ss)
		So(b.Destinations[0].Global, ShouldBeTrue)

		//3
		s1 = `QUERY 'name' EXTERN 'sourcee'
		FROM GLOBAL AS 'source' (
			thing''
		) INTO CONNECTION destination
		WITH (opt1 = 'val', opt2 = 1234)
		`
		b = &Query{}
		err = parser.ParseString(s1, b)
		So(err, ShouldBeNil)
		So(b.Name, ShouldEqual, "name")
		ss = "sourcee"
		So(b.Extern, ShouldResemble, &ss)
		So(strings.TrimSpace(b.Content), ShouldEqual, "thing''")
		So(b.Sources, ShouldHaveLength, 1)
		So(b.Sources[0].Global, ShouldBeTrue)
		So(*b.Sources[0].Alias, ShouldEqual, "source")
		So(*b.Destinations[0].Database, ShouldEqual, "destination")
		So(b.Options, ShouldHaveLength, 2)
		So(b.Options[0].Key, ShouldEqual, "opt1")
		f := 1234.0
		So(b.Options[1].Value.Number, ShouldResemble, &f)

	})
}

func TestInclude(t *testing.T) {
	parser, err := participle.Build(&Include{}, &definition{})
	if err != nil {
		panic(err)
	}
	Convey("It should parse Include blocks successfully", t, func() {
		s1 := `INCLUDE 'source.txt'`
		b := &Include{}
		err = parser.ParseString(s1, b)
		So(b.Source, ShouldEqual, "source.txt")

	})
}

func TestScript(t *testing.T) {
	parser, err := participle.Build(&Script{}, &definition{})
	if err != nil {
		panic(err)
	}
	Convey("It should parse script blocks successfully", t, func() {
		//1
		s1 := `SCRIPT 'name' FROM CONNECTION source (
			query_source()
		) INTO CONNECTION destination
		`
		b := &Script{}
		err = parser.ParseString(s1, b)
		So(err, ShouldBeNil)
		So(b.Name, ShouldEqual, "name")
		So(strings.TrimSpace(b.Content), ShouldEqual, "query_source()")
		So(b.Sources, ShouldHaveLength, 1)
		s := "source"
		So(b.Sources[0].Database, ShouldResemble, &s)
		So(*b.Destinations[0].Database, ShouldEqual, "destination")
	})
}

func TestTest(t *testing.T) {
	parser, err := participle.Build(&Test{}, &definition{})
	if err != nil {
		panic(err)
	}
	Convey("It should parse test blocks successfully", t, func() {
		//1
		s1 := `TEST SCRIPT 'name' FROM CONNECTION source (
			query_source()
		)
		`
		b := &Test{}
		err = parser.ParseString(s1, b)
		So(err, ShouldBeNil)
		So(b.Name, ShouldEqual, "name")
		So(strings.TrimSpace(b.Content), ShouldEqual, "query_source()")
		So(b.Sources, ShouldHaveLength, 1)
		s := "source"
		So(b.Sources[0].Database, ShouldResemble, &s)
		So(b.Script, ShouldBeTrue)
		So(b.Query, ShouldBeFalse)
	})
}

func TestGlobal(t *testing.T) {
	parser, err := participle.Build(&Global{}, &definition{})
	if err != nil {
		panic(err)
	}
	Convey("It should parse global blocks successfully", t, func() {
		//1
		s1 := `GLOBAL 'name' (
			query_source()
		)
		`
		b := &Global{}
		err = parser.ParseString(s1, b)
		So(err, ShouldBeNil)
		So(b.Name, ShouldEqual, "name")
		So(strings.TrimSpace(b.Content), ShouldEqual, "query_source()")
	})
}

func TestDescription(t *testing.T) {
	parser, err := participle.Build(&Description{}, &definition{})
	if err != nil {
		panic(err)
	}
	Convey("It should parse description blocks successfully", t, func() {
		//1
		s1 := `DESCRIPTION 'This is a
		description'
		`
		b := &Description{}
		err = parser.ParseString(s1, b)
		So(err, ShouldBeNil)
		So(b.Content, ShouldEqual, `This is a
		description`)
	})
}

func TestConnection(t *testing.T) {
	parser, err := participle.Build(&UnparsedConnection{}, &definition{})
	if err != nil {
		panic(err)
	}
	Convey("It should process connection blocks successfully", t, func() {
		//1
		s1 := `CONNECTION 'test' (
			Driver = 'MSSQL'
			ConnectionString = 'asdf'
		)
		`
		b := &UnparsedConnection{}
		err = parser.ParseString(s1, b)
		So(err, ShouldBeNil)
		So(b.Name, ShouldEqual, "test")
		So(b.Content, ShouldEqual, `
			Driver = 'MSSQL'
			ConnectionString = 'asdf'
		`)
	})
	Convey("It should parse connection block contents successfully", t, func() {
		b := UnparsedConnection{
			Name: "Test",
			Content: `
			Driver = 'MSSQL',
			ConnectionString = 'asdf'
		`,
		}
		bb, err := parseConnections([]UnparsedConnection{b})
		So(err, ShouldBeNil)
		So(bb, ShouldHaveLength, 1)
		So(bb[0].Driver, ShouldEqual, "MSSQL")
		So(bb[0].Name, ShouldEqual, "Test")
		So(bb[0].ConnectionString, ShouldEqual, "asdf")
	})
}

func TestResolveIncludes(t *testing.T) {
	Convey("Given a script with an INCLUDE statement", t, func() {
		q := `INCLUDE 'testing/2.txt'`
		b, err := ParseString(q)
		So(err, ShouldBeNil)
		Convey("It should correctly resolve the included resources", func() {
			err = b.ResolveExternalContent()
			So(err, ShouldBeNil)
			So(b.Queries, ShouldHaveLength, 2)
			So(b.Description, ShouldNotBeNil)
			So(b.Queries[0].Name, ShouldEqual, "b")
			So(b.Queries[0].Content, ShouldEqual, "TEST EXTERNAL CONTENT")
			So(b.Queries[1].Name, ShouldEqual, "q1")
			So(*b.Queries[1].Destinations[0].Database, ShouldEqual, "d1")
		})
	})
}

func TestParameterEvaluation(t *testing.T) {
	Convey("Given a script with parameters", t, func() {
		q := `QUERY 'a' FROM GLOBAL (
			SELECT * FROM {{ .Table }}
		) INTO GLOBAL
		WITH (Table = 'Something')`
		b, err := ParseString(q)
		So(err, ShouldBeNil)
		Convey("It should correctly evaluate the content", func() {
			err = b.EvaluateParametrizedContent(nil)
			So(err, ShouldBeNil)
			So(b.Queries, ShouldHaveLength, 1)
			So(b.Queries[0].Options, ShouldHaveLength, 1)
			So(b.Queries[0].Content, ShouldEqual, `
			SELECT * FROM Something
		`)
		})
	})
}