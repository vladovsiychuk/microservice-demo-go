package post

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vladovsiychuk/microservice-demo-go/pkg/errors"
)

type PostRouter struct {
	service *PostService
}

type PostRequest struct {
	Content   string `json:"content" binding:"required"`
	IsPrivate bool   `json:"isPrivate"`
}

func NewRouter(service *PostService) *PostRouter {
	return &PostRouter{
		service: service,
	}
}

func (h *PostRouter) RegisterRoutes(r *gin.Engine) {
	postGroup := r.Group("v1/posts")
	{
		postGroup.POST("/", h.create)
		postGroup.GET("/", h.getAllPosts)
	}
}

func (h *PostRouter) create(c *gin.Context) {
	var req PostRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	post, err := h.service.CreatePost(req)
	if err != nil {
		statusCode, response := errors.HandleError(err)
		c.JSON(statusCode, response)
		return
	}

	c.JSON(http.StatusCreated, post)
}

func (h *PostRouter) getAllPosts(c *gin.Context) {
	posts := h.service.GetAllPosts()
	c.JSON(200, gin.H{"posts": posts})
}
