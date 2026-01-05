package app

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hahaclassic/orpheon/pkg/http/router"
	"github.com/hahaclassic/orpheon/services/content-msv/config"
	album_ctrl "github.com/hahaclassic/orpheon/services/content-msv/internal/controller/http/api/album"
	artist_ctrl "github.com/hahaclassic/orpheon/services/content-msv/internal/controller/http/api/artist"
	genre_ctrl "github.com/hahaclassic/orpheon/services/content-msv/internal/controller/http/api/genre"
	license_ctrl "github.com/hahaclassic/orpheon/services/content-msv/internal/controller/http/api/license"
	playlist_ctrl "github.com/hahaclassic/orpheon/services/content-msv/internal/controller/http/api/playlist"
	search_ctrl "github.com/hahaclassic/orpheon/services/content-msv/internal/controller/http/api/search"
	track_ctrl "github.com/hahaclassic/orpheon/services/content-msv/internal/controller/http/api/track"
	album_router "github.com/hahaclassic/orpheon/services/content-msv/internal/controller/http/router-registrators/album"
	artist_router "github.com/hahaclassic/orpheon/services/content-msv/internal/controller/http/router-registrators/artist"
	playlist_router "github.com/hahaclassic/orpheon/services/content-msv/internal/controller/http/router-registrators/playlist"
	track_router "github.com/hahaclassic/orpheon/services/content-msv/internal/controller/http/router-registrators/track"
	"github.com/hahaclassic/orpheon/services/content-msv/internal/providers/http/middleware"
)

const apiBasePath = "api/v1"

func InitServerHTTP(cfg *config.Config, s *Services) *http.Server {
	claimsRequired, claimsOptional := middleware.ClaimsRequired(), middleware.ClaimsOptional()

	// Album
	albumMetaController := album_ctrl.NewAlbumMetaController(s.albumMeta, s.aggregator)
	albumTrackController := album_ctrl.NewAlbumTrackController(s.albumTrack, s.aggregator)
	albumCoverController := album_ctrl.NewAlbumCoverController(s.albumCover)

	// Artist
	artistMetaController := artist_ctrl.NewArtistMetaController(s.artistMeta)
	artistAvatarController := artist_ctrl.NewArtistAvatarController(s.artistAvatar)
	artistAssignController := artist_ctrl.NewArtistAssignController(s.artistAssign, s.aggregator)

	// Track
	trackMetaController := track_ctrl.NewTrackMetaController(s.trackMeta, s.aggregator)
	trackAudioController := track_ctrl.NewTrackAudioController(s.trackAudio)
	trackSegmentController := track_ctrl.NewTrackSegmentController(s.trackSegment)
	statController := track_ctrl.NewStatController(s.trackStat)

	// Playlist
	playlistMetaController := playlist_ctrl.NewPlaylistMetaController(
		s.playlistMeta,
		s.playlistDeleter,
		s.playlistPrivacy,
		s.playlistAggregator,
	)
	playlistCoverController := playlist_ctrl.NewPlaylistCoverController(s.playlistCover)
	playlistTrackController := playlist_ctrl.NewPlaylistTrackController(
		s.playlistTracks, s.aggregator)
	playlistFavoriteController := playlist_ctrl.NewPlaylistFavoritesController(
		s.playlistFavorites, s.playlistAggregator)

	// Genre
	genreController := genre_ctrl.NewGenreController(s.genreMeta, claimsRequired)
	genreAssignController := genre_ctrl.NewGenreAssignController(s.genreAssign)

	// License
	licenseController := license_ctrl.NewLicenseController(s.licenseMeta, claimsRequired)

	// Search
	searchController := search_ctrl.NewSearchController(s.search, s.aggregator,
		s.playlistAggregator, claimsOptional)

	// Routers
	albumRouter := album_router.NewAlbumRouter(
		albumMetaController, albumCoverController,
		albumTrackController, genreAssignController, claimsRequired)

	artistRouter := artist_router.NewArtistRouter(
		artistMetaController, artistAvatarController,
		artistAssignController, claimsRequired)

	trackRouter := track_router.NewTrackRouter(trackMetaController,
		trackSegmentController, trackAudioController,
		statController, artistAssignController, claimsRequired)

	playlistRouter := playlist_router.NewPlaylistRouter(
		playlistMetaController, playlistTrackController,
		playlistCoverController, playlistFavoriteController, claimsRequired)

	// Initialize router
	router := router.SetupRouter(apiBasePath,
		[]router.RoutersRegistrator{
			genreController,
			licenseController,
			searchController,

			albumRouter,
			artistRouter,
			playlistRouter,
			trackRouter,
		},
		[]gin.HandlerFunc{
			// loggerMiddleware,
			// middleware.CORSMiddleware(),
		})

	addr := net.JoinHostPort(cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	return server
}
