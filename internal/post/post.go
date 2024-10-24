package post

import (
	"github.com/google/uuid"
	"github.com/vladovsiychuk/microservice-demo-go/pkg/customErrors"
)

type Post struct {
	Id        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	IsPrivate bool      `json:"isPrivate"`
}

func CreatePost(req PostRequest) (*Post, error) {
	if len(req.Content) > 100 {
		return nil, customErrors.NewBadRequestError("content exceeds 100 characters")
	}

	return &Post{
		uuid.New(),
		req.Content,
		req.IsPrivate,
	}, nil
}
