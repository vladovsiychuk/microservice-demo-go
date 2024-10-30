package post

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	customErrors "github.com/vladovsiychuk/microservice-demo-go/pkg/custom-errors"
)

type PostRouter struct {
	service PostServiceI
}

type PostRequest struct {
	Content   string `json:"content" binding:"required"`
	IsPrivate bool   `json:"isPrivate"`
}

func NewRouter(service PostServiceI) *PostRouter {
	return &PostRouter{
		service: service,
	}
}

func (h *PostRouter) RegisterRoutes(r *gin.Engine) {
	postGroup := r.Group("v1/posts")
	{
		postGroup.POST("/", h.create)
		postGroup.PUT("/:postId", h.update)
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
		statusCode, response := customErrors.HandleError(err)
		c.JSON(statusCode, response)
		return
	}

	c.JSON(http.StatusCreated, post)
}

func (h *PostRouter) update(c *gin.Context) {
	var req PostRequest
	postIdStr := c.Param("postId")

	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID format"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	post, err := h.service.UpdatePost(postId, req)
	if err != nil {
		statusCode, response := customErrors.HandleError(err)
		c.JSON(statusCode, response)
		return
	}

	c.JSON(http.StatusOK, post)
}
