package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vladovsiychuk/microservice-demo-go/internal/comment"
	"github.com/vladovsiychuk/microservice-demo-go/internal/post"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()

	dsn := "host=localhost user=root password=rootpassword dbname=postgres port=5432 sslmode=disable"
	postgresDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	postService := post.NewService(postgresDB)
	postHandler := post.NewRouter(postService)
	postHandler.RegisterRoutes(r)

	commentService := comment.NewService()
	commentHandler := comment.NewRouter(commentService)
	commentHandler.RegisterRoutes(r)

	r.Run(":8080")
}
