package middleware

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	jwttokens "github.com/hahaclassic/orpheon/backend/internal/adapters/tokens/jwt"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/auth"
)

func AuthMiddleware(authService auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.Info("Middleware triggered")
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		fmt.Println("Authorization header:", authHeader)

		token := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := authService.GetClaims(c, token)
		if err != nil {
			var h gin.H
			if errors.Is(err, jwttokens.ErrExpired) {
				h = gin.H{"error": "token expired"}
			} else {
				h = gin.H{"error": "invalid token"}
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, h)
			return
		}

		fmt.Println("Claims:", claims)

		c.Set("claims", claims)
		c.Next()
	}
}
