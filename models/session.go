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
	// TODO: store session in DB

	return &session, nil
}

// TODO: implement
func (ss *SessionService) User(token string) (*User, error) {
	return nil, nil
}

func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
