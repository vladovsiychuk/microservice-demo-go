package post

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostRepository struct {
	postgresDB *gorm.DB
}

type PostRepositoryI interface {
	Create(post *Post) error
	FindByKey(post *Post, postId uuid.UUID) error
	Update(post *Post) error
}

func NewPostRepository(postgresDB *gorm.DB) *PostRepository {
	return &PostRepository{
		postgresDB: postgresDB,
	}
}

func (r *PostRepository) Create(post *Post) error {
	return r.postgresDB.Create(post).Error
}

func (r *PostRepository) FindByKey(post *Post, postId uuid.UUID) error {
	return r.postgresDB.Take(post, postId).Error
}

func (r *PostRepository) Update(post *Post) error {
	return r.postgresDB.Save(&post).Error
}
