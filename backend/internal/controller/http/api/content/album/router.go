package album_ctrl

import (
	"github.com/gin-gonic/gin"
)

type AlbumRouter struct {
	albumController      *AlbumController
	albumCoverController *AlbumCoverController
	albumTrackController *AlbumTrackController
	authMiddleware       gin.HandlerFunc
}

func NewAlbumRouter(
	albumController *AlbumController,
	albumCoverController *AlbumCoverController,
	albumTrackController *AlbumTrackController,
	authMiddleware gin.HandlerFunc,
) *AlbumRouter {
	return &AlbumRouter{
		albumController:      albumController,
		albumCoverController: albumCoverController,
		albumTrackController: albumTrackController,
		authMiddleware:       authMiddleware,
	}
}

func (r *AlbumRouter) RegisterRoutes(router *gin.RouterGroup) {
	albumGroup := router.Group("/albums")

	albumGroup.GET("/:id", r.albumController.GetAlbum)
	albumGroup.GET("/:id/tracks", r.albumTrackController.GetAlbumTracks)

	albumProtected := albumGroup.Group("")
	albumProtected.Use(r.authMiddleware)
	{
		albumProtected.POST("", r.albumController.CreateAlbum)
		albumProtected.PUT("/:id", r.albumController.UpdateAlbum)
		albumProtected.DELETE("/:id", r.albumController.DeleteAlbum)
	}

	coverGroup := albumGroup.Group("/:id/cover")
	coverGroup.GET("", r.albumCoverController.GetCover)

	coverProtected := coverGroup.Group("")
	coverProtected.Use(r.authMiddleware)
	{
		coverProtected.POST("", r.albumCoverController.UploadCover)
		coverProtected.DELETE("", r.albumCoverController.DeleteCover)
	}
}
