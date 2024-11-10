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

var CreateSessionToken = func() SessionTokenI {
	return &SessionToken{
		uuid.New(),
		time.Now().Add(5 * time.Minute),
	}
}
