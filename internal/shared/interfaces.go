package shared

import (
	"github.com/google/uuid"
)

type PostServiceI interface {
	IsPrivate(postId uuid.UUID) (bool, error)
}
