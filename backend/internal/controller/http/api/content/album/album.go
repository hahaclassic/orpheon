package album

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/controller/http/utils"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/album"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
)

type AlbumController struct {
	albumService   album.AlbumService
	authMiddleware gin.HandlerFunc
}

func New(albumService album.AlbumService, authMiddleware gin.HandlerFunc) *AlbumController {
	return &AlbumController{
		albumService:   albumService,
		authMiddleware: authMiddleware,
	}
}

func (c *AlbumController) RegisterRoutes(router *gin.RouterGroup) {
	albums := router.Group("/albums")
	{
		albums.GET("/:id", c.GetAlbum)

		protected := albums.Group("")
		protected.Use(c.authMiddleware)
		{
			protected.POST("", c.CreateAlbum)
			protected.PUT("/:id", c.UpdateAlbum)
			protected.DELETE("/:id", c.DeleteAlbum)
		}
	}
}

func (c *AlbumController) GetAlbum(ctx *gin.Context) {
	albumID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
		return
	}

	album, err := c.albumService.GetAlbum(ctx.Request.Context(), albumID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get album"})
		return
	}

	if album == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
		return
	}

	ctx.JSON(http.StatusOK, album)
}

func (c *AlbumController) CreateAlbum(ctx *gin.Context) {
	claims := utils.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var album entity.AlbumMeta
	if err := ctx.ShouldBindJSON(&album); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := c.albumService.CreateAlbum(ctx.Request.Context(), claims, &album); err != nil {
		if errors.Is(err, commonerr.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create album"})
		return
	}

	ctx.JSON(http.StatusCreated, album.ID)
}

func (c *AlbumController) UpdateAlbum(ctx *gin.Context) {
	claims := utils.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	albumID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
		return
	}

	var album entity.AlbumMeta
	if err := ctx.ShouldBindJSON(&album); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	album.ID = albumID

	if err := c.albumService.UpdateAlbum(ctx.Request.Context(), claims, &album); err != nil {
		if errors.Is(err, commonerr.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update album"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Album updated successfully"})
}

func (c *AlbumController) DeleteAlbum(ctx *gin.Context) {
	claims := utils.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	albumID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
		return
	}

	if err := c.albumService.DeleteAlbum(ctx.Request.Context(), claims, albumID); err != nil {
		if errors.Is(err, commonerr.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete album"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Album deleted successfully"})
}
