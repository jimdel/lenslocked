package models

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	DefaultResetDuration = 1 * time.Hour
)

type PasswordReset struct {
	ID     int
	UserID int
	// Only set when creating a new password_reset request
	// When looking up a password_reset only the token hash will be available
	Token     string
	TokenHash string
	ExpiresAt time.Time
}

type PasswordResetService struct {
	DB            *sql.DB
	BytesPerToken int
	// The amount of time that a pw reset is valid for.
	// Defaults to DefaultResetDuration
	Duration time.Duration
}

func (svc *PasswordResetService) Create(email string) (*PasswordReset, error) {
	return nil, fmt.Errorf("TODO: implement PasswordResetService.Create")
}

func (svc *PasswordResetService) Consume(token string) (*User, error) {
	return nil, fmt.Errorf("TODO: implement PasswordResetService.Consume")
}
