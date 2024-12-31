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
	UserService    *models.UserService
	SessionService *models.SessionService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}

	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, r, data)
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}

	data.Email = r.FormValue("email")
	u.Templates.SignIn.Execute(w, r, data)
}

func (u Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	unauthenticatedUser := models.UnauthenticatedUser{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	user, err := u.UserService.Authenticate(unauthenticatedUser)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Sign in failed...", http.StatusInternalServerError)
		return
	}

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		// TODO: implement error handling to handle sign in err
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	SetCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	token, err := ReadCookie(r, CookieSession)
	if err != nil {
		fmt.Println(err)
		// TODO: implement error handling to handle sign in err
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	user, err := u.SessionService.User(token)
	if err != nil {
		fmt.Println(err)
		// TODO: implement error handling to handle sign in err (INVALID TOKEN)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	fmt.Fprintf(w, "Current user: %v+", user)
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
		return
	}

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		// TODO: implement error handling to handle sign in err
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	SetCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}
