package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const PORT = ":42069"

func defaultHandlerFunc(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Happy, healthy, wealthy Jimbroski!</h1>")
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

func main() {

	fmt.Printf("Server listening on port %v\n", PORT)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 - Not Found!", http.StatusNotFound)
	})

	r.Get("/faq", faqHandlerFunc)
	r.Get("/", defaultHandlerFunc)

	err := http.ListenAndServe(PORT, r)
	if err != nil {
		fmt.Printf("<< SERVER ERROR >>")
		panic(err)
	}
}
