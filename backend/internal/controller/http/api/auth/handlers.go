package auth_ctrl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hahaclassic/orpheon/backend/internal/config"
	"github.com/hahaclassic/orpheon/backend/internal/controller/http/middleware"
	"github.com/hahaclassic/orpheon/backend/internal/controller/http/utils"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/auth"
)

const (
	refreshCookieName = "refresh_token"
	accessCookieName  = "access_token"
)

type AuthController struct {
	service      auth.AuthService
	cookieConfig *config.CookieConfig
}

func NewAuthController(service auth.AuthService,
	cookieConfig *config.CookieConfig) *AuthController {
	return &AuthController{service: service, cookieConfig: cookieConfig}
}

func (ac *AuthController) RegisterRoutes(router *gin.RouterGroup) {
	authGroup := router.Group("/auth")
	authGroup.POST("/register", ac.register)
	authGroup.POST("/login", ac.login)
	authGroup.POST("/refresh", ac.refresh)
	authGroup.POST("/logout", ac.logout)

	passwordGroup := authGroup.Group("/password").Use(middleware.AuthMiddlewareRequired(ac.service))
	passwordGroup.POST("/update", ac.updatePassword)
}

func (ac *AuthController) register(c *gin.Context) {
	var creds entity.UserCredentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	tokens, err := ac.service.RegisterUser(c.Request.Context(), &creds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ac.packTokens(c, tokens)
}

func (ac *AuthController) login(c *gin.Context) {
	var creds entity.UserCredentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	tokens, err := ac.service.Login(c.Request.Context(), &creds)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ac.packTokens(c, tokens)
}

func (ac *AuthController) refresh(c *gin.Context) {
	refreshToken, err := c.Cookie(refreshCookieName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no refresh token"})
		return
	}

	tokens, err := ac.service.RefreshTokens(c.Request.Context(), refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	ac.packTokens(c, tokens)
}

func (ac *AuthController) logout(c *gin.Context) {
	refreshToken, err := c.Cookie(refreshCookieName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no refresh token"})
		return
	}

	if err := ac.service.Logout(c.Request.Context(), refreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (ac *AuthController) updatePassword(c *gin.Context) {
	claims := utils.GetClaims(c)
	if claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var pw entity.UserPasswords
	if err := c.ShouldBindJSON(&pw); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := ac.service.UpdatePassword(c.Request.Context(), claims.UserID, &pw); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (ac AuthController) packTokens(c *gin.Context, tokens *entity.AuthTokens) {
	c.SetCookie(
		refreshCookieName,
		tokens.Refresh,
		int(ac.cookieConfig.RefreshTTL.Seconds()),
		ac.cookieConfig.Path,   // path
		ac.cookieConfig.Domain, // domain ("" = current)
		ac.cookieConfig.Secure, // secure (set to false if testing locally w/o HTTPS)
		true,                   // httpOnly
	)

	c.SetCookie(
		accessCookieName,
		tokens.Access,
		int(ac.cookieConfig.AccessTTL.Seconds()),
		ac.cookieConfig.Path,     // path
		ac.cookieConfig.Domain,   // domain ("" = current)
		ac.cookieConfig.Secure,   // secure (set to false if testing locally w/o HTTPS)
		ac.cookieConfig.HttpOnly, // httpOnly
	)
}
