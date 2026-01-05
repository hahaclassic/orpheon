package middleware

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/services/content-msv/internal/domain/entity"
)

const (
	headerUserID    = "X-User-Id"
	headerAccessLvl = "X-Access-Level"
)

var (
	ErrNoData           = errors.New("there is no user_id or access_lvl in the request")
	ErrInvalidUserID    = errors.New("invalid user id")
	ErrInvalidAccessLvl = errors.New("invalid access lvl")
)

func ClaimsRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.GetHeader(headerUserID)
		lvlStr := c.GetHeader(headerAccessLvl)

		slog.Info(
			"[USER-MSV] Middleware ClaimsRequired headers:",
			"X-User-Id", userIDStr,
			"X-Access-Level", lvlStr,
		)

		if userIDStr == "" || lvlStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrNoData.Error())
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrInvalidUserID.Error())
		}

		lvlInt, err := strconv.Atoi(lvlStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrInvalidAccessLvl.Error())
			return
		}

		lvl := entity.AccessLevel(lvlInt)
		if !lvl.IsValid() {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrInvalidAccessLvl.Error())
			return
		}

		claims := &entity.Claims{
			UserID:    userID,
			AccessLvl: lvl,
		}

		c.Set("claims", claims)
		c.Next()
	}
}

func ClaimsOptional() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.GetHeader(headerUserID)
		lvlStr := c.GetHeader(headerAccessLvl)

		if userIDStr == "" || lvlStr == "" {
			c.Next()
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.Next()
		}

		lvlInt, err := strconv.Atoi(lvlStr)
		if err != nil {
			c.Next()
		}

		lvl := entity.AccessLevel(lvlInt)
		if !lvl.IsValid() {
			c.Next()
		}

		claims := &entity.Claims{
			UserID:    userID,
			AccessLvl: lvl,
		}

		c.Set("claims", claims)
		c.Next()
	}
}
