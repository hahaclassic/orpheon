package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReadOnly(readOnly bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if readOnly && c.Request.Method != http.MethodGet {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Server is in read-only mode",
				"message": "Only GET requests are allowed",
				"method":  c.Request.Method,
				"allowed": []string{"GET"},
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
