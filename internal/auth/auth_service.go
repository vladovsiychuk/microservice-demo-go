package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/vladovsiychuk/microservice-demo-go/pkg/helper"
)

type AuthService struct {
	keyRepository          KeyRepositoryI
	sessionTokenRepository SessionTokenRepositoryI
}

type AuthServiceI interface {
	GenerateJwtAndSessionTokens(email string) (string, SessionTokenI, error)
	RefreshJwtAndSessionTokens(sessionTokenId uuid.UUID) (string, SessionTokenI, error)
	TokenIsValid(tokenStr string) bool
}

func NewService(keyRepository KeyRepositoryI, sessionTokenRepository SessionTokenRepositoryI) *AuthService {
	return &AuthService{
		keyRepository,
		sessionTokenRepository,
	}
}

func (s *AuthService) Init() {
	initOauthProviders()

	keys := CreateKeys()

	if err := s.keyRepository.Update(keys); err != nil {
		panic("Can't save keys in the repository.")
	}

	s.startKeyRotation()
}

func (s *AuthService) startKeyRotation() {
	go func() {
		ticker := time.NewTicker(JWT_KEYS_DURATION)
		defer ticker.Stop()

		for range ticker.C {
			keys, _ := s.keyRepository.GetKeys()
			keys.Rotate()
			s.keyRepository.Update(keys)
		}
	}()
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

func (s *AuthService) GenerateJwtAndSessionTokens(email string) (string, SessionTokenI, error) {
	jwtTokenStr, err := s.generateJwtTokenStr(email)
	if err != nil {
		return "", nil, err
	}

	sessionToken := CreateSessionToken(email)
	return jwtTokenStr, sessionToken, err
}

func (s *AuthService) RefreshJwtAndSessionTokens(sessionTokenId uuid.UUID) (string, SessionTokenI, error) {
	currentSessionTokenI, err := s.sessionTokenRepository.FindById(sessionTokenId)
	currentSessionToken := currentSessionTokenI.(*SessionToken)
	if err != nil {
		// Something went wrong, the session token was most likely stolen.
		// Because it's quite strange that frontend has sent a token that was deleted
		// if it was deleted means that it was used.
		// The best thing to do here is to delete all current user sessions (all session tokens by user id/email)
		// from the repository.
		// And redirect the user to the login page (return error)
		return "", nil, err
	}

	newSessionToken := CreateSessionToken(currentSessionToken.Email)
	newJwtTokenStr, err := s.generateJwtTokenStr(currentSessionToken.Email)
	if err != nil {
		return "", nil, err
	}

	if err := s.sessionTokenRepository.Delete(currentSessionTokenI); err != nil {
		return "", nil, err
	}

	if err := s.sessionTokenRepository.Create(newSessionToken); err != nil {
		return "", nil, err
	}

	return newJwtTokenStr, newSessionToken, nil
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

func (s *AuthService) generateJwtTokenStr(email string) (string, error) {
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
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
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
