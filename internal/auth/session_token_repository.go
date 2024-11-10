package auth

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionTokenRepository struct {
	postgresDB *gorm.DB
}

type SessionTokenRepositoryI interface {
	Create(sessionToken SessionTokenI) error
	FindById(sessionTokenId uuid.UUID) (SessionTokenI, error)
	Delete(sessionToken SessionTokenI) error
}

func NewSessionTokenRepository(postgresDB *gorm.DB) *SessionTokenRepository {
	return &SessionTokenRepository{
		postgresDB,
	}
}

func (r *SessionTokenRepository) Create(sessionToken SessionTokenI) error {
	return r.postgresDB.Create(sessionToken).Error
}

func (r *SessionTokenRepository) FindById(sessionTokenId uuid.UUID) (SessionTokenI, error) {
	var sessionToken SessionToken
	err := r.postgresDB.Take(&sessionToken, sessionTokenId).Error
	return &sessionToken, err
}

func (r *SessionTokenRepository) Delete(sessionToken SessionTokenI) error {
	return r.postgresDB.Delete(sessionToken).Error
}
