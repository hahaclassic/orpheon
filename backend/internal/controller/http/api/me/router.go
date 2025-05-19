package me_ctrl

import (
	"github.com/gin-gonic/gin"
	playlist_ctrl "github.com/hahaclassic/orpheon/backend/internal/controller/http/api/content/playlist"
	user_ctrl "github.com/hahaclassic/orpheon/backend/internal/controller/http/api/user"
)

type MeRouter struct {
	playlistMetaController      *playlist_ctrl.PlaylistMetaController
	userController              *user_ctrl.UserController
	playlistFavoritesController *playlist_ctrl.PlaylistFavoritesController
	authMiddleware              gin.HandlerFunc
}

func NewMeRouter(playlistMetaController *playlist_ctrl.PlaylistMetaController,
	userController *user_ctrl.UserController,
	playlistFavoritesController *playlist_ctrl.PlaylistFavoritesController,
	authMiddleware gin.HandlerFunc) *MeRouter {

	return &MeRouter{
		playlistMetaController:      playlistMetaController,
		userController:              userController,
		playlistFavoritesController: playlistFavoritesController,
		authMiddleware:              authMiddleware,
	}
}

func (r *MeRouter) RegisterRoutes(router *gin.RouterGroup) {
	me := router.Group("/me")
	me.Use(r.authMiddleware)
	{
		me.GET("/playlists", r.playlistMetaController.GetMyPlaylists)
		me.GET("/favorites", r.playlistFavoritesController.GetFavoritePlaylists)
		me.POST("/favorites/:playlist_id", r.playlistFavoritesController.AddToFavorites)
		me.DELETE("/favorites/:playlist_id", r.playlistFavoritesController.RemoveFromFavorites)
		me.GET("", r.userController.GetMe)
		me.PUT("", r.userController.UpdateMe)
	}

	user := router.Group("/users")
	{
		user.GET("/:id", r.userController.GetUser)
		user.GET("/:id/playlists", r.playlistMetaController.GetUserPlaylists)
	}
}
