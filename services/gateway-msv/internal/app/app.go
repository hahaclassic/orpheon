package app

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hahaclassic/orpheon/backend/internal/config"
	auth_ctrl "github.com/hahaclassic/orpheon/backend/internal/controller/http/api/auth"
	album_ctrl "github.com/hahaclassic/orpheon/backend/internal/controller/http/api/content/album"
	artist_ctrl "github.com/hahaclassic/orpheon/backend/internal/controller/http/api/content/artist"
	genre_ctrl "github.com/hahaclassic/orpheon/backend/internal/controller/http/api/content/genre"
	license_ctrl "github.com/hahaclassic/orpheon/backend/internal/controller/http/api/content/license"
	playlist_ctrl "github.com/hahaclassic/orpheon/backend/internal/controller/http/api/content/playlist"
	search_ctrl "github.com/hahaclassic/orpheon/backend/internal/controller/http/api/content/search"
	track_ctrl "github.com/hahaclassic/orpheon/backend/internal/controller/http/api/content/track"
	stats_ctrl "github.com/hahaclassic/orpheon/backend/internal/controller/http/api/stat"
	user_ctrl "github.com/hahaclassic/orpheon/backend/internal/controller/http/api/user"
	"github.com/hahaclassic/orpheon/backend/internal/controller/http/middleware"
	"github.com/hahaclassic/orpheon/backend/internal/controller/http/router"
	album_router "github.com/hahaclassic/orpheon/backend/internal/controller/http/router/router-registrators/album"
	artist_router "github.com/hahaclassic/orpheon/backend/internal/controller/http/router/router-registrators/artist"
	playlist_router "github.com/hahaclassic/orpheon/backend/internal/controller/http/router/router-registrators/playlist"
	track_router "github.com/hahaclassic/orpheon/backend/internal/controller/http/router/router-registrators/track"
	user_me_router "github.com/hahaclassic/orpheon/backend/internal/controller/http/router/router-registrators/user-me"
	"github.com/hahaclassic/orpheon/backend/internal/controller/http/utils/cookie"
)

func Run(conf *config.Config) {
	ctx := context.Background()

	cookieTokensSetter := cookie.NewCookieTokensSetter(&conf.Cookie)

	authMiddleware := middleware.NewAuthMiddleware(authService, cookieTokensSetter)
	authMiddlewareRequired := authMiddleware.Optional() //authMiddleware.Required()
	authMiddlewareOptional := authMiddleware.Optional()

	authController := auth_ctrl.NewAuthController(authService, cookieTokensSetter, authMiddlewareRequired)

	genreController := genre_ctrl.NewGenreController(genreService, authMiddlewareRequired)
	genreAssignController := genre_ctrl.NewGenreAssignController(genreAssignService)
	licenseController := license_ctrl.NewLicenseController(licenseService, authMiddlewareRequired)
	artistMetaController := artist_ctrl.NewArtistMetaController(artistMetaService)
	artistAvatarController := artist_ctrl.NewArtistAvatarController(artistAvatarService)
	artistAssignController := artist_ctrl.NewArtistAssignController(artistAssignService, contentAggregator)
	albumMetaController := album_ctrl.NewAlbumMetaController(albumMetaService, contentAggregator)
	albumTrackController := album_ctrl.NewAlbumTrackController(albumTrackService, contentAggregator)
	albumCoverController := album_ctrl.NewAlbumCoverController(albumCoverService)
	trackMetaController := track_ctrl.NewTrackMetaController(trackService, contentAggregator)
	trackAudioController := track_ctrl.NewTrackAudioController(trackAudioService)
	searchController := search_ctrl.NewSearchController(searchService, contentAggregator, playlistAggregator, authMiddlewareOptional)
	userController := user_ctrl.NewUserController(userService)
	playlistMetaController := playlist_ctrl.NewPlaylistMetaController(playlistMetaService,
		playlistDeletionService,
		playlistPrivacyService,
		playlistAggregator,
	)
	playlistCoverController := playlist_ctrl.NewPlaylistCoverController(playlistCoverService)
	playlistTrackController := playlist_ctrl.NewPlaylistTrackController(playlistTrackService, contentAggregator)
	playlistFavoriteController := playlist_ctrl.NewPlaylistFavoritesController(playlistFavoriteService, playlistAggregator)
	trackSegmentController := track_ctrl.NewTrackSegmentController(segmentService)
	statController := stats_ctrl.NewStatController(listeningStatService)

	albumRouter := album_router.NewAlbumRouter(
		albumMetaController, albumCoverController,
		albumTrackController, genreAssignController, authMiddlewareRequired)

	artistRouter := artist_router.NewArtistRouter(
		artistMetaController, artistAvatarController,
		artistAssignController, authMiddlewareRequired)

	playlistRouter := playlist_router.NewPlaylistRouter(
		playlistMetaController, playlistTrackController, playlistCoverController, authMiddlewareRequired)

	trackRouter := track_router.NewTrackRouter(trackMetaController,
		trackSegmentController, trackAudioController, statController, artistAssignController, authMiddlewareRequired)

	meRouter := user_me_router.NewMeRouter(playlistMetaController, userController,
		playlistFavoriteController, authMiddlewareRequired)

	loggerMiddleware, err := middleware.SetupLoggerMiddleware(conf.Logger.Path, conf.Logger.Level)
	if err != nil {
		slog.Error("failed to create logger middleware", "err", err)
		return
	}

	// Initialize router
	router := router.SetupRouter(
		[]router.RoutersRegistrator{
			authController,
			genreController,
			licenseController,
			searchController,

			albumRouter,
			artistRouter,
			playlistRouter,
			trackRouter,
			meRouter,
		},
		[]gin.HandlerFunc{
			loggerMiddleware,
			middleware.CORSMiddleware(),
		})

	// addr := net.JoinHostPort(conf.HTTP.Host, conf.HTTP.Port)
	// if err := router.Run(addr); err != nil {
	// 	slog.Error("failed to start HTTP server", "err", err)
	// }
	addr := net.JoinHostPort(conf.HTTP.Host, conf.HTTP.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("starting server", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "err", err)
		}
	}()

	<-ctx.Done()
	slog.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server forced to shutdown", "err", err)
	} else {
		slog.Info("server exited properly")
	}
}
