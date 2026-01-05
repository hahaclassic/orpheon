package playlist_router

import "github.com/gin-gonic/gin"

type PlaylistMetaController interface {
	GetPlaylist(c *gin.Context)
	GetPlaylists(c *gin.Context)
	CreatePlaylist(c *gin.Context)
	UpdatePlaylist(c *gin.Context)
	DeletePlaylist(c *gin.Context)
	UpdatePlaylistPrivacy(c *gin.Context)
}

type PlaylistCoverController interface {
	GetCover(c *gin.Context)
	UploadCover(c *gin.Context)
	DeleteCover(c *gin.Context)
}

type PlaylistTrackController interface {
	GetPlaylistTracks(c *gin.Context)
	AddTrackToPlaylist(c *gin.Context)
	DeleteTrackFromPlaylist(c *gin.Context)
	ChangeTrackPosition(c *gin.Context)
}

type PlaylistFavoritesController interface {
	GetFavoritePlaylists(c *gin.Context)
	AddToFavorites(c *gin.Context)
	RemoveFromFavorites(c *gin.Context)
}

type PlaylistRouter struct {
	playlistMetaController      PlaylistMetaController
	playlistTrackController     PlaylistTrackController
	playlistCoverController     PlaylistCoverController
	playlistFavoritesController PlaylistFavoritesController
	authMiddleware              gin.HandlerFunc
}

func NewPlaylistRouter(
	playlistMetaController PlaylistMetaController,
	playlistTrackController PlaylistTrackController,
	playlistCoverController PlaylistCoverController,
	playlistFavoritesController PlaylistFavoritesController,
	authMiddleware gin.HandlerFunc,
) *PlaylistRouter {
	return &PlaylistRouter{
		playlistMetaController:      playlistMetaController,
		playlistTrackController:     playlistTrackController,
		playlistCoverController:     playlistCoverController,
		playlistFavoritesController: playlistFavoritesController,
		authMiddleware:              authMiddleware,
	}
}

func (r *PlaylistRouter) RegisterRoutes(router *gin.RouterGroup) {
	playlistGroup := router.Group("/playlists")
	playlistGroup.Use(r.authMiddleware)
	{
		playlistGroup.GET("", r.playlistMetaController.GetPlaylists)
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

	router.GET("/favorites", r.playlistFavoritesController.GetFavoritePlaylists)
	router.POST("/favorites/:playlist_id", r.playlistFavoritesController.AddToFavorites)
	router.DELETE("/favorites/:playlist_id", r.playlistFavoritesController.RemoveFromFavorites)
}
