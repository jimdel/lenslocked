package controllers

import (
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/jimdel/lenslocked/data"
)

type PageMeta struct {
	Title string
}

func StaticHandler(tpl Template, title string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pm := PageMeta{title}
		tpl.Execute(w, r, pm)
	}
}

func FAQStaticHandler(tpl Template, title string) http.HandlerFunc {
	contents, err := fs.ReadFile(data.FS, "faq.txt")
	rawFaq := strings.Split(string(contents), "\n")

	type Question struct {
		Q string
		A string
	}
	type Questions []Question

	var questions Questions

	for idx := 0; idx <= len(rawFaq)-2; {
		question := Question{
			Q: rawFaq[idx],
			A: rawFaq[idx+1],
		}
		questions = append(questions, question)
		idx += 2
	}
	return func(w http.ResponseWriter, r *http.Request) {

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Something terrible occured!", http.StatusInternalServerError)
			return
		}
		tpl.Execute(w, r, struct {
			Questions Questions
			Title     string
		}{questions, title})
	}
}
