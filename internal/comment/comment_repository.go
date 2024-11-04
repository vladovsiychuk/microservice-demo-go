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
	FindById(commentId uuid.UUID) (CommentI, error)
	FindCommentsByPostId(postId uuid.UUID) ([]CommentI, error)
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

func (r *CommentRepository) FindById(commentId uuid.UUID) (CommentI, error) {
	var comment Comment
	err := r.postgresDB.Take(&comment, commentId).Error
	return &comment, err
}

func (r *CommentRepository) FindCommentsByPostId(postId uuid.UUID) ([]CommentI, error) {
	var comments []Comment
	err := r.postgresDB.Find(&comments, "post_id = ?", postId).Error

	var result []CommentI
	for _, comment := range comments {
		result = append(result, &comment)
	}

	return result, err
}

func (r *CommentRepository) Update(comment CommentI) error {
	return r.postgresDB.Save(comment).Error
}
