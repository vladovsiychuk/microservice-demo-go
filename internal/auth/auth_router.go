package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/markbates/goth/gothic"
)

type AuthRouter struct {
	service AuthServiceI
}

var AUTH_TOKEN = "auth_token"
var SESSION_TOKEN_ID = "session_token_id"

func NewRouter(service AuthServiceI) *AuthRouter {
	return &AuthRouter{
		service: service,
	}
}

func (h *AuthRouter) RegisterRoutes(r *gin.Engine) {
	postGroup := r.Group("auth")
	{
		postGroup.GET("/login", h.login)
		postGroup.GET("/callback", h.callback)
		postGroup.GET("/logout", h.logout)
		postGroup.GET("/refresh", h.refresh)
	}
}

func (h *AuthRouter) refresh(c *gin.Context) {
	sessionTokenID, err := c.Cookie(SESSION_TOKEN_ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Session token missing"})
		return
	}

	sessionTokenIdStr, err := uuid.Parse(sessionTokenID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "session token UUID is not valid"})
		return
	}

	jwtTokenStr, sessionToken, err := h.service.RefreshJwtAndSessionTokens(sessionTokenIdStr)
	if err != nil {
		clearCookie(c, AUTH_TOKEN)
		clearCookie(c, SESSION_TOKEN_ID)
		c.Redirect(http.StatusNotFound, "http://localhost:3000")
		return
	}

	setCookie(c, SESSION_TOKEN_ID, sessionToken.(*SessionToken).Id.String(), SESSION_TOKEN_DURATION)
	setCookie(c, AUTH_TOKEN, jwtTokenStr, JWT_KEYS_DURATION)

	c.Status(http.StatusOK)
}

func (h *AuthRouter) login(c *gin.Context) {
	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "provider", "google"))
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func (h *AuthRouter) callback(c *gin.Context) {
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	jwtTokenStr, sessionToken, err := h.service.GenerateJwtAndSessionTokens(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	setCookie(c, AUTH_TOKEN, jwtTokenStr, JWT_KEYS_DURATION)
	setCookie(c, SESSION_TOKEN_ID, sessionToken.(*SessionToken).Id.String(), SESSION_TOKEN_DURATION)

	c.Redirect(http.StatusFound, "http://localhost:3000/dashboard")
}

func (h *AuthRouter) logout(c *gin.Context) {
	err := gothic.Logout(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Logout failed"})
		return
	}

	clearCookie(c, AUTH_TOKEN)
	c.Redirect(http.StatusFound, "http://localhost:3000")
}

func setCookie(c *gin.Context, name, value string, duration time.Duration) {
	var httpOnly bool
	if name == AUTH_TOKEN {
		httpOnly = false
	} else {
		httpOnly = true
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: httpOnly,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(duration),
	})
}

func clearCookie(c *gin.Context, name string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
	})
}
