package http_ctrl

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	commonerr "github.com/hahaclassic/orpheon/pkg/commonerr"
	"github.com/hahaclassic/orpheon/services/user-msv/internal/domain/entity"
	"github.com/hahaclassic/orpheon/services/user-msv/internal/domain/usecase"
	ctxclaims "github.com/hahaclassic/orpheon/services/user-msv/internal/providers/http/ctxclaims"
)

type UserController struct {
	userService      usecase.UserService
	claimsMiddleware gin.HandlerFunc
}

func NewUserController(userService usecase.UserService,
	claimsMiddleware gin.HandlerFunc) *UserController {
	return &UserController{
		userService:      userService,
		claimsMiddleware: claimsMiddleware,
	}
}

func (c *UserController) RegisterRoutes(g *gin.RouterGroup) {
	usersGroup := g.Group("/users")
	usersGroup.Use(c.claimsMiddleware)
	usersGroup.GET("/:id", c.GetUser)
	usersGroup.DELETE("/:id", c.DeleteUser)

	meGroup := usersGroup.Group("/me")
	meGroup.GET("/", c.GetMe)
	meGroup.PUT("/", c.UpdateMe)
}

func (c *UserController) GetMe(ctx *gin.Context) {
	var err error

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

	var userInfo entity.User
	if err := ctx.ShouldBindJSON(&userInfo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userInfo.ID = claims.UserID

	if err := c.userService.UpdateUser(ctx.Request.Context(), claims, &userInfo); err != nil {
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
