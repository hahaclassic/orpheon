package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hahaclassic/orpheon/services/gateway-msv/internal/controller/http/dto"
	"github.com/hahaclassic/orpheon/services/gateway-msv/internal/controller/http/utils/cookie"
)

var (
	ErrNoAccessToken = errors.New("no access token")
	ErrNoTokens      = errors.New("no tokens")
	ErrExpiredToken  = errors.New("token expired")
	ErrInvalidToken  = errors.New("invalid token")
)

type AuthService interface {
	RefreshTokens(ctx context.Context, refreshToken string) (*dto.AuthTokens, error)
	GetClaims(ctx context.Context, accessToken string) (*dto.Claims, error)
}

type AuthMiddleware struct {
	authService        AuthService
	cookieTokensSetter *cookie.CookieTokensSetter
}

func NewAuthMiddleware(authService AuthService, cookieTokensSetter *cookie.CookieTokensSetter) *AuthMiddleware {
	return &AuthMiddleware{
		authService:        authService,
		cookieTokensSetter: cookieTokensSetter,
	}
}

func (a *AuthMiddleware) Required() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := a.setClaims(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})

			return
		}

		c.Next()
	}
}

func (a *AuthMiddleware) Optional() gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = a.setClaims(c)
		c.Next()
	}
}

func (a *AuthMiddleware) setClaims(c *gin.Context) error {
	accessToken, err := c.Cookie(cookie.AccessCookieName)
	if err != nil || accessToken == "" {
		tokens, err := a.refresh(c)
		if err != nil {
			return err
		}
		accessToken = tokens.Access
	}

	claims, err := a.authService.GetClaims(c, accessToken)
	if err != nil {
		tokens, err := a.refresh(c)
		if err != nil {
			return err
		}
		claims, err = a.authService.GetClaims(c, tokens.Access)
		if err != nil {
			return err
		}
	}

	c.Set("claims", claims)

	return nil
}

func (a *AuthMiddleware) refresh(c *gin.Context) (*dto.AuthTokens, error) {
	refreshToken, err := c.Cookie(cookie.RefreshCookieName)
	if err != nil || refreshToken == "" {
		return nil, ErrNoTokens
	}

	tokens, err := a.authService.RefreshTokens(c.Request.Context(), refreshToken)
	if err != nil {
		return nil, err
	}

	a.cookieTokensSetter.SetAll(c, tokens)

	return tokens, nil
}
