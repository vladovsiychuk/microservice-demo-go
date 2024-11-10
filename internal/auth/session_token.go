package auth

import (
	"time"

	"github.com/google/uuid"
)

type SessionToken struct {
	Id        uuid.UUID
	ExpiresAt time.Time
}

type SessionTokenI interface{}

var SESSION_TOKEN_DURATION = 10 * time.Second

var CreateSessionToken = func() SessionTokenI {
	return &SessionToken{
		uuid.New(),
		time.Now().Add(SESSION_TOKEN_DURATION),
	}
}
