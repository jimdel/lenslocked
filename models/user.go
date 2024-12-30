package models

import (
	"database/sql"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UnauthenticatedUser struct {
	Email    string
	Password string
}

type User struct {
	ID           int
	Email        string
	PasswordHash string
}

type UserService struct {
	DB *sql.DB
}

func (us *UserService) Create(nu UnauthenticatedUser) (*User, error) {

	email := strings.ToLower(nu.Email)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	hashedPassword := string(hashedBytes)
	user := User{
		Email:        email,
		PasswordHash: hashedPassword,
	}

	if err != nil {
		fmt.Printf("An error occured hashing password: %v\n", err)
		return nil, err
	}

	row := us.DB.QueryRow(`
		INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id;
	`, email, hashedPassword)

	var userId int
	err = row.Scan(&userId)

	if err != nil {
		fmt.Printf("An error occured scanning row: %v\n", err)
		return nil, err
	}

	user.ID = userId
	return &user, nil
}

func (us *UserService) Authenticate(nu UnauthenticatedUser) (*User, error) {
	email := strings.ToLower(nu.Email)

	user := User{
		Email: email,
	}
	row := us.DB.QueryRow(`
		SELECT id, password_hash from users WHERE email=$1;
	`, email)

	err := row.Scan(&user.ID, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("authentication error: %v", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(nu.Password))

	if err != nil {
		return nil, fmt.Errorf("invalid password: %v", err)
	}
	return &user, err
}
