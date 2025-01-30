package models

import (
	"database/sql"
	"fmt"
	"strings"
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
	// Verify we have a valid email address for a user & grab userId
	email = strings.ToLower(email)
	var userID int
	row := svc.DB.QueryRow(`SELECT id FROM users WHERE email = $1;`, email)
	err := row.Scan(&userID)
	if err != nil {
		// TODO: consider returning a specific error when the user does not exist
		return nil, fmt.Errorf("create: %w", err)
	}

	// Build the PasswordReset struct
	tm := TokenManager{}
	token, tokenHash, err := tm.New()
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	duration := svc.Duration
	if duration == 0 {
		duration = DefaultResetDuration
	}

	pwr := PasswordReset{
		UserID:    userID,
		Token:     token,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(duration),
	}
	// Insert pwr into the DB
	row = svc.DB.QueryRow(`
		INSERT INTO password_resets (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3) ON CONFLICT (user_id) DO
		UPDATE
		SET token_hash = $2, expires_at = $3
		RETURNING id;
	`, pwr.UserID, pwr.TokenHash, pwr.ExpiresAt)

	err = row.Scan(&pwr.ID)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	return &pwr, nil
}

func (svc *PasswordResetService) Consume(token string) (*User, error) {
	//Hash token - query db for user
	tm := TokenManager{}
	var user User
	var pwr PasswordReset

	tokenHash := tm.Hash(token)

	row := svc.DB.QueryRow(`
		SELECT 
		   password_resets.id,
		   password_resets.expires_at,
		   users.id,
		   users.email,
		   users.password_hash
		FROM 
			password_resets
		JOIN 
			users on users.id = password_resets.user_id
		WHERE 
			password_resets.token_hash = $1;`, tokenHash)

	err := row.Scan(&pwr.ID, &pwr.ExpiresAt, &user.ID, &user.Email, &user.PasswordHash)

	// Handle can't find token
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}

	//Handle expired token
	if time.Now().After(pwr.ExpiresAt) {
		return nil, fmt.Errorf("token expired: %v", err)
	}

	// Delete pwr from DB after finding user
	err = svc.delete(pwr.ID)
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}

	return &user, nil
}

func (svc *PasswordResetService) delete(id int) error {
	_, err := svc.DB.Exec(`
		DELETE FROM password_resets
		WHERE id = $1;`, id)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}
