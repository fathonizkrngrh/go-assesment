package utils

import (
	"github.com/o1egl/paseto"

	"time"
)

func GenerateToken(userID string, roleID string, privateKey string) (string, error) {
	// Create a new PASETO token with an expiration time (e.g., 1 hour)
	token := paseto.JSONToken{
		Expiration: time.Now().Add(1 * time.Hour),
	}

	// Add the userID to the token's "subject" claim
	token.Subject = userID
	token.Set("roleID", roleID)

	// Create a PASETO token using the private key
	encoder := paseto.NewV2()
	tokenString, err := encoder.Sign(privateKey, token, nil)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
