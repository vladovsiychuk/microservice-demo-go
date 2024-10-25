package post

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostService struct {
	postgresDB *gorm.DB
}

func NewService(postgresDB *gorm.DB) *PostService {
	return &PostService{
		postgresDB: postgresDB,
	}
}

func (s *PostService) CreatePost(req PostRequest) (*Post, error) {
	post, err := CreatePost(req)
	if err != nil {
		return nil, err
	}

	result := s.postgresDB.Create(post)
	if result.Error != nil {
		return nil, result.Error
	}

	return post, nil
}

func (s *PostService) UpdatePost(postId uuid.UUID, req PostRequest) (*Post, error) {
	var post Post

	if err := s.postgresDB.Take(&post, postId).Error; err != nil {
		return nil, errors.New("Post not found")
	}

	if err := post.Update(req); err != nil {
		return nil, err
	}

	if err := s.postgresDB.Save(&post).Error; err != nil {
		return nil, err
	}

	return &post, nil
}
