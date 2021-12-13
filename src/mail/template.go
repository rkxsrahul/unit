package mail

import (
	"bytes"
	"log"
	"text/template"
)

// EmailTemplate is a method for parsing email template
// convert it into string and set variables values
func EmailTemplate(tmplPath string, data map[string]interface{}) string {
	// parsing template file
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Println(err)
		return ""
	}
	// creating new buffer as io writer
	buf := new(bytes.Buffer)
	// pasing above template with data and result data in buffer
	err = tmpl.Execute(buf, data)
	if err != nil {
		log.Println(err)
		return ""
	}
	// return buffer in string
	return buf.String()
}
