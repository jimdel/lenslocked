package controllers

import (
	"fmt"
	"net/http"
)

type User struct {
	Templates struct {
		New Template
	}
}

func (u User) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, data)
}

func (u User) Create(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "Creating user...", r.FormValue("email"), r.FormValue("password"))
	fmt.Fprint(w, r.FormValue("email"))
	fmt.Fprint(w, r.FormValue("password"))
}
