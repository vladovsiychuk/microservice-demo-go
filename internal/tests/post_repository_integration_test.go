package tests

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/vladovsiychuk/microservice-demo-go/configs"
	"github.com/vladovsiychuk/microservice-demo-go/internal/post"
	pgDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestPostRepository(t *testing.T) {
	ctx := context.Background()

	dbName := "testdb"
	dbUser := "user"
	dbPassword := "password"

	postgresContainer, err := postgres.Run(ctx,
		"postgres:latest",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	defer func() {
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return
	}

	host, _ := postgresContainer.Host(ctx)
	port, _ := postgresContainer.MappedPort(ctx, "5432")

	dsn := fmt.Sprintf("host=%s user=user password=password dbname=testdb port=%s sslmode=disable", host, port.Port())
	postgresDB, err := gorm.Open(pgDriver.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	configs.SetupDbMigration(postgresDB)
	repository := post.NewPostRepository(postgresDB)

	newPost, err := post.CreatePost(post.PostRequest{Content: "foo", IsPrivate: false})
	if err != nil {
		panic(err)
	}

	if err := repository.Create(newPost); err != nil {
		panic(err)
	}

	savedPost, err := repository.FindByKey(newPost.(*post.Post).Id)
	if err != nil {
		panic(err)
	}

	assert.NotNil(t, savedPost)
}
