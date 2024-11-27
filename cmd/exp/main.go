package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type DBConf struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func (c DBConf) generateDbConnStr() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

func main() {
	var conf = DBConf{
		Host:     "localhost",
		Port:     "8081",
		User:     "baloo",
		Password: "junglebook",
		DBName:   "lenslocked",
		SSLMode:  "disable",
	}

	db, err := sql.Open("pgx", conf.generateDbConnStr())
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("CONNECTED SUCCESSFULLY")
	defer db.Close()

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT UNIQUE,
		createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
	);

	CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		user_id INT,
		amount INT,
		description TEXT,
		createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`)
	username := "Jim Del"
	email := "jimdel@github.com"
	row := db.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, name;", username, email)
	var user struct {
		id   int
		name string
	}

	err = row.Scan(&user)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%d", user)
}
