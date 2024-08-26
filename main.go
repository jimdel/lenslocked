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

func executeTemplate(w http.ResponseWriter, tplPath string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tpl, err := template.ParseFiles(tplPath)
	if err != nil {
		log.Printf("Parsing template: %v", err)
		http.Error(w, "There was an error parsing the template", http.StatusInternalServerError)
		return
	}

	err = tpl.Execute(w, data)
	if err != nil {
		log.Printf("Executing template: %v", err)
		http.Error(w, "There was an error executing the template", http.StatusInternalServerError)
		return
	}
}

func defaultHandlerFunc(w http.ResponseWriter, _ *http.Request) {
	tplPath := filepath.Join("templates", "home.gohtml")
	executeTemplate(w, tplPath, nil)
}

func contactHandlerFunc(w http.ResponseWriter, _ *http.Request) {
	tplPath := filepath.Join("templates", "contact.gohtml")
	executeTemplate(w, tplPath, nil)
}

func faqHandlerFunc(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "faq.gohtml")
	contents, err := os.ReadFile("faq.txt")

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
	executeTemplate(w, tplPath, questions)
}

func urlParamHandler(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "rank")
	fmt.Fprint(w, "<h1>"+param+"</h1>")
}

func peopleHandlerFunc(w http.ResponseWriter, r *http.Request) {
	type ContactInfo struct {
		Email       string
		Address     string
		ShowAddress bool
	}
	type Person struct {
		Name    string
		Contact ContactInfo
	}
	var people []Person = []Person{
		{
			Name: "james",
			Contact: ContactInfo{
				Email:       "jimdel@github.com",
				Address:     "300 Herb Street, L.I.N.Y",
				ShowAddress: false,
			},
		},
		{
			Name: "tim",
			Contact: ContactInfo{
				Email:       "timdel@github.com",
				Address:     "301 Herb Street, L.I.N.Y",
				ShowAddress: true,
			},
		},
	}
	var aliens map[string]string = map[string]string{
		"1": "Romulus",
		"2": "Remulus",
		"3": "Reemulus",
		"4": "Carl",
	}

	tplPath := filepath.Join("templates", "people.gohtml")
	executeTemplate(w, tplPath, struct {
		People []Person
		Aliens map[string]string
	}{People: people, Aliens: aliens})
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
	r.Get("/contact", contactHandlerFunc)
	r.Get("/people", peopleHandlerFunc)
	r.Get("/", defaultHandlerFunc)

	err := http.ListenAndServe(PORT, r)
	if err != nil {
		fmt.Printf("<< SERVER ERROR >>")
		panic(err)
	}
}
