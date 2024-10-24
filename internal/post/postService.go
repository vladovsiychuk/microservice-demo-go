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

	if err := s.postgresDB.Create(post); err != nil {
		return nil, err.Error
	}

	return post, nil
}
