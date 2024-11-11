package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vladovsiychuk/microservice-demo-go/internal/auth"
)

func TestRotation(t *testing.T) {
	keys := auth.CreateKeys()

	initialPrivateKey, _ := keys.GetPrivateKey()
	initialPublicKey, _ := keys.GetPublicKey()

	keys.Rotate()

	privateKey, _ := keys.GetPrivateKey()
	secondaryPublicKey, _ := keys.GetSecondaryPulicKey()

	assert.Equal(t, secondaryPublicKey, initialPublicKey)
	assert.NotEqual(t, privateKey, initialPrivateKey)
}
