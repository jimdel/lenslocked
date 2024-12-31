package models

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type PostgresConf struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func DefaultPostgresConfig() PostgresConf {
	var conf = PostgresConf{
		Host:     "localhost",
		Port:     "8081",
		User:     "baloo",
		Password: "junglebook",
		DBName:   "lenslocked",
		SSLMode:  "disable",
	}
	return conf
}

func (c PostgresConf) generateDbConnStr() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// Open will open a SQL conn - you will have to call Close() afterward
func Open(config PostgresConf) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.generateDbConnStr())
	if err != nil {
		panic(err)
	} else {
		fmt.Println("CONNECTED SUCCESSFULLY")
		return db, err
	}
}
