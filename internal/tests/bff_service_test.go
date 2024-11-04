package tests

import (
	"testing"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	bff "github.com/vladovsiychuk/microservice-demo-go/internal/backendforfrontend"
	"github.com/vladovsiychuk/microservice-demo-go/internal/comment"
	"github.com/vladovsiychuk/microservice-demo-go/internal/post"
	"github.com/vladovsiychuk/microservice-demo-go/mocks"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestGetPostAggregate(t *testing.T) {
	redisCache := mocks.NewRedisRepositoryI(t)
	repository := mocks.NewPostAggregateRepositoryI(t)
	postService := mocks.NewPostServiceI(t)
	commentService := mocks.NewCommentServiceI(t)
	mockedPost := mocks.NewPostI(t)
	mockedComments := []comment.CommentI{mocks.NewCommentI(t)}
	mockedPostAggregate := mocks.NewPostAggregateI(t)

	service := bff.NewService(repository, redisCache, postService, commentService)

	originalCreatePostAggregateWithComments := bff.CreatePostAggregateWithComments
	defer func() { bff.CreatePostAggregateWithComments = originalCreatePostAggregateWithComments }()

	tests := []struct {
		name           string
		setupMocks     func()
		expectedError  error
		expectedResult bff.PostAggregateI
	}{
		{
			name: "Success with cached post aggregate",
			setupMocks: func() {
				redisCache.On("FindByPostId", mock.Anything).Return(mockedPostAggregate, nil).Once()
			},
			expectedError:  nil,
			expectedResult: mockedPostAggregate,
		},
		{
			name: "Success with post aggregate from Mongo",
			setupMocks: func() {
				redisCache.On("FindByPostId", mock.Anything).Return(nil, redis.Nil).Once()
				repository.On("FindById", mock.Anything).Return(mockedPostAggregate, nil).Once()
			},
			expectedError:  nil,
			expectedResult: mockedPostAggregate,
		},
		{
			name: "Success with post aggregate created from post and comments",
			setupMocks: func() {
				redisCache.On("FindByPostId", mock.Anything).Return(nil, redis.Nil).Once()
				repository.On("FindById", mock.Anything).Return(nil, mongo.ErrNoDocuments).Once()
				postService.On("FindById", mock.Anything).Return(mockedPost, nil).Once()
				commentService.On("FindCommentsByPostId", mock.Anything).Return(mockedComments, nil).Once()
				bff.CreatePostAggregateWithComments = func(postI post.PostI, commentsI []comment.CommentI) (bff.PostAggregateI, error) {
					return mockedPostAggregate, nil
				}
			},
			expectedError:  nil,
			expectedResult: mockedPostAggregate,
		},
		{
			name: "With Redis error",
			setupMocks: func() {
				redisCache.On("FindByPostId", mock.Anything).Return(nil, errFoo).Once()
			},
			expectedError: errFoo,
		},
		{
			name: "With Mongo error",
			setupMocks: func() {
				redisCache.On("FindByPostId", mock.Anything).Return(nil, redis.Nil).Once()
				repository.On("FindById", mock.Anything).Return(nil, errFoo).Once()
			},
			expectedError: errFoo,
		},
		{
			name: "With post service error",
			setupMocks: func() {
				redisCache.On("FindByPostId", mock.Anything).Return(nil, redis.Nil).Once()
				repository.On("FindById", mock.Anything).Return(nil, mongo.ErrNoDocuments).Once()
				postService.On("FindById", mock.Anything).Return(nil, errFoo).Once()
				commentService.On("FindCommentsByPostId", mock.Anything).Return(nil, errFoo).Once()
			},
			expectedError: errFoo,
		},
		{
			name: "With comment service error",
			setupMocks: func() {
				redisCache.On("FindByPostId", mock.Anything).Return(nil, redis.Nil).Once()
				repository.On("FindById", mock.Anything).Return(nil, mongo.ErrNoDocuments).Once()
				postService.On("FindById", mock.Anything).Return(mockedPost, nil).Once()
				commentService.On("FindCommentsByPostId", mock.Anything).Return(nil, errFoo).Once()
			},
			expectedError: errFoo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			response, err := service.GetPostAggregate(uuid.New())

			assert.Equal(t, tt.expectedError, err)
			if tt.expectedError != nil {
				assert.Nil(t, response)
			} else {
				assert.Equal(t, tt.expectedResult, response)
			}

			redisCache.AssertExpectations(t)
			repository.AssertExpectations(t)
			postService.AssertExpectations(t)
			commentService.AssertExpectations(t)
		})
	}
}
