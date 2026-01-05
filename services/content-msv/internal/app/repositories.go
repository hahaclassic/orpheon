package app

import (
	"context"

	"github.com/hahaclassic/orpheon/services/content-msv/config"
	album_cover_minio "github.com/hahaclassic/orpheon/services/content-msv/internal/repository/album/cover/minio"
	album_meta_postgres "github.com/hahaclassic/orpheon/services/content-msv/internal/repository/album/meta/postgres"
	album_tracks_postgres "github.com/hahaclassic/orpheon/services/content-msv/internal/repository/album/tracks/postgres"
	assign_postgres "github.com/hahaclassic/orpheon/services/content-msv/internal/repository/artist/assign/postgres"
	avatar_minio "github.com/hahaclassic/orpheon/services/content-msv/internal/repository/artist/avatar/minio"
	artist_meta_postgres "github.com/hahaclassic/orpheon/services/content-msv/internal/repository/artist/meta/postgres"
	genre_assign_postgres "github.com/hahaclassic/orpheon/services/content-msv/internal/repository/genre/assign/postgres"
	genre_postgres "github.com/hahaclassic/orpheon/services/content-msv/internal/repository/genre/meta/postgres"
	license_postgres "github.com/hahaclassic/orpheon/services/content-msv/internal/repository/license/postgres"
	access_cache_local "github.com/hahaclassic/orpheon/services/content-msv/internal/repository/playlist/access-cache/local"
	access_cache_redis "github.com/hahaclassic/orpheon/services/content-msv/internal/repository/playlist/access-cache/redis"
	access_meta_postgres "github.com/hahaclassic/orpheon/services/content-msv/internal/repository/playlist/access-meta/default/postgres"
	access_meta "github.com/hahaclassic/orpheon/services/content-msv/internal/repository/playlist/access-meta/with-cache"
	playlist_cover_minio "github.com/hahaclassic/orpheon/services/content-msv/internal/repository/playlist/cover/minio"
	favorites_postgres "github.com/hahaclassic/orpheon/services/content-msv/internal/repository/playlist/favorites/postgres"
	playlist_meta_postgres "github.com/hahaclassic/orpheon/services/content-msv/internal/repository/playlist/meta/postgres"
	playlist_tracks_postgres "github.com/hahaclassic/orpheon/services/content-msv/internal/repository/playlist/tracks/postgres"
	search_postgres "github.com/hahaclassic/orpheon/services/content-msv/internal/repository/search/postgres"
	audio_minio "github.com/hahaclassic/orpheon/services/content-msv/internal/repository/track/audio/minio"
	track_meta_postgres "github.com/hahaclassic/orpheon/services/content-msv/internal/repository/track/meta/postgres"
	segment_postgres "github.com/hahaclassic/orpheon/services/content-msv/internal/repository/track/segment/postgres"
)

// TODO: change implementations to interfaces
type Repositories struct {
	// Album
	albumMeta  *album_meta_postgres.AlbumRepository
	albumTrack *album_tracks_postgres.AlbumTrackRepository
	albumCover *album_cover_minio.AlbumCoverRepository

	// Artist
	artistMeta   *artist_meta_postgres.ArtistMetaRepository
	artistAvatar *avatar_minio.ArtistAvatarRepository
	artistAssign *assign_postgres.ArtistAssignRepository

	// Track
	trackMeta    *track_meta_postgres.TrackMetaRepository
	trackSegment *segment_postgres.TrackSegmentRepository
	trackAudio   *audio_minio.AudioFileRepository

	// Playlist
	playlistMeta     *playlist_meta_postgres.PlaylistMetaRepository
	playlistTracks   *playlist_tracks_postgres.PlaylistTracksRepository
	playlistFavorite *favorites_postgres.PlaylistFavoriteRepository
	playlistCover    *playlist_cover_minio.PlaylistCoverRepository
	//playlistAccessLocal     *access_cache_local.AccessCache
	//playlistAccessRedis     *access_cache_redis.AccessCache
	//playlistAccessPostgres  *access_meta_postgres.PlaylistAccessRepository
	playlistAccessWithCache *access_meta.PlaylistAccessRepoWithCache

	// Genre
	genreMeta   *genre_postgres.GenreRepository
	genreAssign *genre_assign_postgres.GenreAssignRepository

	// License
	licenseMeta *license_postgres.LicenseRepository

	// Search
	search *search_postgres.SearchRepository
}

