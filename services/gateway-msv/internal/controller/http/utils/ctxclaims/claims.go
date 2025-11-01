package ctxclaims

import (
	"github.com/gin-gonic/gin"
	"github.com/hahaclassic/orpheon/services/gateway-msv/internal/controller/http/dto"
)

func GetClaims(c *gin.Context) *dto.Claims {
	claims, exists := c.Get("claims")
	if !exists {
		return nil
	}

	parsedClaims, ok := claims.(*dto.Claims)
	if !ok {
		return nil
	}

	return parsedClaims
}
