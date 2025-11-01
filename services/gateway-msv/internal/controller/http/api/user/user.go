package user_ctrl

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/pkg/commonerr"
	"github.com/hahaclassic/orpheon/services/gateway-msv/internal/controller/http/dto"
	"github.com/hahaclassic/orpheon/services/gateway-msv/internal/controller/http/utils/ctxclaims"
)

type UserService interface {
	GetUser(ctx context.Context, userID uuid.UUID) (*dto.User, error)
	UpdateUser(ctx context.Context, claims *dto.Claims, user *dto.EditUserRequest) error
	DeleteUser(ctx context.Context, claims *dto.Claims, userID uuid.UUID) error
}

type UserController struct {
	userService UserService
}

func NewUserController(userService UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (c *UserController) GetMe(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userInfo, err := c.userService.GetUser(ctx.Request.Context(), claims.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	ctx.JSON(http.StatusOK, userInfo)
}

func (c *UserController) GetUser(ctx *gin.Context) {
	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userInfo, err := c.userService.GetUser(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	if userInfo == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, userInfo)
}

func (c *UserController) UpdateMe(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user dto.EditUserRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := c.userService.UpdateUser(ctx.Request.Context(), claims, &user); err != nil {
		if errors.Is(err, commonerr.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (c *UserController) DeleteUser(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := c.userService.DeleteUser(ctx.Request.Context(), claims, userID); err != nil {
		if errors.Is(err, commonerr.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
