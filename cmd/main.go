package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vladovsiychuk/microservice-demo-go/internal/comment"
	"github.com/vladovsiychuk/microservice-demo-go/internal/post"
)

func main() {
	r := gin.Default()

	postService := post.NewService()
	postHandler := post.NewRouter(postService)
	postHandler.RegisterRoutes(r)

	commentService := comment.NewService()
	commentHandler := comment.NewRouter(commentService)
	commentHandler.RegisterRoutes(r)

	r.Run(":8080")
}
