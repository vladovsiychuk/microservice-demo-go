package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/vladovsiychuk/microservice-demo-go/configs"
	backendforfrontend "github.com/vladovsiychuk/microservice-demo-go/internal/backendforfrontend"
	"github.com/vladovsiychuk/microservice-demo-go/internal/comment"
	"github.com/vladovsiychuk/microservice-demo-go/internal/post"
	"github.com/vladovsiychuk/microservice-demo-go/internal/shared"
	websocketserver "github.com/vladovsiychuk/microservice-demo-go/internal/websocket_server"
	eventbus "github.com/vladovsiychuk/microservice-demo-go/pkg/event-bus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()

	postgresDB := setupPostgres()
	redisClient := setupRedis()
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

	websocketService := websocketserver.NewService()
	websocketHander := websocketserver.NewRouter(websocketService)
	websocketHander.RegisterRoutes(r)
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

func setupPostgres() *gorm.DB {
	host := getEnv("POSTGRES_HOST", "localhost")
	user := getEnv("POSTGRES_USER", "root")
	password := getEnv("POSTGRES_PASSWORD", "rootpassword")
	dbname := getEnv("POSTGRES_DB_NAME", "postgres")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", host, user, password, dbname)
	postgresDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return postgresDB
}

func setupRedis() *redis.Client {
	host := getEnv("REDIS_HOST", "localhost")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:6379", host),
		Password: "",
		DB:       0,
		Protocol: 2,
	})
	return redisClient
}

func setupMongoDb() *mongo.Database {
	host := getEnv("MONGODB_HOST", "localhost")

	uri := fmt.Sprintf("mongodb://root:example@%s:27017/test?authSource=admin", host)
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	return client.Database("test")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
