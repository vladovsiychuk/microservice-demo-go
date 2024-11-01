package backendtofrontend

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type PostAggregateRepository struct {
	mongoColl *mongo.Collection
}

type PostAggregateRepositoryI interface {
	Create(post PostAggregateI) error
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
