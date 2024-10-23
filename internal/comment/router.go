package comment

import "github.com/gin-gonic/gin"

type CommentRouter struct {
	service *CommentService
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
	// posts := h.service.GetAllPosts()
	// c.JSON(200, gin.H{"posts": posts})
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
