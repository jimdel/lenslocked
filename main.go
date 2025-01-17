package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/jimdel/lenslocked/controllers"
	mw "github.com/jimdel/lenslocked/middleware"
	"github.com/jimdel/lenslocked/migrations"
	"github.com/jimdel/lenslocked/models"
	"github.com/jimdel/lenslocked/templates"
	"github.com/jimdel/lenslocked/views"
	"github.com/joho/godotenv"
)

type config struct {
	PSQL models.PostgresConfig
	SMTP models.SMTPConfig
	CSRF struct {
		Key        string
		SecureMode bool
	}
	Server struct {
		Address string
	}
}

func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load(".env")
	if err != nil {
		return cfg, err
	}
	//PSQL
	cfg.PSQL.DBName = os.Getenv("DB_NAME")
	cfg.PSQL.Host = os.Getenv("DB_HOST")
	cfg.PSQL.Port = os.Getenv("DB_PORT")
	cfg.PSQL.User = os.Getenv("DB_USERNAME")
	cfg.PSQL.Password = os.Getenv("DB_PASSWORD")

	// SMTP
	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")
	SMTP_PORT := os.Getenv("SMTP_PORT")
	SMTP_PORT_STR, err := strconv.Atoi(SMTP_PORT)
	if err != nil {
		return cfg, fmt.Errorf("unable to convert smtp port to int: %w", err)
	}
	cfg.SMTP.Port = SMTP_PORT_STR

	// CSRF
	cfg.CSRF.Key = os.Getenv("APP_CSRF_KEY")
	cfg.CSRF.SecureMode = os.Getenv("APP_CSRF_MODE") == "true"

	// Server
	cfg.Server.Address = string(":" + os.Getenv("APP_PORT"))
	return cfg, nil
}

func main() {

	// Load env
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", cfg)

	psqlConfig := cfg.PSQL
	db, err := models.Open(psqlConfig)
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
	passwordResetSvc := &models.PasswordResetService{
		DB: db,
	}
	emailSvc, err := models.NewEmailService(cfg.SMTP)
	if err != nil {
		panic(err)
	}

	//END - Instantiate db svcs

	//Middleware
	umw := mw.UserMiddleware{
		SessionService: sessionSvc,
	}
	csrfMiddleware := csrf.Protect([]byte(cfg.CSRF.Key), csrf.Path("/"), csrf.Secure(cfg.CSRF.SecureMode))
	//END - Middleware

	// Controllers
	userController := controllers.Users{
		UserService:          userSvc,
		SessionService:       sessionSvc,
		EmailService:         emailSvc,
		PasswordResetService: passwordResetSvc,
	}
	userController.Templates.New = views.Must(views.ParseFS(templates.FS, "site-layout.gohtml", "signup.gohtml"))
	userController.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "site-layout.gohtml", "signin.gohtml"))
	userController.Templates.CurrentUser = views.Must(views.ParseFS(templates.FS, "site-layout.gohtml", "current-user.gohtml"))
	userController.Templates.ForgotPassword = views.Must(views.ParseFS(templates.FS, "site-layout.gohtml", "forgot-password.gohtml"))
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

	r.Get("/forgot-pw", userController.ForgotPassword)
	r.Post("/forgot-pw", userController.ProcessForgotPassword)

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

	fmt.Printf("Server listening on port %v\n", cfg.Server.Address)
	err = http.ListenAndServe(cfg.Server.Address, r)
	if err != nil {
		panic(err)
	}
}
