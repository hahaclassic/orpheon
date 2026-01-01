package middleware

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/services/auth-msv/internal/domain/entity"
)

const (
	HeaderUserID    = "X-User-Id"
	HeaderAccessLvl = "X-Access-Level"
)

var (
	ErrNoData           = errors.New("there is no user_id or access_lvl in the request")
	ErrInvalidUserID    = errors.New("invalid user id")
	ErrInvalidAccessLvl = errors.New("invalid access lvl")
)

func ClaimsRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.GetHeader(HeaderUserID)
		lvlStr := c.GetHeader(HeaderAccessLvl)

		if userIDStr == "" || lvlStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrNoData)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrInvalidUserID)
		}

		lvlInt, err := strconv.Atoi(lvlStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrInvalidAccessLvl)
			return
		}

		lvl := entity.AccessLevel(lvlInt)
		if !lvl.IsValid() {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrInvalidAccessLvl)
			return
		}

		claims := entity.Claims{
			UserID:    userID,
			AccessLvl: lvl,
		}

		c.Set("claims", claims)
		c.Next()
	}
}

func ClaimsOptional() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.GetHeader("X-User-Id")
		lvlStr := c.GetHeader("X-Access-Level")

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

		claims := entity.Claims{
			UserID:    userID,
			AccessLvl: lvl,
		}

		c.Set("claims", claims)
		c.Next()
	}
}
