package main

import (
	"errors"
	"fmt"
)

func main() {
	err := CreateOrg()
	fmt.Println(err)
}

func Connect() error {
	return errors.New("ERR CONNECTION FAILED")
}

func CreateUser() error {
	err := Connect()
	if err != nil {
		return fmt.Errorf("ERR CREATE USER -> %w", err)
	}
	return nil
}

func CreateOrg() error {
	err := CreateUser()
	if err != nil {
		return fmt.Errorf("ERR CREATE ORG -> %w", err)
	}
	return nil
}
