package main

import (
	"embed"

	"github.com/gin-gonic/gin"
	"github.com/pressly/goose/v3"
	backendtofrontend "github.com/vladovsiychuk/microservice-demo-go/internal/backendToFrontend"
	"github.com/vladovsiychuk/microservice-demo-go/internal/comment"
	"github.com/vladovsiychuk/microservice-demo-go/internal/post"
	eventbus "github.com/vladovsiychuk/microservice-demo-go/pkg/eventBus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	r := gin.Default()

	eventBus := eventbus.NewEventBus()

	userRegisteredChan := make(chan eventbus.Event)
	eventBus.Subscribe("UserRegistered", userRegisteredChan)
	go backendtofrontend.UserRegisteredHandler(userRegisteredChan)

	userRegisteredChan2 := make(chan eventbus.Event)
	eventBus.Subscribe("UserRegistered", userRegisteredChan2)
	go backendtofrontend.UserRegisteredHandler2(userRegisteredChan2)

	// setup postgres DB
	dsn := "host=localhost user=root password=rootpassword dbname=postgres port=5432 sslmode=disable"
	postgresDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	setupDbMigration(postgresDB)

	postService := post.NewService(postgresDB)
	postHandler := post.NewRouter(postService)
	postHandler.RegisterRoutes(r)

	commentService := comment.NewService(postgresDB, postService, eventBus)
	commentHandler := comment.NewRouter(commentService)
	commentHandler.RegisterRoutes(r)

	r.Run(":8080")
}

func setupDbMigration(postgresDB *gorm.DB) {
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
