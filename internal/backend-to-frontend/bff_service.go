package backendtofrontend

import (
	"github.com/vladovsiychuk/microservice-demo-go/internal/comment"
	"github.com/vladovsiychuk/microservice-demo-go/internal/post"
)

type BffService struct {
}

type BffServiceI interface {
	CreatePostAggregate(*post.Post)
	UpdatePostAggregate(*post.Post)
	AddCommentToPostAggregate(*comment.Comment)
	UpdateCommentInPostAggregate(*comment.Comment)
}

func NewService() *BffService {
	return &BffService{}
}

func (s *BffService) CreatePostAggregate(*post.Post) {

}
func (s *BffService) UpdatePostAggregate(*post.Post) {

}
func (s *BffService) AddCommentToPostAggregate(*comment.Comment) {

}
func (s *BffService) UpdateCommentInPostAggregate(*comment.Comment) {

}
