package comment

import (
	"github.com/google/uuid"
	"github.com/vladovsiychuk/microservice-demo-go/internal/post"
	"gorm.io/gorm"
)

type CommentService struct {
	postgresDB  *gorm.DB
	postService *post.PostService
}

func NewService(postgresDB *gorm.DB, postService *post.PostService) *CommentService {
	return &CommentService{
		postgresDB:  postgresDB,
		postService: postService,
	}
}

func (s *CommentService) CreateComment(req CommentRequest, postId uuid.UUID) (*Comment, error) {
	postIsPrivate, err := s.postService.IsPrivate(postId)
	if err != nil {
		return nil, err
	}

	comment, err := CreateComment(req, postId, postIsPrivate)
	if err != nil {
		return nil, err
	}

	result := s.postgresDB.Create(comment)
	if result.Error != nil {
		return nil, result.Error
	}

	return comment, nil
}
