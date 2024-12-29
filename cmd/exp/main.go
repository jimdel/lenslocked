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

type User struct {
	ID    int
	Name  string
	Email string
}

type Order struct {
	ID          int
	UserID      int
	Amount      int
	Description string
	CreatedAt   string
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
	// username := "Jim Del"
	// email := "jimdel@github.com"
	// _ = db.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, name, email;", username, email)

	// var name, email string
	// row := db.QueryRow("SELECT name, email FROM users WHERE id = $1;", 1)
	// err = row.Scan(&name, &email)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(name, email)

	// for i := 1; i <= 10; i++ {
	// 	var order Order
	// 	order.UserID = 1
	// 	order.Amount = 100 * i
	// 	order.Description = fmt.Sprintf("Fake order %d", i)
	// 	_, err := db.Exec(`
	// 	INSERT INTO orders (user_id, amount, description) VALUES ($1, $2, $3);
	// 	`, order.UserID, order.Amount, order.Description)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	fmt.Println("INSERTED ORDERS SUCCESSFULLY")
	var orders []Order
	rows, err := db.Query(`
		SELECT id, amount, description, createdAt FROM orders WHERE user_id = $1;
	`, 1)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var order Order
		order.UserID = 1
		err := rows.Scan(&order.ID, &order.Amount, &order.Description, &order.CreatedAt)
		orders = append(orders, order)
		if err != nil {
			panic(err)
		}
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	fmt.Println(orders)

}
