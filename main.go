package main

import (
	"fmt"
	"net/http"
)

const PORT = ":42069"

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Happy, healthy, wealthy Jimbroski!</h1>")
}

func main() {

	fmt.Printf("Server listening on port %v\n", PORT)
	http.HandleFunc("/", handlerFunc)

	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		fmt.Printf("<< SERVER ERROR >>")
		panic(err)
	}
}
