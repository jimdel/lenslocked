package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/jimdel/lenslocked/rand"
)

const (
	// Min # of bytes for session token
	MinBytesPerToken = 32
)

type Session struct {
	ID     int
	UserID int
	// Only set when creating a new session
	// When looking up a session only the token hash will be available
	Token     string
	TokenHash string
}

type SessionService struct {
	DB *sql.DB
	// Min # of bytes for session token, if not set default to MinBytesPerToken
	BytesPerToken int
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	bytesPerToken := ss.BytesPerToken

	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}

	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	session := Session{
		UserID:    userID,
		Token:     token,
		TokenHash: ss.hash(token),
	}

	// 1. Try to update session
	// 2. If no session create one
	row := ss.DB.QueryRow(`UPDATE sessions SET token_hash=$2 WHERE user_id=$1 RETURNING id;`, session.UserID, session.TokenHash)

	err = row.Scan(&session.ID)
	if err == sql.ErrNoRows {
		row = ss.DB.QueryRow(`INSERT INTO sessions (user_id, token_hash) VALUES ($1, $2) RETURNING id;`, session.UserID, session.TokenHash)
		err = row.Scan(&session.ID)
	}

	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	var user User

	// 1. Hash the session token
	tokenHash := ss.hash(token)

	// 2. Query for the session with the token hash
	row := ss.DB.QueryRow(`SELECT user_id FROM sessions WHERE token_hash=$1;`, tokenHash)
	err := row.Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}

	// 3. Using UserID from session query for the user
	row = ss.DB.QueryRow(`SELECT email, password_hash FROM users WHERE id=$1`, user.ID)
	err = row.Scan(&user.Email, &user.PasswordHash)

	// 4. Return user
	return &user, err
}

func (ss *SessionService) Delete(token string) error {
	tokenHash := ss.hash(token)
	_, err := ss.DB.Exec(`DELETE FROM sessions WHERE token_hash=$1;`, tokenHash)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
