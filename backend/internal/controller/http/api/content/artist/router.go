package artist_ctrl

import "github.com/gin-gonic/gin"

type ArtistRouter struct {
	artistMetaController   *ArtistMetaController
	artistAvatarController *ArtistAvatarController
	artistAssignController *ArtistAssignController
	authMiddleware         gin.HandlerFunc
}

func NewArtistRouter(
	artistMetaController *ArtistMetaController,
	artistAvatarController *ArtistAvatarController,
	artistAssignController *ArtistAssignController,
	authMiddleware gin.HandlerFunc,
) *ArtistRouter {
	return &ArtistRouter{
		artistMetaController:   artistMetaController,
		artistAvatarController: artistAvatarController,
		artistAssignController: artistAssignController,
		authMiddleware:         authMiddleware,
	}
}

func (r *ArtistRouter) RegisterRoutes(router *gin.RouterGroup) {
	artistGroup := router.Group("/artists")

	artistGroup.GET("", r.artistMetaController.GetAllArtists)
	artistGroup.GET("/:id", r.artistMetaController.GetArtist)
	artistGroup.GET("/:id/albums", r.artistAssignController.GetAlbumsByArtist)
	artistGroup.GET("/:id/tracks", r.artistAssignController.GetTracksByArtist)

	artistProtected := artistGroup.Group("")
	artistProtected.Use(r.authMiddleware)
	{
		artistProtected.POST("", r.artistMetaController.CreateArtist)
		artistProtected.PUT("/:id", r.artistMetaController.UpdateArtist)
		artistProtected.DELETE("/:id", r.artistMetaController.DeleteArtist)

		artistProtected.POST("/:id/tracks/:track_id", r.artistAssignController.AssignArtistToTrack)
		artistProtected.DELETE("/:id/tracks/:track_id", r.artistAssignController.UnassignArtistFromTrack)
		artistProtected.POST("/:id/albums/:album_id", r.artistAssignController.AssignArtistToAlbum)
		artistProtected.DELETE("/:id/albums/:album_id", r.artistAssignController.UnassignArtistFromAlbum)
	}

	avatarGroup := artistGroup.Group("/:id/avatar")
	avatarGroup.GET("", r.artistAvatarController.GetAvatar)

	avatarProtected := avatarGroup.Group("")
	avatarProtected.Use(r.authMiddleware)
	{
		avatarProtected.POST("", r.artistAvatarController.UploadAvatar)
		avatarProtected.DELETE("", r.artistAvatarController.DeleteAvatar)
	}
}
