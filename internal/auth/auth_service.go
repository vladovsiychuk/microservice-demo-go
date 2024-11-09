package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/vladovsiychuk/microservice-demo-go/pkg/helper"
)

type AuthService struct {
	keyRepository KeyRepositoryI
}

type AuthServiceI interface {
	GenerateJwt(email string) (string, error)
	TokenIsValid(tokenStr string) bool
}

func NewService(keyRepository KeyRepositoryI) *AuthService {
	return &AuthService{
		keyRepository,
	}
}

func (s *AuthService) Init() {
	initOauthProviders()

	keys := CreateKeys()

	if err := s.keyRepository.Update(keys); err != nil {
		panic("Can't save keys in the repository.")
	}
}

func initOauthProviders() {
	sessionSecret := helper.GetEnv("SESSION_SECRET", "")
	googleClientKey := helper.GetEnv("GOOGLE_OAUTH_CLIENT_KEY", "")
	googleSecret := helper.GetEnv("GOOGLE_OAUTH_SECRET", "")

	gothic.Store = sessions.NewCookieStore([]byte(sessionSecret))

	goth.UseProviders(
		google.New(
			googleClientKey,
			googleSecret,
			"http://localhost:8080/auth/callback",
			"email", "profile",
		),
	)
}

func (s *AuthService) GenerateJwt(email string) (string, error) {
	keys, err := s.keyRepository.GetKeys()
	if err != nil {
		return "", err
	}

	privateKey, err := decodeBase64PrivateKey(keys.(*Keys).PrivateKey)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

func (s *AuthService) TokenIsValid(tokenStr string) bool {
	keys, err := s.keyRepository.GetKeys()
	if err != nil {
		fmt.Printf("Failed to load keys from repository.")
		return false
	}

	publicKey, err := decodeBase64PublicKey(keys.(*Keys).PublicKey)
	if err != nil {
		fmt.Printf("Failed to decode public key.")
		return false
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil || !token.Valid {
		fmt.Printf("Token is not valid.")
		return false
	}

	return true
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
