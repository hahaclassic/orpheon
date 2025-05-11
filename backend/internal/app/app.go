package app

import (
	"log/slog"
	"net"

	bcrypt_hasher "github.com/hahaclassic/orpheon/backend/internal/adapters/password-hasher/bcrypt-hasher"
	jwttokens "github.com/hahaclassic/orpheon/backend/internal/adapters/tokens/jwt"
	"github.com/hahaclassic/orpheon/backend/internal/config"
	auth_ctrl "github.com/hahaclassic/orpheon/backend/internal/controller/http/api/auth"
	albumcontroller "github.com/hahaclassic/orpheon/backend/internal/controller/http/api/content/album"
	artistcontroller "github.com/hahaclassic/orpheon/backend/internal/controller/http/api/content/artist"
	genrecontroller "github.com/hahaclassic/orpheon/backend/internal/controller/http/api/content/genre"
	licensecontroller "github.com/hahaclassic/orpheon/backend/internal/controller/http/api/content/license"
	searchcontroller "github.com/hahaclassic/orpheon/backend/internal/controller/http/api/content/search"
	trackcontroller "github.com/hahaclassic/orpheon/backend/internal/controller/http/api/content/track"
	"github.com/hahaclassic/orpheon/backend/internal/controller/http/middleware"
	"github.com/hahaclassic/orpheon/backend/internal/controller/http/router"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/auth"
	albummeta "github.com/hahaclassic/orpheon/backend/internal/domain/services/content/album/meta"
	artistmeta "github.com/hahaclassic/orpheon/backend/internal/domain/services/content/artist/meta"
	genreservice "github.com/hahaclassic/orpheon/backend/internal/domain/services/content/genre"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/content/license"
	searchservice "github.com/hahaclassic/orpheon/backend/internal/domain/services/content/search"
	trackmeta "github.com/hahaclassic/orpheon/backend/internal/domain/services/content/track/meta"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/user"
	"github.com/hahaclassic/orpheon/backend/internal/infrastructure/postgres"
	"github.com/hahaclassic/orpheon/backend/internal/infrastructure/redis"
	auth_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/auth/auth-repo/postgres"
	refresh_redis "github.com/hahaclassic/orpheon/backend/internal/repository/auth/refresh-token/redis"
	album_meta_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/content/album/meta/postgres"
	artist_meta_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/content/artist/meta/postgres"
	genre_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/content/genre/postgres"
	license_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/content/license/postgres"
	search_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/content/search/postgres"
	track_meta_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/content/track/meta/postgres"
	user_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/user/postgres"
)

func Run(conf *config.Config) {
	pgxpool := postgres.NewPostgresPool(conf.Postgres)
	defer pgxpool.Close()

	redisClient, err := redis.NewRedisClient(conf.Redis)
	if err != nil {
		slog.Error("redisClient", "err", err)
		return
	}
	defer redisClient.Close()

	// Initialize repositories
	authRepo := auth_postgres.NewAuthRepository(pgxpool)
	userRepo := user_postgres.NewUserRepository(pgxpool)
	refreshRepo := refresh_redis.NewRefreshTokenRepository(redisClient, &conf.RefreshToken)
	trackRepo := track_meta_postgres.NewTrackMetaRepository(pgxpool)
	albumRepo := album_meta_postgres.NewAlbumRepository(pgxpool)
	artistRepo := artist_meta_postgres.NewArtistMetaRepository(pgxpool)
	genreRepo := genre_postgres.NewGenreRepository(pgxpool)
	licenseRepo := license_postgres.NewLicenseRepository(pgxpool)
	searchRepo := search_postgres.NewSearchRepository(pgxpool)
	// playlistRepo := playlist_meta_postgres.NewPlaylistMetaRepository(pgxpool)

	// playlistAccessRepo := access_meta_postgres.NewPlaylistAccessRepository(pgxpool)
	// accessCacheLocal, err := access_cache_local.NewAccessCache(conf.LocalAccessMetaCache.Size)
	// if err != nil {
	// 	slog.Error("failed to create access cache local", "err", err)
	// 	return
	// }
	// accessMetaCacheRedis := access_cache_redis.NewAccessCache(redisClient, &conf.RedisAccessMetaCache)
	// playlistAccessRepoWithCache := access_meta.New(playlistAccessRepo,
	// 	access_meta.WithL1Cache(accessCacheLocal),
	// 	access_meta.WithL2Cache(accessMetaCacheRedis))

	// Initialize services
	hasher := bcrypt_hasher.New(conf.PasswordHasher.Cost)
	tokenService := jwttokens.New(conf.AccessToken)
	userService := user.New(userRepo)
	authService := auth.NewAuthService(authRepo, refreshRepo, userService, hasher, tokenService)
	//playlistPolicyService := policy.New(playlistAccessRepoWithCache)

	// Initialize content services
	trackService := trackmeta.NewTrackMetaService(trackRepo)
	albumService := albummeta.New(albumRepo)
	artistService := artistmeta.New(artistRepo)
	//playlistService := playlistmeta.NewPlaylistMetaService(playlistRepo, playlistPolicyService)
	genreService := genreservice.NewGenreService(genreRepo)
	searchService := searchservice.NewSearchService(searchRepo)
	licenseService := license.NewLicenseService(licenseRepo)

	// Initialize controllers
	authMiddleware := middleware.AuthMiddleware(authService)

	authController := auth_ctrl.NewAuthController(authService)
	trackController := trackcontroller.New(trackService, authMiddleware)
	albumController := albumcontroller.New(albumService, authMiddleware)
	artistController := artistcontroller.New(artistService, authMiddleware)
	//playlistController := playlistcontroller.New(playlistService)
	genreController := genrecontroller.New(genreService, authMiddleware)
	licenseController := licensecontroller.New(licenseService, authMiddleware)
	searchController := searchcontroller.New(searchService)

	// Initialize router
	router := router.SetupRouter(
		authController,
		trackController,
		albumController,
		artistController,
		genreController,
		// playlistController,q
		searchController,
		licenseController,
	)

	addr := net.JoinHostPort(conf.HTTP.Host, conf.HTTP.Port)
	if err := router.Run(addr); err != nil {
		slog.Error("failed to start HTTP server", "err", err)
	}
}
