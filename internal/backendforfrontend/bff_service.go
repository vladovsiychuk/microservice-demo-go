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
	repository     PostAggregateRepositoryI
	redisCache     RedisRepositoryI
	postService    post.PostServiceI
	commentService comment.CommentServiceI
}

type BffServiceI interface {
	GetPostAggregate(postId uuid.UUID) (PostAggregateI, error)
	CreatePostAggregate(*post.Post)
	UpdatePostAggregate(*post.Post)
	AddCommentToPostAggregate(*comment.Comment)
	UpdateCommentInPostAggregate(*comment.Comment)
}

func NewService(
	repository PostAggregateRepositoryI,
	redisCache RedisRepositoryI,
	postService post.PostServiceI,
	commentService comment.CommentServiceI,
) *BffService {
	return &BffService{
		repository,
		redisCache,
		postService,
		commentService,
	}
}

func (s *BffService) GetPostAggregate(postId uuid.UUID) (PostAggregateI, error) {
	cachedPostAgg, err := s.redisCache.FindByPostId(postId)
	if err == nil {
		return cachedPostAgg, nil
	} else if err != redis.Nil {
		return nil, err
	}

	postAggFromMongo, err := s.repository.FindById(postId)
	if err == nil {
		return postAggFromMongo, nil
	} else if err != mongo.ErrNoDocuments {
		return nil, err
	}

	postChannel := make(chan post.PostI, 1)
	commentsChannel := make(chan []comment.CommentI, 1)
	errorChannel := make(chan error, 2)

	go func() {
		post, err := s.postService.FindById(postId)
		if err != nil {
			errorChannel <- err
			return
		}
		postChannel <- post
	}()

	go func() {
		comments, err := s.commentService.FindCommentsByPostId(postId)
		if err != nil {
			errorChannel <- err
			return
		}
		commentsChannel <- comments
	}()

	var post post.PostI
	var comments []comment.CommentI
	for i := 0; i < 2; i++ {
		select {
		case p := <-postChannel:
			post = p
		case c := <-commentsChannel:
			comments = c
		case err := <-errorChannel:
			return nil, err
		}
	}

	return CreatePostAggregateWithComments(post, comments)
}

func (s *BffService) CreatePostAggregate(post *post.Post) {
	postAggregate, err := CreatePostAggregate(post)
	if err != nil {
		fmt.Printf("Error occured during the creation of post aggregate: " + err.Error())
		return
	}

	if err := s.repository.Create(postAggregate); err != nil {
		fmt.Printf("Error when saving to mongo db: " + err.Error())
	}
}

func (s *BffService) UpdatePostAggregate(post *post.Post) {
	postAgg, err := s.repository.FindById(post.Id)
	if err != nil {
		fmt.Printf("Error during post aggregate query: " + err.Error())
		return
	}

	postAgg.Update(post)
	if err := s.repository.Update(postAgg); err != nil {
		fmt.Printf("Error during post update: " + err.Error())
		return
	}

	s.redisCache.UpdateCache(postAgg)
}

func (s *BffService) AddCommentToPostAggregate(comment *comment.Comment) {
	postAgg, err := s.repository.FindById(comment.PostId)
	if err != nil {
		fmt.Printf("Error during post aggregate query: " + err.Error())
		return
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
		return
	}

	postAgg.UpdateComment(comment)
	if err := s.repository.Update(postAgg); err != nil {
		fmt.Printf("Error during post update: " + err.Error())
	}
}
