package playlist

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
)

type PlaylistFavoritesController struct {
	favoritesService playlist.FavoritesService
}

func NewFavoritesController(favoritesService playlist.FavoritesService) *PlaylistFavoritesController {
	return &PlaylistFavoritesController{
		favoritesService: favoritesService,
	}
}

func (c *PlaylistFavoritesController) RegisterRoutes(router *gin.RouterGroup) {
	playlists := router.Group("/playlists")
	{
		// Protected routes
		protected := playlists.Group("")
		protected.Use(middleware.Auth())
		{
			protected.GET("/favorites", c.GetFavoritePlaylists)
			protected.POST("/:id/favorite", c.AddToFavorites)
			protected.DELETE("/:id/favorite", c.RemoveFromFavorites)
		}
	}
}

// GetFavoritePlaylists godoc
// @Summary Get favorite playlists
// @Description Get all favorite playlists for the current user
// @Tags playlists
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entity.Playlist
// @Failure 500 {object} gin.H
// @Router /api/v1/playlists/favorites [get]
func (c *PlaylistFavoritesController) GetFavoritePlaylists(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	playlists, err := c.favoritesService.GetFavoritePlaylists(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get favorite playlists"})
		return
	}

	ctx.JSON(http.StatusOK, playlists)
}

// AddToFavorites godoc
// @Summary Add playlist to favorites
// @Description Add a playlist to user's favorites
// @Tags playlists
// @Produce json
// @Param id path string true "Playlist ID"
// @Security BearerAuth
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/playlists/{id}/favorite [post]
func (c *PlaylistFavoritesController) AddToFavorites(ctx *gin.Context) {
	playlistID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID"})
		return
	}

	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err = c.favoritesService.AddToFavorites(ctx.Request.Context(), userID, playlistID)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Failed to add playlist to favorites"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Playlist added to favorites"})
}

// RemoveFromFavorites godoc
// @Summary Remove playlist from favorites
// @Description Remove a playlist from user's favorites
// @Tags playlists
// @Produce json
// @Param id path string true "Playlist ID"
// @Security BearerAuth
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/playlists/{id}/favorite [delete]
func (c *PlaylistFavoritesController) RemoveFromFavorites(ctx *gin.Context) {
	playlistID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID"})
		return
	}

	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err = c.favoritesService.RemoveFromFavorites(ctx.Request.Context(), userID, playlistID)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Failed to remove playlist from favorites"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Playlist removed from favorites"})
}
