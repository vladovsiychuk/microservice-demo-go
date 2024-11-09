package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vladovsiychuk/microservice-demo-go/internal/auth"
)

func TestRotation(t *testing.T) {
	var keys = &auth.Keys{PrivateKey: "a", PublicKey: "b", SecondaryPublicKey: "c"}

	keys.Rotate()

	assert.Equal(t, keys.SecondaryPublicKey, "b")
	assert.NotEqual(t, keys.PrivateKey, "a")
}
