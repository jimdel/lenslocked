package views

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/jimdel/lenslocked/context"
	"github.com/jimdel/lenslocked/models"
)

type Template struct {
	htmlTpl *template.Template
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	tpl := template.New(patterns[0])

	//Register all template functions
	tpl = tpl.Funcs(
		template.FuncMap{
			"CSRF": func() (template.HTML, error) {
				return `<!-- placeholder -->`, fmt.Errorf("CSRF not implemented")
			},
			"currentUser": func() (template.HTML, error) {
				return `<!-- placeholder -->`, fmt.Errorf("currentUser not implemented")
			},
		},
	)

	tpl, err := tpl.ParseFS(fs, patterns...)

	if err != nil {
		return Template{}, fmt.Errorf("ERR PARSE FS TEMPLATE: %v", err)
	}

	return Template{tpl}, nil
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}) {

	tpl, err := t.htmlTpl.Clone()
	if err != nil {
		fmt.Printf("Error cloning template: %v", err)
	}

	// Implement registered template functions
	tpl = tpl.Funcs(
		template.FuncMap{
			"CSRF": func() template.HTML {
				return csrf.TemplateField(r)
			},
			"currentUser": func() *models.User {
				return context.User(r.Context())
			},
		},
	)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Write data to memory
	var buf bytes.Buffer
	// Check for template parsing errors
	err = tpl.Execute(&buf, data)

	if err != nil {
		fmt.Printf("Error executing template: %v", err)
		http.Error(w, "There was an error executing the template", http.StatusInternalServerError)
		return
	}

	//Copy data from the buffer to the response writer
	io.Copy(w, &buf)

	// INFO: may be troublesome with LARGE templates
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}
