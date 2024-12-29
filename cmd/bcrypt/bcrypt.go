package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func hash(pw string) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(hashedPassword))
}

func compare(pw, hash string) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
	if err != nil {
		fmt.Printf("Incorrect password %v\n", pw)
	} else {
		fmt.Println("Password matches!")
	}
}

func main() {

	switch os.Args[1] {
	case "hash":
		hash(os.Args[2])
	case "compare":
		compare(os.Args[2], os.Args[3])
	default:
		fmt.Printf("Invalid command %s\n", os.Args[1])
	}
}
