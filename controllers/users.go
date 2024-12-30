package controllers

import (
	"fmt"
	"net/http"

	"github.com/jimdel/lenslocked/models"
)

type Users struct {
	Templates struct {
		New    Template
		SignIn Template
	}
	UserService *models.UserService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, data)
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}

	data.Email = r.FormValue("email")
	u.Templates.SignIn.Execute(w, data)
}

func (u Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	unauthenticatedUser := models.UnauthenticatedUser{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	user, err := u.UserService.Authenticate(unauthenticatedUser)
	if err != nil {
		http.Error(w, "Sign in failed...", http.StatusInternalServerError)
	} else {
		fmt.Fprintf(w, "User authenticated: %v+", user)
	}
	// u.Templates.SignIn.Execute(w, &user)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {

	unauthenticatedUser := models.UnauthenticatedUser{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	user, err := u.UserService.Create(unauthenticatedUser)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong...", http.StatusInternalServerError)
	}

	fmt.Println("User created")
	fmt.Fprint(w, user)

}
