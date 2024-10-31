package post

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostRepository struct {
	postgresDB *gorm.DB
}

type PostRepositoryI interface {
	Create(post PostI) error
	FindByKey(postId uuid.UUID) (PostI, error)
	Update(post PostI) error
}

func NewPostRepository(postgresDB *gorm.DB) *PostRepository {
	return &PostRepository{
		postgresDB: postgresDB,
	}
}

func (r *PostRepository) Create(post PostI) error {
	return r.postgresDB.Create(post).Error
}

func (r *PostRepository) FindByKey(postId uuid.UUID) (PostI, error) {
	var post Post
	err := r.postgresDB.Take(&post, postId).Error
	return &post, err
}

func (r *PostRepository) Update(post PostI) error {
	return r.postgresDB.Save(post).Error
}
