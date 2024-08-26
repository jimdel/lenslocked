package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jimdel/lenslocked/controllers"
	"github.com/jimdel/lenslocked/views"
)

const PORT = ":42069"

func main() {

	fmt.Printf("Server listening on port %v\n", PORT)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 - Not Found!", http.StatusNotFound)
	})

	// Home
	tpl, err := views.Parse("home.gohtml")
	if err != nil {
		panic(err)
	}
	r.Get("/", controllers.StaticHandler(tpl))

	// Contact
	tpl, err = views.Parse("contact.gohtml")
	if err != nil {
		panic(err)
	}
	r.Get("/contact", controllers.StaticHandler(tpl))

	// FAQ
	tpl, err = views.Parse("faq.gohtml")
	if err != nil {
		panic(err)
	}
	r.Get("/faq", controllers.FAQStaticHandler(tpl))

	err = http.ListenAndServe(PORT, r)
	if err != nil {
		fmt.Printf("<< SERVER ERROR >>")
		panic(err)
	}
}
