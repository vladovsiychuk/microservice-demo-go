package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
)

type Keys struct {
	PrivateKey         string
	PublicKey          string
	SecondaryPublicKey string
}

type KeysI interface {
	Rotate()
}

var CreateKeys = func() KeysI {
	privateKeyStr, publicKeyStr := generateRandomKeys()

	return &Keys{
		privateKeyStr,
		publicKeyStr,
		publicKeyStr,
	}
}

func (k *Keys) Rotate() {
	privateKeyStr, publicKeyStr := generateRandomKeys()

	k.PrivateKey = privateKeyStr
	k.SecondaryPublicKey = k.PublicKey
	k.PublicKey = publicKeyStr
}

func generateRandomKeys() (string, string) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic("Failed to generate private key: " + err.Error())
	}

	publicKey := &privateKey.PublicKey

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	publicKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)

	privateKeyStr := base64.StdEncoding.EncodeToString(privateKeyBytes)
	publicKeyStr := base64.StdEncoding.EncodeToString(publicKeyBytes)
	return privateKeyStr, publicKeyStr
}
