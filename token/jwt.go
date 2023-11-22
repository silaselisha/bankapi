package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)
const (
	SECRETE_KEY_LENGTH=32
)
type JwtMaker struct {
	SecreteKey string 
}

func NewJwtMaker(key string) (Maker, error) {
	if len(key) < SECRETE_KEY_LENGTH {
		return nil, fmt.Errorf("invalid jwt secrete key")
	}

	return &JwtMaker {
		SecreteKey: key,
	}, nil
}

func(m *JwtMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	result, err := token.SignedString([]byte(m.SecreteKey))
	if err != nil {
		return "", err
	}
	return result, nil
}

func(m *JwtMaker) VerifyToken(tok string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("invalid token header")
		}

		return []byte(m.SecreteKey), nil
	}

	token, err := jwt.ParseWithClaims(tok, &Payload{}, keyFunc)
	if err != nil {
		vErr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(vErr.Inner, fmt.Errorf("expired token")) {
			return nil, fmt.Errorf("expired token")
		}
		
		return nil, fmt.Errorf("invalid token")
	}

	payload, ok := token.Claims.(*Payload)
	if !ok {
		fmt.Println("here...")
		return nil, fmt.Errorf("invalid token")
	}
	
	return payload, nil
}