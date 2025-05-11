package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/user"
)

type UserController struct {
	userService user.UserService
}

func New(userService user.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (c *UserController) RegisterRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.GET("/:id", c.GetUser)
		users.PUT("/:id", c.UpdateUser)
		users.DELETE("/:id", c.DeleteUser)
	}
}

// // CreateUser godoc
// // @Summary Create a new user
// // @Description Create a new user with the provided information
// // @Tags users
// // @Accept json
// // @Produce json
// // @Param user body entity.UserInfo true "User information"
// // @Success 201 {object} uuid.UUID
// // @Failure 400 {object} gin.H
// // @Failure 500 {object} gin.H
// // @Router /api/v1/users [post]
// func (c *UserController) CreateUser(ctx *gin.Context) {
// 	var userInfo entity.UserInfo
// 	if err := ctx.ShouldBindJSON(&userInfo); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
// 		return
// 	}

// 	userID, err := c.userService.CreateUser(ctx.Request.Context(), &userInfo)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
// 		return
// 	}

// 	ctx.JSON(http.StatusCreated, userID)
// }

// GetUser godoc
// @Summary Get user information
// @Description Get user information by ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} entity.UserInfo
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/users/{id} [get]
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

// UpdateUser godoc
// @Summary Update user information
// @Description Update user information by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body entity.UserInfo true "User information"
// @Security BearerAuth
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/users/{id} [put]
func (c *UserController) UpdateUser(ctx *gin.Context) {
	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var userInfo entity.UserInfo
	if err := ctx.ShouldBindJSON(&userInfo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userInfo.ID = userID

	if err := c.userService.UpdateUser(ctx.Request.Context(), claims.(*entity.Claims), &userInfo); err != nil {
		if err == user.ErrUpdateUser {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete user by ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Security BearerAuth
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/users/{id} [delete]
func (c *UserController) DeleteUser(ctx *gin.Context) {
	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := c.userService.DeleteUser(ctx.Request.Context(), claims.(*entity.Claims), userID); err != nil {
		if err == user.ErrDeleteUser {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
