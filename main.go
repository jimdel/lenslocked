package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const PORT = ":42069"

func defaultHandlerFunc(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tplPath := filepath.Join("templates", "home.gohtml")
	tpl, err := template.ParseFiles(tplPath)
	if err != nil {
		log.Printf("Parsing template: %v", err)
		http.Error(w, "There was an error parsing the template", http.StatusInternalServerError)
		return
	}
	err = tpl.Execute(w, struct{ Hi string }{Hi: "Hey"})
	if err != nil {
		log.Printf("Executing template: %v", err)
		http.Error(w, "There was an error executing the template", http.StatusInternalServerError)
		return
	}
}

func faqHandlerFunc(w http.ResponseWriter, r *http.Request) {
	contents, err := os.ReadFile("faq.txt")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something terrible occured!", http.StatusInternalServerError)
	} else {
		rawFaq := string(contents)
		var faqHtml string
		for _, line := range strings.Split(rawFaq, "\n") {
			faqHtml += "<p>" + line + "</p>"
		}
		w.Header().Set("Content-Type", "text/html;charset=utf-8")
		fmt.Fprint(w, faqHtml)
	}
}

func urlParamHandler(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "rank")
	fmt.Fprint(w, "<h1>"+param+"</h1>")
}

func main() {

	fmt.Printf("Server listening on port %v\n", PORT)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 - Not Found!", http.StatusNotFound)
	})

	// Apply middleware to single route
	// r.Route("/faq", func(r chi.Router) {
	// 	r.Use(middleware.Logger)
	// 	r.Get("/faq", faqHandlerFunc)
	// })

	r.Get("/jimdel/{rank}", urlParamHandler)
	r.Get("/faq", faqHandlerFunc)
	r.Get("/", defaultHandlerFunc)

	err := http.ListenAndServe(PORT, r)
	if err != nil {
		fmt.Printf("<< SERVER ERROR >>")
		panic(err)
	}
}
