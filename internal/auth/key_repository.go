package auth

import "gorm.io/gorm"

type KeyRepository struct {
	postgresDB *gorm.DB
}

type KeyRepositoryI interface{}

func NewKeyRepository(postgresDB *gorm.DB) *KeyRepository {
	return &KeyRepository{
		postgresDB,
	}
}

func (r *KeyRepository) Update(keys KeysI) error {
	r.postgresDB.Delete(&Keys{}, "1 = 1")
	return r.postgresDB.Create(keys).Error
}

func (r *KeyRepository) GetKeys() (KeysI, error) {
	var keys Keys
	err := r.postgresDB.First(&keys).Error
	return &keys, err
}
