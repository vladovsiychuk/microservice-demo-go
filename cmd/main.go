package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
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

var jwtPrivateKey = getEnv("JWT_PRIVATE_KEY", "")
var jwtSecret = []byte(jwtPrivateKey)

func generateJWT(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func main() {
	r := gin.Default()

	sessionSecret := getEnv("SESSION_SECRET", "")
	googleClientKey := getEnv("GOOGLE_OAUTH_CLIENT_KEY", "")
	googleSecret := getEnv("GOOGLE_OAUTH_SECRET", "")

	gothic.Store = sessions.NewCookieStore([]byte(sessionSecret))

	goth.UseProviders(
		google.New(
			googleClientKey,
			googleSecret,
			"http://localhost:8080/auth/callback",
			"email", "profile",
		),
	)

	r.GET("/auth/login", func(c *gin.Context) {
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "provider", "google"))
		gothic.BeginAuthHandler(c.Writer, c.Request)
	})

	r.GET("/auth/callback", func(c *gin.Context) {
		user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			return
		}

		token, err := generateJWT(user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
			return
		}

		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "auth_token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,                    // Can't be accessed by JavaScript
			Secure:   true,                    // Use Secure if using HTTPS
			SameSite: http.SameSiteStrictMode, // Optional, for CSRF protection
			MaxAge:   3600,                    // Token expiry (1 hour)
		})

		c.Redirect(http.StatusFound, "http://localhost:3000/dashboard")
	})

	r.GET("/protected", jwtAuthMiddleware, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the protected route!"})
	})

	postgresDB := setupPostgres()
	redisClient := setupRedis()
	mongoDB := setupMongoDb()

	eventBus := eventbus.NewEventBus()

	configs.SetupDbMigration(postgresDB)
	injectDependencies(postgresDB, mongoDB, redisClient, eventBus, r)

	r.Run(":8080")
}

func jwtAuthMiddleware(c *gin.Context) {
	tokenStr := c.GetHeader("Authorization")

	if !strings.HasPrefix(tokenStr, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	c.Next()
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
