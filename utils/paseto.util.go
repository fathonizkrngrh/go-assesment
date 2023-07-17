package utils

import (
	"fmt"
	"github.com/aead/chacha20poly1305"
	"github.com/gocroot/gocroot/config"
	"github.com/o1egl/paseto"
	"log"
	"time"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey []byte) (*PasetoMaker, error) {
	log.Println(len(symmetricKey))
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

func (maker *PasetoMaker) CreateToken(userId string, roleId string) (string, error) {
	payload, err := NewPayload(userId, roleId)
	if err != nil {
		return "", err
	}
	log.Println(payload)

	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	return token, err
}

func (maker *PasetoMaker) ValidateToken(token string) (*Payload, error) {
	var payload Payload
	symnetricKey := []byte(config.PrivateKey)

	err := maker.paseto.Decrypt(token, symnetricKey, &payload, nil)
	if err != nil {
		return nil, err
	}

	// Verify token expiration
	if time.Now().After(payload.ExpiredAt) {
		return nil, fmt.Errorf("token has expired")
	}

	return &payload, nil
}
