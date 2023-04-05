package main

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

func parseTemplate(jobs []Job) string {
	tmplBytes, err := os.ReadFile("prometheus-scrape.tmpl")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	t := template.Must(template.New("tmpl").Parse(string(tmplBytes)))

	buff := bytes.NewBuffer([]byte{})
	t.Execute(buff, map[string][]Job{
		"jobs": jobs,
	})

	return buff.String()
}
