package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/vladovsiychuk/microservice-demo-go/configs"
	backendforfrontend "github.com/vladovsiychuk/microservice-demo-go/internal/backendforfrontend"
	"github.com/vladovsiychuk/microservice-demo-go/internal/comment"
	"github.com/vladovsiychuk/microservice-demo-go/internal/post"
	"github.com/vladovsiychuk/microservice-demo-go/internal/shared"
	eventbus "github.com/vladovsiychuk/microservice-demo-go/pkg/event-bus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // No password set
		DB:       0,  // Use default DB
		Protocol: 2,  // Connection protocol
	})

	mongoDB := setupMongoDb()

	eventBus := eventbus.NewEventBus()

	configs.SetupDbMigration(postgresDB)
	injectDependencies(postgresDB, mongoDB, redisClient, eventBus, r)

	r.Run(":8080")
}

func injectDependencies(
	postgresDB *gorm.DB,
	mongoDB *mongo.Database,
	redisClient *redis.Client,
	eventBus *eventbus.EventBus,
	r *gin.Engine,
) {
	postRepository := post.NewPostRepository(postgresDB)
	postService := post.NewService(postRepository, eventBus)
	postHandler := post.NewRouter(postService)
	postHandler.RegisterRoutes(r)

	commentRepository := comment.NewCommentRepository(postgresDB)
	commentService := comment.NewService(commentRepository, postService, eventBus)
	commentHandler := comment.NewRouter(commentService)
	commentHandler.RegisterRoutes(r)

	postAggregateRepository := backendforfrontend.NewPostAggregateRepository(mongoDB)
	redisCache := backendforfrontend.NewRedisRepository(redisClient)
	bffService := backendforfrontend.NewService(postAggregateRepository, redisCache, postService, commentService)
	bffRouter := backendforfrontend.NewRouter(bffService)
	bffRouter.RegisterRoutes(r)
	eventHandler := backendforfrontend.NewEventHandler(bffService)
	setupSubscribers(eventBus, eventHandler)
}

func setupSubscribers(eventBus *eventbus.EventBus, eventHandler *backendforfrontend.EventHandler) {
	postCreatedChan := make(chan eventbus.Event)
	eventBus.Subscribe(shared.PostCreatedEventType, postCreatedChan)
	go eventHandler.PostCreatedHandler(postCreatedChan)

	postUpdatedChan := make(chan eventbus.Event)
	eventBus.Subscribe(shared.PostUpdatedEventType, postUpdatedChan)
	go eventHandler.PostUpdatedHandler(postUpdatedChan)

	commentCreatedChan := make(chan eventbus.Event)
	eventBus.Subscribe(shared.CommentCreatedEventType, commentCreatedChan)
	go eventHandler.CommentCreatedHandler(commentCreatedChan)

	commentUpdatedChan := make(chan eventbus.Event)
	eventBus.Subscribe(shared.CommentUpdatedEventType, commentUpdatedChan)
	go eventHandler.CommentUpdatedHandler(commentUpdatedChan)
}

func setupMongoDb() *mongo.Database {
	uri := "mongodb://root:example@localhost:27017/test?authSource=admin"
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	return client.Database("test")
}
