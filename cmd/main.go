package main

import (
	"embed"

	"github.com/gin-gonic/gin"
	"github.com/pressly/goose/v3"
	"github.com/vladovsiychuk/microservice-demo-go/internal/comment"
	"github.com/vladovsiychuk/microservice-demo-go/internal/post"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	r := gin.Default()

	dsn := "host=localhost user=root password=rootpassword dbname=postgres port=5432 sslmode=disable"
	postgresDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	dbMigration(postgresDB)

	postService := post.NewService(postgresDB)
	postHandler := post.NewRouter(postService)
	postHandler.RegisterRoutes(r)

	commentService := comment.NewService()
	commentHandler := comment.NewRouter(commentService)
	commentHandler.RegisterRoutes(r)

	r.Run(":8080")
}

func dbMigration(postgresDB *gorm.DB) {
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}
	db, err := postgresDB.DB()
	if err != nil {
		panic(err)
	}
	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}
}
