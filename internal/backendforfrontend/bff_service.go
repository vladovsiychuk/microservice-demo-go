package backendforfrontend

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/vladovsiychuk/microservice-demo-go/internal/comment"
	"github.com/vladovsiychuk/microservice-demo-go/internal/post"
	"go.mongodb.org/mongo-driver/mongo"
)

type BffService struct {
	repository PostAggregateRepositoryI
	redisCache RedisRepositoryI
}

type BffServiceI interface {
	GetPostAggregate(postId uuid.UUID) (PostAggregateI, error)
	CreatePostAggregate(*post.Post)
	UpdatePostAggregate(*post.Post)
	AddCommentToPostAggregate(*comment.Comment)
	UpdateCommentInPostAggregate(*comment.Comment)
}

func NewService(repository PostAggregateRepositoryI, redisCache RedisRepositoryI) *BffService {
	return &BffService{
		repository: repository,
		redisCache: redisCache,
	}
}

func (s *BffService) GetPostAggregate(postId uuid.UUID) (PostAggregateI, error) {
	cachedPostAgg, err := s.redisCache.FindByPostId(postId)
	if err == nil {
		return cachedPostAgg, nil
	} else if err != redis.Nil {
		return nil, err
	}

	return s.repository.FindById(postId)
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
	if err == mongo.ErrNoDocuments {
		s.CreatePostAggregate(post)
		return
	} else if err != nil {
		fmt.Printf("Error during post aggregate query: " + err.Error())
	}

	postAgg.Update(post)
	if err := s.repository.Update(postAgg); err != nil {
		fmt.Printf("Error during post update: " + err.Error())
	}

	s.redisCache.UpdateCache(postAgg)
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
