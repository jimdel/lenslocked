package main

import (
	"fmt"
	"net/http"
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

func routerFunc(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/path":
		pathHandlerFunc(w, r)
	case "/contact":
		contactHandlerFunc(w, r)
	case "/":
		defaultHandlerFunc(w, r)
	default:
		http.Error(w, "<h1>404 NOT FOUND</h1>", http.StatusNotFound)
	}
}

func main() {

	fmt.Printf("Server listening on port %v\n", PORT)
	// Routes
	http.HandleFunc("/", routerFunc)

	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		fmt.Printf("<< SERVER ERROR >>")
		panic(err)
	}
}
