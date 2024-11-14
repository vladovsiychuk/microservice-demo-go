package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"time"
)

type Keys struct {
	PrivateKey         string
	PublicKey          string
	SecondaryPublicKey string
}

type KeysI interface {
	Rotate()
	GetPrivateKey() (*rsa.PrivateKey, error)
	GetPublicKey() (*rsa.PublicKey, error)
	GetSecondaryPulicKey() (*rsa.PublicKey, error)
}

var JWT_KEYS_DURATION = 30 * time.Second

var CreateKeys = func() KeysI {
	privateKeyStr, publicKeyStr := generateRandomKeyStrs()

	return &Keys{
		privateKeyStr,
		publicKeyStr,
		publicKeyStr,
	}
}

func (k *Keys) Rotate() {
	privateKeyStr, publicKeyStr := generateRandomKeyStrs()

	k.PrivateKey = privateKeyStr
	k.SecondaryPublicKey = k.PublicKey
	k.PublicKey = publicKeyStr
}

func (k *Keys) GetPrivateKey() (*rsa.PrivateKey, error) {
	return decodeBase64PrivateKey(k.PrivateKey)
}

func (k *Keys) GetPublicKey() (*rsa.PublicKey, error) {
	return decodeBase64PublicKey(k.PublicKey)
}

func (k *Keys) GetSecondaryPulicKey() (*rsa.PublicKey, error) {
	return decodeBase64PublicKey(k.SecondaryPublicKey)
}

func generateRandomKeyStrs() (string, string) {
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

func decodeBase64PrivateKey(encodedKey string) (*rsa.PrivateKey, error) {
	decodedKey, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		return nil, err
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(decodedKey)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func decodeBase64PublicKey(encodedKey string) (*rsa.PublicKey, error) {
	decodedKey, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		return nil, err
	}
	publicKey, err := x509.ParsePKCS1PublicKey(decodedKey)
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}
