package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

const PORT = ":42069"

func defaultHandlerFunc(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Happy, healthy, wealthy Jimbroski!</h1>")
}

func contactHandlerFunc(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	fmt.Fprint(w, "<h1>Contact</h1><p>Get in touch at <a href=\"mailto:jimdel@gmail.com\">jimdel</a></p>")
}

func pathHandlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, r.URL.RawPath)
}

func faqHandlerFunc(w http.ResponseWriter, r *http.Request) {
	contents, err := os.ReadFile("faq.txt")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something terrible occured!", http.StatusInternalServerError)
	}
	rawFaq := string(contents)
	var faqHtml string
	for _, line := range strings.Split(rawFaq, "\n") {
		faqHtml += "<p>" + line + "</p>"
	}

	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	fmt.Fprint(w, faqHtml)
}

// METHOD 1 - create a Router type & implement ServeHTTP method

// type Router struct{}

// func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	switch r.URL.Path {
// 	case "/path":
// 		pathHandlerFunc(w, r)
// 	case "/contact":
// 		contactHandlerFunc(w, r)
// 	case "/":
// 		defaultHandlerFunc(w, r)
// 	default:
// 		http.Error(w, "404 NOT FOUND", http.StatusNotFound)
// 	}
// }

func routerFunc(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/path":
		pathHandlerFunc(w, r)
	case "/contact":
		contactHandlerFunc(w, r)
	case "/faq":
		faqHandlerFunc(w, r)
	case "/":
		defaultHandlerFunc(w, r)
	default:
		http.Error(w, "404 NOT FOUND", http.StatusNotFound)
	}
}

func main() {

	fmt.Printf("Server listening on port %v\n", PORT)
	// Routes

	// METHOD 1 - create a Router type & implement ServeHTTP method
	// var router Router
	// err := http.ListenAndServe(PORT, router)

	// METHOD 2 - create an
	// var router http.HandlerFunc = routerFunc

	// METHOD 3
	err := http.ListenAndServe(PORT, http.HandlerFunc(routerFunc))

	if err != nil {
		fmt.Printf("<< SERVER ERROR >>")
		panic(err)
	}
}
