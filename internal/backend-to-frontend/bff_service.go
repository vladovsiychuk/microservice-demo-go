package backendtofrontend

import (
	"fmt"

	"github.com/vladovsiychuk/microservice-demo-go/internal/comment"
	"github.com/vladovsiychuk/microservice-demo-go/internal/post"
)

type BffService struct {
	repository PostAggregateRepositoryI
}

type BffServiceI interface {
	CreatePostAggregate(*post.Post)
	UpdatePostAggregate(*post.Post)
	AddCommentToPostAggregate(*comment.Comment)
	UpdateCommentInPostAggregate(*comment.Comment)
}

func NewService(repository PostAggregateRepositoryI) *BffService {
	return &BffService{
		repository: repository,
	}
}

func (s *BffService) CreatePostAggregate(post *post.Post) {
	postAggregate, err := CreatePostAggregate(post)
	if err != nil {
		fmt.Printf("Error occured during the creation of post aggregate: " + err.Error())
	}

	if err := s.repository.Create(postAggregate); err != nil {
		fmt.Printf("Error when saving to mongo db: " + err.Error())
	}
}
func (s *BffService) UpdatePostAggregate(*post.Post) {

}
func (s *BffService) AddCommentToPostAggregate(*comment.Comment) {

}
func (s *BffService) UpdateCommentInPostAggregate(*comment.Comment) {

}
