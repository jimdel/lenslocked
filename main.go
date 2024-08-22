package main

import (
	"fmt"
	"net/http"
)

const PORT = ":42069"

func defaultHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Happy, healthy, wealthy Jimbroski!</h1>")
}

func contactHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	fmt.Fprint(w, "<h1>Contact</h1><p>Get in touch at <a href=\"mailto:jimdel@gmail.com\">jimdel</a></p>")
}

func pathHandlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, r.URL.Path)
}

func main() {

	fmt.Printf("Server listening on port %v\n", PORT)
	// Routes
	http.HandleFunc("/contact", contactHandlerFunc)
	http.HandleFunc("/path", pathHandlerFunc)
	http.HandleFunc("/", defaultHandlerFunc)

	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		fmt.Printf("<< SERVER ERROR >>")
		panic(err)
	}
}
