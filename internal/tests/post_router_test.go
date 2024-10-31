package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vladovsiychuk/microservice-demo-go/internal/post"
	"github.com/vladovsiychuk/microservice-demo-go/mocks"
)

func TestCreatePostAndReturnJsonResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := mocks.NewPostServiceI(t)
	router := gin.Default()

	postRouter := post.NewRouter(mockService)
	postRouter.RegisterRoutes(router)

	fooPost := &post.Post{
		Id:        uuid.New(),
		Content:   "hello",
		IsPrivate: false,
	}

	reqBody := map[string]interface{}{
		"content":   "hello",
		"isPrivate": false,
	}

	expectedResponse := map[string]interface{}{
		"id":        fooPost.Id.String(),
		"content":   "hello",
		"isPrivate": false,
	}

	mockService.On("CreatePost", mock.Anything).Return(fooPost, nil)

	jsonReq, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/v1/posts/", bytes.NewBuffer(jsonReq))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	var responseBody map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, responseBody)

	mockService.AssertExpectations(t)
}
