package playlist

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
)

type PlaylistCoverController struct {
	coverService playlist.CoverService
}

func NewCoverController(coverService playlist.CoverService) *PlaylistCoverController {
	return &PlaylistCoverController{
		coverService: coverService,
	}
}

func (c *PlaylistCoverController) RegisterRoutes(router *gin.RouterGroup) {
	playlists := router.Group("/playlists")
	{
		// Public routes
		playlists.GET("/:id/cover", c.GetCover)

		// Protected routes
		protected := playlists.Group("")
		protected.Use(middleware.Auth())
		{
			protected.POST("/:id/cover", c.UploadCover)
			protected.DELETE("/:id/cover", c.DeleteCover)
		}
	}
}

// GetCover godoc
// @Summary Get playlist cover
// @Description Get the cover image for a playlist
// @Tags playlists
// @Produce image/*
// @Param id path string true "Playlist ID"
// @Success 200 {file} binary
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/playlists/{id}/cover [get]
func (c *PlaylistCoverController) GetCover(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID"})
		return
	}

	cover, err := c.coverService.GetCover(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Cover not found"})
		return
	}

	ctx.DataFromReader(http.StatusOK, int64(len(cover)), "image/jpeg", cover, nil)
}

// UploadCover godoc
// @Summary Upload playlist cover
// @Description Upload a cover image for a playlist
// @Tags playlists
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Playlist ID"
// @Param cover formData file true "Cover image file"
// @Security BearerAuth
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/playlists/{id}/cover [post]
func (c *PlaylistCoverController) UploadCover(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID"})
		return
	}

	file, err := ctx.FormFile("cover")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No cover file provided"})
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer openedFile.Close()

	err = c.coverService.UploadCover(ctx.Request.Context(), id, openedFile)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Failed to upload cover"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Cover uploaded successfully"})
}

// DeleteCover godoc
// @Summary Delete playlist cover
// @Description Delete the cover image for a playlist
// @Tags playlists
// @Produce json
// @Param id path string true "Playlist ID"
// @Security BearerAuth
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/playlists/{id}/cover [delete]
func (c *PlaylistCoverController) DeleteCover(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID"})
		return
	}

	err = c.coverService.DeleteCover(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Failed to delete cover"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Cover deleted successfully"})
}
