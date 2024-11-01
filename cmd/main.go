package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vladovsiychuk/microservice-demo-go/configs"
	backendtofrontend "github.com/vladovsiychuk/microservice-demo-go/internal/backend-to-frontend"
	"github.com/vladovsiychuk/microservice-demo-go/internal/comment"
	"github.com/vladovsiychuk/microservice-demo-go/internal/post"
	"github.com/vladovsiychuk/microservice-demo-go/internal/shared"
	eventbus "github.com/vladovsiychuk/microservice-demo-go/pkg/event-bus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()

	// setup postgres DB
	dsn := "host=localhost user=root password=rootpassword dbname=postgres port=5432 sslmode=disable"
	postgresDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	eventBus := eventbus.NewEventBus()

	setupSubscribers(eventBus)
	configs.SetupDbMigration(postgresDB)
	injectDependencies(postgresDB, eventBus, r)

	r.Run(":8080")
}

func setupSubscribers(eventBus *eventbus.EventBus) {
	postCreatedChan := make(chan eventbus.Event)
	eventBus.Subscribe(shared.PostCreatedEventType, postCreatedChan)
	go backendtofrontend.PostCreatedHandler(postCreatedChan)

	postUpdatedChan := make(chan eventbus.Event)
	eventBus.Subscribe(shared.PostUpdatedEventType, postUpdatedChan)
	go backendtofrontend.PostUpdatedHandler(postUpdatedChan)

	commentCreatedChan := make(chan eventbus.Event)
	eventBus.Subscribe(shared.CommentCreatedEventType, commentCreatedChan)
	go backendtofrontend.CommentCreatedHandler(commentCreatedChan)

	commentUpdatedChan := make(chan eventbus.Event)
	eventBus.Subscribe(shared.CommentUpdatedEventType, commentUpdatedChan)
	go backendtofrontend.CommentUpdatedHandler(commentUpdatedChan)
}

func injectDependencies(postgresDB *gorm.DB, eventBus *eventbus.EventBus, r *gin.Engine) {
	postRepository := post.NewPostRepository(postgresDB)
	postService := post.NewService(postRepository, eventBus)
	postHandler := post.NewRouter(postService)
	postHandler.RegisterRoutes(r)

	commentRepository := comment.NewCommentRepository(postgresDB)
	commentService := comment.NewService(commentRepository, postService, eventBus)
	commentHandler := comment.NewRouter(commentService)
	commentHandler.RegisterRoutes(r)
}
