package auth

import (
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
	// keys, err := s.keyRepository.GetKeys()
	// if err != nil {
	// 	return "", err
	// }

	// privateKey := []byte(keys.(*Keys).PrivateKey)
	privateKey := []byte(helper.GetEnv("JWT_PRIVATE_KEY", ""))

	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(privateKey)
}

func (s *AuthService) TokenIsValid(tokenStr string) bool {
	privateKey := helper.GetEnv("JWT_PRIVATE_KEY", "")

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return privateKey, nil
	})

	if err != nil || !token.Valid {
		return false
	}

	return true
}
