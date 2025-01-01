package models

import (
	"database/sql"
	"fmt"
	"io/fs"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
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

// Open will open a SQL conn - you will have to call defer db.Close() in main where conn is instantiated
func Open(config PostgresConf) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.generateDbConnStr())
	if err != nil {
		panic(err)
	} else {
		return db, err
	}
}

func Migrate(db *sql.DB, migrationsDir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	err = goose.Up(db, migrationsDir)
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	return nil
}

func MigrateFS(db *sql.DB, migrationsFS fs.FS, migrationsDir string) error {
	if migrationsDir == "" {
		migrationsDir = "."
	}
	goose.SetBaseFS(migrationsFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()
	return Migrate(db, migrationsDir)
}
