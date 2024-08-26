package controllers

import (
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/jimdel/lenslocked/data"
	"github.com/jimdel/lenslocked/views"
)

func StaticHandler(tpl views.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, nil)
	}
}

func FAQStaticHandler(tpl views.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contents, err := fs.ReadFile(data.FS, "faq.txt")

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Something terrible occured!", http.StatusInternalServerError)
			return
		}

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
		tpl.Execute(w, questions)
	}
}
