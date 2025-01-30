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

	if err != nil {
		fmt.Printf("An error occured hashing password: %v\n", err)
		return nil, err
	}
	hashedPassword := string(hashedBytes)
	user := User{
		Email:        email,
		PasswordHash: hashedPassword,
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

func (us *UserService) UpdatePassword(userID int, pw string) error {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)

	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}

	hashedPassword := string(hashedBytes)

	_, err = us.DB.Exec(`
		UPDATE users SET password_hash = $2 WHERE id = $1
	`, userID, hashedPassword)

	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	return nil
}
