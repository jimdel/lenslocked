package controllers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/jimdel/lenslocked/context"
	"github.com/jimdel/lenslocked/models"
)

type Users struct {
	Templates struct {
		New            Template
		SignIn         Template
		CurrentUser    Template
		ForgotPassword Template
		CheckYourEmail Template
	}
	UserService          *models.UserService
	SessionService       *models.SessionService
	PasswordResetService *models.PasswordResetService
	EmailService         *models.EmailService
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

func (u Users) ProcessSignOut(w http.ResponseWriter, r *http.Request) {
	token, err := ReadCookie(r, CookieSession)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	err = u.SessionService.Delete(token)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	DeleteCookie(w, CookieSession)
	http.Redirect(w, r, "/signin", http.StatusFound)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	u.Templates.CurrentUser.Execute(w, r, user)
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

func (u Users) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.ForgotPassword.Execute(w, r, data)
}

func (u Users) ProcessForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	pwReset, err := u.PasswordResetService.Create(data.Email)
	if err != nil {
		// TODO: implement better err handling (ex. no user with email)
		fmt.Println(err)
		http.Error(w, "Something went wrong!", http.StatusInternalServerError)
		return
	}
	vals := url.Values{
		"token": {pwReset.Token},
	}
	resetUrl := "https://localhost:42069/reset-pw?" + vals.Encode()
	err = u.EmailService.ForgotPassword(data.Email, resetUrl)
	if err != nil {
		// TODO: implement better err handling (ex. no user with email)
		fmt.Println(err)
		http.Error(w, "Something went wrong!", http.StatusInternalServerError)
		return
	}

	// Don't render the token - this should be sent via email only for id verification
	u.Templates.CheckYourEmail.Execute(w, r, data)
}
