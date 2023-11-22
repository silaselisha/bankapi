package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Payload struct {
	Id        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	ExpiresAt time.Time `json:"expires_at"`
	IssuedAt  time.Time `json:"issued_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return &Payload{}, nil
	}
	return &Payload{
		Id: id,
		Username: username,
		ExpiresAt: time.Now().Add(duration),
		IssuedAt: time.Now(),
	}, nil
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiresAt) {
		err := errors.New("expired token")
		return fmt.Errorf("%w", err)
	}
	return nil
}