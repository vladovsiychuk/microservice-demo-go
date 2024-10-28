package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vladovsiychuk/microservice-demo-go/internal/post"
	"github.com/vladovsiychuk/microservice-demo-go/mocks"
)

func TestCreatePost(t *testing.T) {
	repository := mocks.NewPostRepositoryI(t)
	eventbus := mocks.NewEventBusI(t)

	repository.On("Create", mock.Anything).Return(nil).Once()
	eventbus.On("Publish", mock.Anything)

	service := post.NewService(repository, eventbus)
	response, err := service.CreatePost(post.PostRequest{Content: "foo", IsPrivate: false})

	assert.Equal(t, err, nil)
	assert.Equal(t, response.Content, "foo")

	repository.AssertExpectations(t)
}
