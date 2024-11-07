package auth

import "gorm.io/gorm"

type KeyRepository struct {
	postgresDB *gorm.DB
}

func (r *KeyRepository) Create(keys KeysI) error {
	return r.postgresDB.Create(keys).Error
}
