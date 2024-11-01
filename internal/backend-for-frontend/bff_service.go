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

func (s *BffService) UpdatePostAggregate(post *post.Post) {
	postAgg, err := s.repository.FindById(post.Id)
	if err != nil {
		fmt.Printf("Error during post aggregate query: " + err.Error())
	}

	postAgg.Update(post)
	if err := s.repository.Update(postAgg); err != nil {
		fmt.Printf("Error during post update: " + err.Error())
	}
}

func (s *BffService) AddCommentToPostAggregate(comment *comment.Comment) {
	postAgg, err := s.repository.FindById(comment.PostId)
	if err != nil {
		fmt.Printf("Error during post aggregate query: " + err.Error())
	}

	postAgg.AddComment(comment)
	if err := s.repository.Update(postAgg); err != nil {
		fmt.Printf("Error during post update: " + err.Error())
	}

}

func (s *BffService) UpdateCommentInPostAggregate(comment *comment.Comment) {
	postAgg, err := s.repository.FindById(comment.PostId)
	if err != nil {
		fmt.Printf("Error during post aggregate query: " + err.Error())
	}

	postAgg.UpdateComment(comment)
	if err := s.repository.Update(postAgg); err != nil {
		fmt.Printf("Error during post update: " + err.Error())
	}
}
