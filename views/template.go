package views

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
)

type Template struct {
	htmlTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := t.htmlTpl.Execute(w, data)

	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "There was an error executing the template", http.StatusInternalServerError)
		return
	}
}

func Parse(templateName string) (Template, error) {
	tplPath := filepath.Join("templates", templateName)
	tpl, err := template.ParseFiles(tplPath)
	if err != nil {
		return Template{}, fmt.Errorf("ERR PARSING TEMPLATE: %v", err)
	}
	tmpl := Template{tpl}
	return tmpl, nil
}

func ParseFS(fs fs.FS, pattern string) (Template, error) {
	tpl, err := template.ParseFS(fs, pattern)
	if err != nil {
		return Template{}, fmt.Errorf("ERR PARSE FS TEMPLATE: %v", err)
	}
	tmpl := Template{tpl}
	return tmpl, nil
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}
