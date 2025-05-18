package playlist_ctrl

import "github.com/gin-gonic/gin"

type PlaylistRouter struct {
	playlistMetaController      *PlaylistMetaController
	playlistTrackController     *PlaylistTrackController
	playlistFavoritesController *PlaylistFavoritesController
	playlistCoverController     *PlaylistCoverController
	authMiddleware              gin.HandlerFunc
}

func NewPlaylistRouter(playlistMetaController *PlaylistMetaController,
	playlistTrackController *PlaylistTrackController,
	playlistFavoritesController *PlaylistFavoritesController,
	playlistCoverController *PlaylistCoverController,
	authMiddleware gin.HandlerFunc) *PlaylistRouter {
	return &PlaylistRouter{
		playlistMetaController:      playlistMetaController,
		playlistTrackController:     playlistTrackController,
		playlistFavoritesController: playlistFavoritesController,
		playlistCoverController:     playlistCoverController,
		authMiddleware:              authMiddleware,
	}
}

func (r *PlaylistRouter) RegisterRoutes(router *gin.RouterGroup) {
	playlistGroup := router.Group("/playlists")
	playlistGroup.Use(r.authMiddleware)
	{
		playlistGroup.GET("/:id", r.playlistMetaController.GetPlaylist)
		playlistGroup.POST("", r.playlistMetaController.CreatePlaylist)
		playlistGroup.PUT("/:id", r.playlistMetaController.UpdatePlaylist)
		playlistGroup.DELETE("/:id", r.playlistMetaController.DeletePlaylist)
		playlistGroup.PATCH("/:id/privacy", r.playlistMetaController.UpdatePlaylistPrivacy)
	}

	coverGroup := playlistGroup.Group("/:id/cover")
	{
		coverGroup.GET("", r.playlistCoverController.GetCover)
		coverGroup.POST("", r.playlistCoverController.UploadCover)
		coverGroup.DELETE("", r.playlistCoverController.DeleteCover)
	}

	tracksGroup := playlistGroup.Group("/:id/tracks")
	{
		tracksGroup.GET("", r.playlistTrackController.GetPlaylistTracks)
		tracksGroup.POST("", r.playlistTrackController.AddTrackToPlaylist)
		tracksGroup.DELETE("/:track_id", r.playlistTrackController.DeleteTrackFromPlaylist)
		tracksGroup.PATCH("/:track_id/position", r.playlistTrackController.ChangeTrackPosition)
	}
}
