package backendtofrontend

import (
	"github.com/google/uuid"
	"github.com/vladovsiychuk/microservice-demo-go/internal/post"
)

type CommentItem struct {
	Id      uuid.UUID
	Content string
}

type PostAggregate struct {
	Id        uuid.UUID `bson:"_id"`
	Content   string
	IsPrivate bool `bson:"is_private"`
	comments  []CommentItem
}

type PostAggregateI interface {
	Update(post *post.Post) error
}

var CreatePostAggregate = func(post *post.Post) (PostAggregateI, error) {
	return &PostAggregate{
		post.Id,
		post.Content,
		post.IsPrivate,
		[]CommentItem{},
	}, nil
}

func (a *PostAggregate) Update(post *post.Post) error {
	a.Content = post.Content
	a.IsPrivate = post.IsPrivate
	return nil
}
