package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// Bytes wraps the crypto/rand pkg to return a random byte slice
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	nRead, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("bytes: %w", err)
	}
	if nRead < n {
		return nil, fmt.Errorf("bytes: didn't read enough bytes")
	}

	return b, nil
}

// String calls Bytes to return a random string of n number of bytes
func String(n int) (string, error) {
	b, err := Bytes(n)
	if err != nil {
		return "", fmt.Errorf("string: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
