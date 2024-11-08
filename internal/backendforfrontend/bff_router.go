package backendforfrontend

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vladovsiychuk/microservice-demo-go/internal/auth"
	customErrors "github.com/vladovsiychuk/microservice-demo-go/pkg/custom-errors"
)

type BffRouter struct {
	service     BffServiceI
	authService auth.AuthServiceI
}

func NewRouter(service BffServiceI, authService auth.AuthServiceI) *BffRouter {
	return &BffRouter{
		service,
		authService,
	}
}

func (h *BffRouter) RegisterRoutes(r *gin.Engine) {
	postGroup := r.Group("v1/posts")
	{
		postGroup.GET("/:postId", h.jwtAuthMiddleware, h.getPostAggregate)
	}
}

func (h *BffRouter) getPostAggregate(c *gin.Context) {
	postIdStr := c.Param("postId")

	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID format"})
		return
	}

	post, err := h.service.GetPostAggregate(postId)
	if err != nil {
		statusCode, response := customErrors.HandleError(err)
		c.JSON(statusCode, response)
		return
	}

	c.JSON(http.StatusCreated, post)
}

func (h *BffRouter) jwtAuthMiddleware(c *gin.Context) {
	tokenStr := c.GetHeader("Authorization")

	if !strings.HasPrefix(tokenStr, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

	tokenIsValid := h.authService.TokenIsValid(tokenStr)

	if !tokenIsValid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	c.Next()
}
