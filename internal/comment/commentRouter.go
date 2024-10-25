package comment

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vladovsiychuk/microservice-demo-go/pkg/customErrors"
)

type CommentRouter struct {
	service *CommentService
}

type CommentRequest struct {
	Content string `json:"content" binding:"required"`
}

func NewRouter(service *CommentService) *CommentRouter {
	return &CommentRouter{
		service: service,
	}
}

func (h *CommentRouter) RegisterRoutes(r *gin.Engine) {
	postGroup := r.Group("v1/posts/:postId/comment")
	{
		postGroup.POST("/", h.create)
		postGroup.PUT("/", h.update)
	}
}

func (h *CommentRouter) create(c *gin.Context) {
	var req CommentRequest
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

	post, err := h.service.CreateComment(req, postId)
	if err != nil {
		statusCode, response := customErrors.HandleError(err)
		c.JSON(statusCode, response)
		return
	}

	c.JSON(http.StatusCreated, post)
}

func (h *CommentRouter) update(c *gin.Context) {
	var newPost struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&newPost); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	// post := h.service.CreatePost(newPost.Content)
	// c.JSON(201, gin.H{"message": post})
}
