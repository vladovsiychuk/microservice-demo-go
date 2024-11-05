package comment

import (
	"github.com/google/uuid"
	customErrors "github.com/vladovsiychuk/microservice-demo-go/pkg/custom-errors"
)

type Comment struct {
	Id      uuid.UUID `json:"id"`
	PostId  uuid.UUID `json:"postId"`
	Content string    `json:"content"`
}

type CommentI interface {
	Update(req CommentRequest) error
}

var CreateComment = func(req CommentRequest, postId uuid.UUID, postIsPrivate bool) (CommentI, error) {
	if postIsPrivate {
		return nil, customErrors.NewBadRequestError("Comments cannot be added to private posts.")
	}

	return &Comment{
		uuid.New(),
		postId,
		req.Content,
	}, nil
}

func (p *Comment) Update(req CommentRequest) error {
	p.Content = req.Content
	return nil
}
