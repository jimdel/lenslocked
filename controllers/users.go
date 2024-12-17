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
	u.Templates.New.Execute(w, nil)
}

func (u User) Create(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "Creating user...", r.FormValue("email"), r.FormValue("password"))
	fmt.Fprint(w, r.FormValue("email"))
	fmt.Fprint(w, r.FormValue("password"))
}
