package backendforfrontend

import (
	"github.com/google/uuid"
	"github.com/vladovsiychuk/microservice-demo-go/internal/comment"
	"github.com/vladovsiychuk/microservice-demo-go/internal/post"
)

type CommentItem struct {
	Id      uuid.UUID `json:"id"`
	Content string    `json:"content"`
}

type PostAggregate struct {
	Id        uuid.UUID     `json:"id" bson:"_id"`
	Content   string        `json:"content"`
	IsPrivate bool          `json:"is_private" bson:"is_private"`
	Comments  []CommentItem `json:"comments"`
}

type PostAggregateI interface {
	Update(post *post.Post) error
	AddComment(comment *comment.Comment) error
	UpdateComment(comment *comment.Comment) error
}

var CreatePostAggregate = func(post *post.Post) (PostAggregateI, error) {
	return &PostAggregate{
		post.Id,
		post.Content,
		post.IsPrivate,
		[]CommentItem{},
	}, nil
}

var CreatePostAggregateWithComments = func(postI post.PostI, commentsI []comment.CommentI) (PostAggregateI, error) {
	var commentItems []CommentItem
	post := postI.(*post.Post)

	for _, commentI := range commentsI {
		comment := commentI.(*comment.Comment)
		commentItems = append(commentItems, CommentItem{comment.Id, comment.Content})
	}

	return &PostAggregate{
		post.Id,
		post.Content,
		post.IsPrivate,
		commentItems,
	}, nil
}

func (a *PostAggregate) Update(post *post.Post) error {
	a.Content = post.Content
	a.IsPrivate = post.IsPrivate
	return nil
}

func (a *PostAggregate) AddComment(comment *comment.Comment) error {
	a.Comments = append(a.Comments, CommentItem{comment.Id, comment.Content})
	return nil
}

func (a *PostAggregate) UpdateComment(updatedComment *comment.Comment) error {
	for i := range a.Comments {
		if a.Comments[i].Id == updatedComment.Id {
			a.Comments[i].Content = updatedComment.Content
			break
		}
	}
	return nil
}
