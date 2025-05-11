package playlist

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
)

type PlaylistController struct {
	playlistService playlist.PlaylistService
}

func New(playlistService playlist.PlaylistService) *PlaylistController {
	return &PlaylistController{
		playlistService: playlistService,
	}
}

func (c *PlaylistController) RegisterRoutes(router *gin.Engine) {
	playlists := router.Group("/api/v1/playlists")
	{
		playlists.GET("/:id", c.GetPlaylist)
		playlists.POST("", c.CreatePlaylist)
		playlists.PUT("/:id", c.UpdatePlaylist)
		playlists.DELETE("/:id", c.DeletePlaylist)
		playlists.GET("/user/:userId", c.GetUserPlaylists)
	}
}

// GetPlaylist godoc
// @Summary Get playlist by ID
// @Description Get playlist details by its ID
// @Tags playlists
// @Produce json
// @Param id path string true "Playlist ID"
// @Success 200 {object} entity.Playlist
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/playlists/{id} [get]
func (c *PlaylistController) GetPlaylist(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID"})
		return
	}

	playlist, err := c.playlistService.GetPlaylist(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Playlist not found"})
		return
	}

	ctx.JSON(http.StatusOK, playlist)
}

// CreatePlaylist godoc
// @Summary Create a new playlist
// @Description Create a new playlist with the provided details
// @Tags playlists
// @Accept json
// @Produce json
// @Param playlist body entity.Playlist true "Playlist object"
// @Success 201 {object} entity.Playlist
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/playlists [post]
func (c *PlaylistController) CreatePlaylist(ctx *gin.Context) {
	var playlist entity.Playlist
	if err := ctx.ShouldBindJSON(&playlist); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdPlaylist, err := c.playlistService.CreatePlaylist(ctx.Request.Context(), &playlist)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Failed to create playlist"})
		return
	}

	ctx.JSON(http.StatusCreated, createdPlaylist)
}

// UpdatePlaylist godoc
// @Summary Update a playlist
// @Description Update an existing playlist with new details
// @Tags playlists
// @Accept json
// @Produce json
// @Param id path string true "Playlist ID"
// @Param playlist body entity.Playlist true "Updated playlist object"
// @Success 200 {object} entity.Playlist
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/playlists/{id} [put]
func (c *PlaylistController) UpdatePlaylist(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID"})
		return
	}

	var playlist entity.Playlist
	if err := ctx.ShouldBindJSON(&playlist); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	playlist.ID = id
	updatedPlaylist, err := c.playlistService.UpdatePlaylist(ctx.Request.Context(), &playlist)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Failed to update playlist"})
		return
	}

	ctx.JSON(http.StatusOK, updatedPlaylist)
}

// DeletePlaylist godoc
// @Summary Delete a playlist
// @Description Delete a playlist by its ID
// @Tags playlists
// @Produce json
// @Param id path string true "Playlist ID"
// @Success 204 "No Content"
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/playlists/{id} [delete]
func (c *PlaylistController) DeletePlaylist(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID"})
		return
	}

	err = c.playlistService.DeletePlaylist(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Failed to delete playlist"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// GetUserPlaylists godoc
// @Summary Get user's playlists
// @Description Get all playlists created by a specific user
// @Tags playlists
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {array} entity.Playlist
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/playlists/user/{userId} [get]
func (c *PlaylistController) GetUserPlaylists(ctx *gin.Context) {
	userID, err := uuid.Parse(ctx.Param("userId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	playlists, err := c.playlistService.GetUserPlaylists(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user playlists"})
		return
	}

	ctx.JSON(http.StatusOK, playlists)
}
