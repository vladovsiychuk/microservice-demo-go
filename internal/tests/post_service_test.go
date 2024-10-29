package tests

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vladovsiychuk/microservice-demo-go/internal/post"
	"github.com/vladovsiychuk/microservice-demo-go/mocks"
)

func TestCreatePostSuccess(t *testing.T) {
	repository := mocks.NewPostRepositoryI(t)
	eventbus := mocks.NewEventBusI(t)

	repository.On("Create", mock.Anything).Return(nil).Once()
	eventbus.On("Publish", mock.Anything)

	service := post.NewService(repository, eventbus)
	response, err := service.CreatePost(post.PostRequest{Content: "foo", IsPrivate: false})

	assert.Equal(t, err, nil)
	assert.NotNil(t, response)

	repository.AssertExpectations(t)
}

func TestCreatePostFail(t *testing.T) {
	repository := mocks.NewPostRepositoryI(t)
	eventbus := mocks.NewEventBusI(t)

	repository.On("Create", mock.Anything).Return(errors.New("Error")).Once()

	service := post.NewService(repository, eventbus)
	_, err := service.CreatePost(post.PostRequest{Content: "foo", IsPrivate: false})

	assert.Equal(t, err, errors.New("Error"))
}
