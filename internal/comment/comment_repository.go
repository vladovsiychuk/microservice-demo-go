package comment

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentRepository struct {
	postgresDB *gorm.DB
}

type CommentRepositoryI interface {
	Create(comment CommentI) error
	FindByKey(commentId uuid.UUID) (CommentI, error)
	Update(comment CommentI) error
}

func NewCommentRepository(postgresDB *gorm.DB) *CommentRepository {
	return &CommentRepository{
		postgresDB: postgresDB,
	}
}

func (r *CommentRepository) Create(comment CommentI) error {
	return r.postgresDB.Create(comment).Error
}

func (r *CommentRepository) FindByKey(commentId uuid.UUID) (CommentI, error) {
	var comment Comment
	err := r.postgresDB.Take(&comment, commentId).Error
	return &comment, err
}

func (r *CommentRepository) Update(comment CommentI) error {
	return r.postgresDB.Save(comment).Error
}
