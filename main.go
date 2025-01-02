package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/jimdel/lenslocked/controllers"
	mw "github.com/jimdel/lenslocked/middleware"
	"github.com/jimdel/lenslocked/migrations"
	"github.com/jimdel/lenslocked/models"
	"github.com/jimdel/lenslocked/templates"
	"github.com/jimdel/lenslocked/views"
)

const PORT = ":42069"

func main() {

	// Setup DB connection
	config := models.DefaultPostgresConfig()
	db, err := models.Open(config)
	if err != nil {
		panic(err)
	}
	fmt.Println("DB connection successful")
	defer db.Close()
	// END - Setup DB connection

	// Setup migrations
	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}
	// END - Setup migrations

	// Instantiate db svcs
	userSvc := &models.UserService{
		DB: db,
	}
	sessionSvc := &models.SessionService{
		DB: db,
	}
	//END - Instantiate db svcs

	//Middleware
	umw := mw.UserMiddleware{
		SessionService: sessionSvc,
	}
	csrfAuthKey := "abze12fjabze12fjabze12fjabze12fj"
	// TODO: fix Secure before deployment
	csrfMiddleware := csrf.Protect([]byte(csrfAuthKey), csrf.Path("/"), csrf.Secure(false))

	// r.Use(umw.SetUser)
	//END - Middleware

	// Controllers
	userController := controllers.Users{
		UserService:    userSvc,
		SessionService: sessionSvc,
	}
	userController.Templates.New = views.Must(views.ParseFS(templates.FS, "site-layout.gohtml", "signup.gohtml"))
	userController.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "site-layout.gohtml", "signin.gohtml"))
	userController.Templates.CurrentUser = views.Must(views.ParseFS(templates.FS, "site-layout.gohtml", "current-user.gohtml"))
	//END - Controllers

	// Routes
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(csrfMiddleware)
	r.Use(umw.SetUser)

	r.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "site-layout.gohtml", "home.gohtml")), "Home"))
	r.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "site-layout.gohtml", "contact.gohtml")), "Contact"))
	r.Get("/faq", controllers.FAQStaticHandler(views.Must(views.ParseFS(templates.FS, "site-layout.gohtml", "faq.gohtml")), "FAQ"))
	r.Get("/static", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "site-layout.gohtml", "static.gohtml")), "Static"))
	r.Get("/signup", userController.New)
	r.Post("/users", userController.Create)
	r.Get("/signin", userController.SignIn)
	r.Post("/signin", userController.ProcessSignIn)
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", userController.CurrentUser)
		r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "HEY")
		})
	})
	r.Post("/signout", userController.ProcessSignOut)
	// r.Get("/users/me", controllers.Performance(userController.CurrentUser))

	// TODO: add nicer 404 page
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 - Not Found!", http.StatusNotFound)
	})
	//END - Routes

	fmt.Printf("Server listening on port %v\n", PORT)
	err = http.ListenAndServe(PORT, r)
	if err != nil {
		fmt.Printf("<< SERVER ERROR >>")
		panic(err)
	}
}
