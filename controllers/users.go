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
	// err := r.ParseForm()
	// if err != nil {
	// 	http.Error(w, "Bad Request", http.StatusBadRequest)
	// 	return
	// }
	email := r.FormValue("email")
	pw := r.FormValue("password")
	fmt.Println(email, pw)

	fmt.Fprint(w, "TMP")
}
