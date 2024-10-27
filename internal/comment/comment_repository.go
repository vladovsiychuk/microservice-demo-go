package comment

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentRepository struct {
	postgresDB *gorm.DB
}

type CommentRepositoryI interface {
	Create(comment *Comment) error
	FindByKey(comment *Comment, commentId uuid.UUID) error
	Update(comment *Comment) error
}

func NewCommentRepository(postgresDB *gorm.DB) *CommentRepository {
	return &CommentRepository{
		postgresDB: postgresDB,
	}
}

func (r *CommentRepository) Create(comment *Comment) error {
	return r.postgresDB.Create(comment).Error
}

func (r *CommentRepository) FindByKey(comment *Comment, postId uuid.UUID) error {
	return r.postgresDB.Take(comment, postId).Error
}

func (r *CommentRepository) Update(comment *Comment) error {
	return r.postgresDB.Save(&comment).Error
}
