package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/vladovsiychuk/microservice-demo-go/configs"
	"github.com/vladovsiychuk/microservice-demo-go/internal/auth"
	backendforfrontend "github.com/vladovsiychuk/microservice-demo-go/internal/backendforfrontend"
	"github.com/vladovsiychuk/microservice-demo-go/internal/comment"
	"github.com/vladovsiychuk/microservice-demo-go/internal/post"
	"github.com/vladovsiychuk/microservice-demo-go/internal/shared"
	websocketserver "github.com/vladovsiychuk/microservice-demo-go/internal/websocket_server"
	eventbus "github.com/vladovsiychuk/microservice-demo-go/pkg/event-bus"
	"github.com/vladovsiychuk/microservice-demo-go/pkg/helper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()
	r.Use(configs.CORSMiddleware())

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
	sessionTokenRepository := auth.NewSessionTokenRepository(postgresDB)

	authRepository := auth.NewKeyRepository(postgresDB)
	authService := auth.NewService(authRepository, sessionTokenRepository)
	authService.Init()
	authHandler := auth.NewRouter(authService)
	authHandler.RegisterRoutes(r)

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
	bffRouter := backendforfrontend.NewRouter(bffService, authService)
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
	host := helper.GetEnv("POSTGRES_HOST", "localhost")
	user := helper.GetEnv("POSTGRES_USER", "root")
	password := helper.GetEnv("POSTGRES_PASSWORD", "rootpassword")
	dbname := helper.GetEnv("POSTGRES_DB_NAME", "postgres")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", host, user, password, dbname)
	postgresDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return postgresDB
}

func setupRedis() *redis.Client {
	host := helper.GetEnv("REDIS_HOST", "localhost")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:6379", host),
		Password: "",
		DB:       0,
		Protocol: 2,
	})
	return redisClient
}

func setupMongoDb() *mongo.Database {
	host := helper.GetEnv("MONGODB_HOST", "localhost")

	uri := fmt.Sprintf("mongodb://root:example@%s:27017/test?authSource=admin", host)
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	return client.Database("test")
}
