package post

import (
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

func (s *PostService) GetAllPosts() []string {
	return []string{"Post 1", "Post 2"}
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
