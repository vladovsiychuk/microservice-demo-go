package shared

import (
	"github.com/google/uuid"
)

type PostServiceSharedI interface {
	IsPrivate(postId uuid.UUID) (bool, error)
}
