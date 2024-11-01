package backendforfrontend

import (
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostAggregateRepository struct {
	mongoColl *mongo.Collection
}

type PostAggregateRepositoryI interface {
	Create(post PostAggregateI) error
	FindById(postId uuid.UUID) (PostAggregateI, error)
	Update(post PostAggregateI) error
}

func NewPostAggregateRepository(mongoDB *mongo.Database) *PostAggregateRepository {
	return &PostAggregateRepository{
		mongoDB.Collection("posts"),
	}
}

func (a *PostAggregateRepository) Create(post PostAggregateI) error {
	_, err := a.mongoColl.InsertOne(context.TODO(), post)
	return err
}

func (a *PostAggregateRepository) FindById(postId uuid.UUID) (PostAggregateI, error) {
	var post PostAggregate
	err := a.mongoColl.FindOne(context.TODO(), bson.D{{Key: "_id", Value: postId}}).Decode(&post)
	return &post, err
}

func (a *PostAggregateRepository) Update(post PostAggregateI) error {
	_, err := a.mongoColl.ReplaceOne(
		context.TODO(),
		bson.D{{Key: "_id", Value: post.(*PostAggregate).Id}},
		post,
	)

	return err
}
