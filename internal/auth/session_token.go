package auth

import (
	"time"

	"github.com/google/uuid"
)

type SessionToken struct {
	Id        uuid.UUID
	Email     string
	ExpiresAt time.Time
}

type SessionTokenI interface{}

var SESSION_TOKEN_DURATION = 24 * time.Hour

var CreateSessionToken = func(email string) SessionTokenI {
	return &SessionToken{
		uuid.New(),
		email,
		time.Now().Add(SESSION_TOKEN_DURATION),
	}
}
