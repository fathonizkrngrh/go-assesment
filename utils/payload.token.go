package utils

import (
	"github.com/pkg/errors"
	"log"
	"time"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// Payload contains the payload data of the token
type Payload struct {
	UserID    string    `json:"user_id"`
	RoleID    string    `json:"role_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(userId string, roleId string) (*Payload, error) {
	payload := &Payload{
		UserID:    userId,
		RoleID:    roleId,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(24 * time.Hour),
	}
	log.Println(payload)
	return payload, nil
}
