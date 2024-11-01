package backendforfrontend

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	customErrors "github.com/vladovsiychuk/microservice-demo-go/pkg/custom-errors"
)

type BffRouter struct {
	service BffServiceI
}

func NewRouter(service BffServiceI) *BffRouter {
	return &BffRouter{
		service: service,
	}
}

func (h *BffRouter) RegisterRoutes(r *gin.Engine) {
	postGroup := r.Group("v1/posts")
	{
		postGroup.GET("/:postId", h.getPostAggregate)
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
