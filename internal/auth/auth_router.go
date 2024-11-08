package auth

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

type AuthRouter struct {
	service AuthServiceI
}

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
	}
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

	token, err := h.service.GenerateJwt(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,                    // Can't be accessed by JavaScript
		Secure:   true,                    // Use Secure if using HTTPS
		SameSite: http.SameSiteStrictMode, // Optional, for CSRF protection
		MaxAge:   3600,                    // Token expiry (1 hour)
	})

	c.Redirect(http.StatusFound, "http://localhost:3000/dashboard")
}
