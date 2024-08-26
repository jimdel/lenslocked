package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jimdel/lenslocked/controllers"
	"github.com/jimdel/lenslocked/templates"
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

	r.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "home.gohtml"))))
	r.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "contact.gohtml"))))
	r.Get("/faq", controllers.FAQStaticHandler(views.Must(views.ParseFS(templates.FS, "faq.gohtml"))))
	r.Get("/static", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "static.gohtml"))))

	err := http.ListenAndServe(PORT, r)
	if err != nil {
		fmt.Printf("<< SERVER ERROR >>")
		panic(err)
	}
}
