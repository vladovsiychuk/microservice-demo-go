package comment

import (
	"github.com/google/uuid"
	customErrors "github.com/vladovsiychuk/microservice-demo-go/pkg/custom-errors"
)

type Comment struct {
	Id      uuid.UUID `json:"id"`
	PostId  uuid.UUID `json:"post_id"`
	Content string    `json:"content"`
}

func CreateComment(req CommentRequest, postId uuid.UUID, postIsPrivate bool) (*Comment, error) {
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
