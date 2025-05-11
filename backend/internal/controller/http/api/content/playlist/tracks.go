package playlist

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
)

type PlaylistTracksController struct {
	tracksService playlist.TracksService
}

func NewTracksController(tracksService playlist.TracksService) *PlaylistTracksController {
	return &PlaylistTracksController{
		tracksService: tracksService,
	}
}

func (c *PlaylistTracksController) RegisterRoutes(router *gin.RouterGroup) {
	playlists := router.Group("/playlists")
	{
		// Public routes
		playlists.GET("/:id/tracks", c.GetPlaylistTracks)

		// Protected routes
		protected := playlists.Group("")
		protected.Use(middleware.Auth())
		{
			protected.POST("/:id/tracks", c.AddTrackToPlaylist)
			protected.DELETE("/:id/tracks/:trackId", c.RemoveTrackFromPlaylist)
		}
	}
}

// GetPlaylistTracks godoc
// @Summary Get playlist tracks
// @Description Get all tracks in a playlist
// @Tags playlists
// @Produce json
// @Param id path string true "Playlist ID"
// @Success 200 {array} entity.Track
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/playlists/{id}/tracks [get]
func (c *PlaylistTracksController) GetPlaylistTracks(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID"})
		return
	}

	tracks, err := c.tracksService.GetPlaylistTracks(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Failed to get playlist tracks"})
		return
	}

	ctx.JSON(http.StatusOK, tracks)
}

// AddTrackToPlaylist godoc
// @Summary Add track to playlist
// @Description Add a track to an existing playlist
// @Tags playlists
// @Accept json
// @Produce json
// @Param id path string true "Playlist ID"
// @Param trackId body string true "Track ID"
// @Security BearerAuth
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/playlists/{id}/tracks [post]
func (c *PlaylistTracksController) AddTrackToPlaylist(ctx *gin.Context) {
	playlistID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID"})
		return
	}

	var request struct {
		TrackID string `json:"trackId" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trackID, err := uuid.Parse(request.TrackID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid track ID"})
		return
	}

	err = c.tracksService.AddTrackToPlaylist(ctx.Request.Context(), playlistID, trackID)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Failed to add track to playlist"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Track added to playlist"})
}

// RemoveTrackFromPlaylist godoc
// @Summary Remove track from playlist
// @Description Remove a track from a playlist
// @Tags playlists
// @Produce json
// @Param id path string true "Playlist ID"
// @Param trackId path string true "Track ID"
// @Security BearerAuth
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/playlists/{id}/tracks/{trackId} [delete]
func (c *PlaylistTracksController) RemoveTrackFromPlaylist(ctx *gin.Context) {
	playlistID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID"})
		return
	}

	trackID, err := uuid.Parse(ctx.Param("trackId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid track ID"})
		return
	}

	err = c.tracksService.RemoveTrackFromPlaylist(ctx.Request.Context(), playlistID, trackID)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Failed to remove track from playlist"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Track removed from playlist"})
}
