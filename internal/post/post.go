package post

import (
	"github.com/google/uuid"
	customErrors "github.com/vladovsiychuk/microservice-demo-go/pkg/custom-errors"
)

type Post struct {
	Id        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	IsPrivate bool      `json:"isPrivate"`
}

type PostI interface {
	Update(req PostRequest) error
}

var CreatePost = func(req PostRequest) (PostI, error) {
	if len(req.Content) > 10 {
		return nil, customErrors.NewBadRequestError("content exceeds 10 characters")
	}

	return &Post{
		uuid.New(),
		req.Content,
		req.IsPrivate,
	}, nil
}

func (p *Post) Update(req PostRequest) error {
	p.Content = req.Content
	p.IsPrivate = req.IsPrivate
	return nil
}
