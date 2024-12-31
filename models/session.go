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
}

type TokenManager struct {
	BytesPerToken int
}

func (tm *TokenManager) new() (token, tokenHash string, err error) {
	bytesPerToken := tm.BytesPerToken

	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}

	token, err = rand.String(bytesPerToken)
	if err != nil {
		return "", "", fmt.Errorf("create: %w", err)
	}
	tokenHash = tm.hash(token)
	return token, tokenHash, nil
}

func (tm *TokenManager) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	tm := TokenManager{}
	token, tokenHash, err := tm.new()
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	session := Session{
		UserID:    userID,
		Token:     token,
		TokenHash: tokenHash,
	}

	// 1. Try to update session & if no session -> create one
	row := ss.DB.QueryRow(`
		INSERT INTO sessions (user_id, token_hash) VALUES ($1, $2) 
		ON CONFLICT (user_id) 
		DO UPDATE SET token_hash = $2
		RETURNING id;`, session.UserID, session.TokenHash)

	err = row.Scan(&session.ID)

	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	var user User

	// 1. Hash the session token
	tm := TokenManager{}
	tokenHash := tm.hash(token)

	// 2. Query for user using the token hash
	row := ss.DB.QueryRow(`SELECT users.id, users.email, users.password_hash FROM sessions JOIN users ON users.id=sessions.user_id WHERE sessions.token_hash=$1;`, tokenHash)
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash)

	// 3. Return user
	return &user, err
}

func (ss *SessionService) Delete(token string) error {
	tm := TokenManager{}
	tokenHash := tm.hash(token)
	_, err := ss.DB.Exec(`DELETE FROM sessions WHERE token_hash=$1;`, tokenHash)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}
