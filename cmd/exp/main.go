package main

import (
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jimdel/lenslocked/models"
)

func main() {
	config := models.DefaultPostgresConfig()
	db, err := models.Open(config)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("Error pinging db: %v\n", err)
	}
	defer db.Close()
}
