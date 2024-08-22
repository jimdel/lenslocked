package main

import (
	"fmt"
	"html/template"
	"os"
)

type User struct {
	Name  string
	Email string
}

func main() {
	fmt.Println("Experimental main")
	user := User{Name: "James DeLay", Email: "jimdel@github.com"}

	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}
	err = t.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}

}
