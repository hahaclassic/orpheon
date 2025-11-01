package auth_ctrl

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hahaclassic/orpheon/services/gateway-msv/internal/controller/http/dto"
	"github.com/hahaclassic/orpheon/services/gateway-msv/internal/controller/http/utils/cookie"
	"github.com/hahaclassic/orpheon/services/gateway-msv/internal/controller/http/utils/ctxclaims"
)

type AuthService interface {
	RegisterUser(ctx context.Context, request *dto.RegisterRequest) (*dto.AuthTokens, error)
	Login(ctx context.Context, creds *dto.UserCredentials) (*dto.AuthTokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*dto.AuthTokens, error)
	Logout(ctx context.Context, refreshToken string) error
	GetClaims(ctx context.Context, accessToken string) (*dto.Claims, error)
	UpdatePassword(ctx context.Context, claims *dto.Claims, passwords *dto.UserPasswords) error
}

type AuthController struct {
	service            AuthService
	cookieTokensSetter *cookie.CookieTokensSetter
	authMiddleware     gin.HandlerFunc
}

func NewAuthController(service AuthService,
	cookieTokensSetter *cookie.CookieTokensSetter,
	authMiddleware gin.HandlerFunc) *AuthController {

	return &AuthController{
		service:            service,
		cookieTokensSetter: cookieTokensSetter,
		authMiddleware:     authMiddleware,
	}
}

func (ac *AuthController) RegisterRoutes(router *gin.RouterGroup) {
	authGroup := router.Group("/auth")
	authGroup.POST("/register", ac.register)
	authGroup.POST("/login", ac.login)
	authGroup.POST("/refresh", ac.refresh)
	authGroup.POST("/logout", ac.logout)

	passwordGroup := authGroup.Group("/password").Use(ac.authMiddleware)
	passwordGroup.POST("/update", ac.updatePassword)
}

func (ac *AuthController) register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})

		return
	}

	tokens, err := ac.service.RegisterUser(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ac.cookieTokensSetter.SetAll(c, tokens)
}

func (ac *AuthController) login(c *gin.Context) {
	var creds dto.UserCredentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})

		return
	}

	tokens, err := ac.service.Login(c.Request.Context(), &creds)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})

		return
	}

	ac.cookieTokensSetter.SetAll(c, tokens)
}

func (ac *AuthController) refresh(c *gin.Context) {
	refreshToken, err := c.Cookie(cookie.RefreshCookieName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no refresh token"})
		return
	}

	tokens, err := ac.service.RefreshTokens(c.Request.Context(), refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	slog.Info("Tokens refreshed")
	ac.cookieTokensSetter.SetAll(c, tokens)
}

func (ac *AuthController) logout(c *gin.Context) {
	refreshToken, err := c.Cookie(cookie.RefreshCookieName)
	if err == nil {
		if err := ac.service.Logout(c.Request.Context(), refreshToken); err != nil {
			slog.Error("failed to logout", "error", err)
			c.JSON(http.StatusOK, gin.H{"message": err.Error()})
		}

		ac.cookieTokensSetter.ResetRefresh(c)
	}

	ac.cookieTokensSetter.ResetAccess(c)

	c.Status(http.StatusOK)
}

func (ac *AuthController) updatePassword(c *gin.Context) {
	claims := ctxclaims.GetClaims(c)
	if claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var passwords dto.UserPasswords
	if err := c.ShouldBindJSON(&passwords); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := ac.service.UpdatePassword(c.Request.Context(), claims, &passwords); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
