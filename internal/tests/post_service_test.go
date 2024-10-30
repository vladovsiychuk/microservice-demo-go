package tests

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vladovsiychuk/microservice-demo-go/internal/post"
	"github.com/vladovsiychuk/microservice-demo-go/mocks"
)

var anyValidPostRequest = post.PostRequest{Content: "foo", IsPrivate: false}
var anyPost, _ = post.CreatePost(anyValidPostRequest)
var errFoo = errors.New("Error")

func TestCreatePost(t *testing.T) {
	repository := mocks.NewPostRepositoryI(t)
	eventbus := mocks.NewEventBusI(t)
	service := post.NewService(repository, eventbus)

	originalCreatePostModel := post.CreatePost
	defer func() { post.CreatePost = originalCreatePostModel }()

	tests := []struct {
		name          string
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Success",
			setupMocks: func() {
				repository.On("Create", mock.Anything).Return(nil).Once()
				eventbus.On("Publish", mock.Anything)
				post.CreatePost = func(req post.PostRequest) (*post.Post, error) {
					return anyPost, nil
				}
			},
			expectedError: nil,
		},
		{
			name: "With Model Error",
			setupMocks: func() {
				post.CreatePost = func(req post.PostRequest) (*post.Post, error) {
					return nil, errFoo
				}
			},
			expectedError: errFoo,
		},
		{
			name: "With DB error",
			setupMocks: func() {
				repository.On("Create", mock.Anything).Return(errFoo).Once()
				post.CreatePost = func(req post.PostRequest) (*post.Post, error) {
					return anyPost, nil
				}
			},
			expectedError: errFoo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			response, err := service.CreatePost(anyValidPostRequest)

			assert.Equal(t, tt.expectedError, err)
			if tt.expectedError != nil {
				assert.Nil(t, response)
			} else {
				assert.NotNil(t, response)
			}

			repository.AssertExpectations(t)
			eventbus.AssertExpectations(t)
		})
	}
}
