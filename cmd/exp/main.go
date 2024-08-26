package main

import (
	"errors"
	"fmt"
)

var ErrorNotFound = errors.New("ERR NOT FOUND")

func main() {
	err := B()
	if errors.Is(err, ErrorNotFound) {
		fmt.Println("WE GOT THE ERROR CAPTAIN!")
	}
}

func A() error {
	return ErrorNotFound
}

func B() error {
	err := A()
	if err != nil {
		return fmt.Errorf("ERR IN FUNC B: %w", err)
	}
	return nil
}