func InitRepositories(_ context.Context, cfg *config.Config, infra *Infra) (*Repositories, error) {
	r := &Repositories{}

	// Album
	r.albumMeta = album_meta_postgres.NewAlbumRepository(infra.postgres)
	r.albumTrack = album_tracks_postgres.NewAlbumTrackRepository(infra.postgres)

	// Artist
	r.artistMeta = artist_meta_postgres.NewArtistMetaRepository(infra.postgres)
	r.artistAssign = assign_postgres.NewArtistAssignRepository(infra.postgres)

	// Track
	r.trackMeta = track_meta_postgres.NewTrackMetaRepository(infra.postgres)
	r.trackSegment = segment_postgres.NewTrackSegmentRepository(infra.postgres)

	// Playlist
	r.playlistMeta = playlist_meta_postgres.NewPlaylistMetaRepository(infra.postgres)
	r.playlistTracks = playlist_tracks_postgres.NewPlaylistTracksRepository(infra.postgres)
	r.playlistFavorite = favorites_postgres.NewPlaylistFavoriteRepository(infra.postgres)

	// Genre
	r.genreMeta = genre_postgres.NewGenreRepository(infra.postgres)
	r.genreAssign = genre_assign_postgres.NewGenreAssignRepository(infra.postgres)

	// License
	r.licenseMeta = license_postgres.NewLicenseRepository(infra.postgres)

	// Search
	r.search = search_postgres.NewSearchRepository(infra.postgres)

	r.playlistMeta = playlist_meta_postgres.NewPlaylistMetaRepository(infra.postgres)
	r.playlistTracks = playlist_tracks_postgres.NewPlaylistTracksRepository(infra.postgres)
	r.playlistFavorite = favorites_postgres.NewPlaylistFavoriteRepository(infra.postgres)

	// mocks
	r.albumCover = &album_cover_minio.AlbumCoverRepository{}
	r.artistAvatar = &avatar_minio.ArtistAvatarRepository{}
	r.trackAudio = &audio_minio.AudioFileRepository{}
	r.playlistCover = &playlist_cover_minio.PlaylistCoverRepository{}
	// albumCoverRepo, err := album_cover_minio.NewAlbumCoverRepository(ctx, minioClient, conf.MinIOBuckets.BucketAlbum)
	// if err != nil {
	// 	slog.Error("failed to create album cover repository", "err", err)
	// 	return
	// }
	// r.albumCover = albumCover

	// artistAvatarRepo, err := avatar_minio.NewArtistAvatarRepository(ctx, minioClient, conf.MinIOBuckets.BucketArtistAvatar)
	// if err != nil {
	// 	slog.Error("failed to create artist avatar repository", "err", err)
	// 	return
	// }

	// audioRepo, err := audio_minio.NewAudioFileRepository(ctx, minioClient, conf.MinIOBuckets.BucketAudio)
	// if err != nil {
	// 	slog.Error("failed to create audio file repository", "err", err)
	// 	return
	// }

	// playlistCoverRepo, err := playlist_cover_minio.NewPlaylistCoverRepository(ctx, minioClient, conf.MinIOBuckets.BucketPlaylist)
	// if err != nil {
	// 	slog.Error("failed to create playlist cover repository", "err", err)
	// 	return
	// }

	playlistAccessPostgres := access_meta_postgres.NewPlaylistAccessRepository(infra.postgres)
	accessCacheLocal, err := access_cache_local.NewAccessCache(cfg.LocalAccessCache.Size)
	if err != nil {
		return nil, err
	}

	accessMetaCacheRedis := access_cache_redis.NewAccessCache(infra.redis, &cfg.RedisAccessCache)
	r.playlistAccessWithCache = access_meta.New(playlistAccessPostgres,
		access_meta.WithL1Cache(accessCacheLocal),
		access_meta.WithL2Cache(accessMetaCacheRedis))

	return r, nil
}
