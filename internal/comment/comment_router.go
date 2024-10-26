package comment

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	customErrors "github.com/vladovsiychuk/microservice-demo-go/pkg/custom-errors"
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
	postGroup := r.Group("v1/posts/:postId/comments")
	{
		postGroup.POST("/", h.create)
		postGroup.PUT("/:commentId", h.update)
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

	comment, err := h.service.CreateComment(req, postId)
	if err != nil {
		statusCode, response := customErrors.HandleError(err)
		c.JSON(statusCode, response)
		return
	}

	c.JSON(http.StatusCreated, comment)
}

func (h *CommentRouter) update(c *gin.Context) {
	var req CommentRequest
	commentIdStr := c.Param("commentId")

	commentId, err := uuid.Parse(commentIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID format"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	comment, err := h.service.UpdateComment(req, commentId)
	if err != nil {
		statusCode, response := customErrors.HandleError(err)
		c.JSON(statusCode, response)
		return
	}

	c.JSON(http.StatusOK, comment)
}
