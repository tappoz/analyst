package main

import (
	"fmt"
	"github.com/michaelbironneau/analyst/aql"
	"github.com/tealeg/xlsx"
	"github.com/urfave/cli"
	"gopkg.in/cheggaaa/pb.v1"
	"io/ioutil"
)

//Create creates an Excel spreadsheet based on the script
func Create(c *cli.Context) error {
	scriptPath := c.String("script")
	if len(scriptPath) == 0 {
		fmt.Println("Script path not set")
		return nil
	}
	b, err := ioutil.ReadFile(scriptPath)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	script, err := aql.Load(string(b))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	//Parse and set parameters
	params, err := parseParameters(c.String("params"))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for k, v := range params {
		if err := script.SetParameter(k, v); err != nil {
			fmt.Println(err)
			return nil
		}
	}
	//Compile script
	var task *aql.Report
	if task, err = script.ExecuteTemplates(); err != nil {
		fmt.Println(err)
		return nil
	}
	//Get connections
	connections := make(map[string]aql.Connection)
	for _, connection := range task.Connections {
		c, err := parseConn(connection)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		for k, v := range c {
			connections[k] = v
		}
	}
	//Get template
	templateFile, err := xlsx.OpenFile(task.TemplateFile)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	//Run script
	progress := make(chan int)
	errs := make(chan bool)
	bar := pb.StartNew(100)
	fmt.Println("Executing task...")
	var totalProgress int
	go func() {
		for {
			select {
			case p := <-progress:
				totalProgress += p
				bar.Add(p)
				if totalProgress >= 100 {
					return
				}
			case <-errs:
				return
			}
		}
	}()
	report, err := task.Execute(aql.DBQuery, templateFile, connections, progress)
	if err != nil {
		errs <- true
		fmt.Printf("[ERROR] %v", err)
		return nil
	}
	return report.Save(task.OutputFile)
}
