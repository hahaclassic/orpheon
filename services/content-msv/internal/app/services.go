package app

import (
	"net"

	"github.com/hahaclassic/orpheon/services/content-msv/config"
	content_aggregator "github.com/hahaclassic/orpheon/services/content-msv/internal/domain/service/aggregator"
	"github.com/hahaclassic/orpheon/services/content-msv/internal/domain/service/album"
	"github.com/hahaclassic/orpheon/services/content-msv/internal/domain/service/artist"
	"github.com/hahaclassic/orpheon/services/content-msv/internal/domain/service/genre"
	"github.com/hahaclassic/orpheon/services/content-msv/internal/domain/service/license"
	playlist_aggregator "github.com/hahaclassic/orpheon/services/content-msv/internal/domain/service/playlist/aggregator"
	playlist_cover "github.com/hahaclassic/orpheon/services/content-msv/internal/domain/service/playlist/cover"
	playlist_deleter "github.com/hahaclassic/orpheon/services/content-msv/internal/domain/service/playlist/deleter"
	playlist_favorites "github.com/hahaclassic/orpheon/services/content-msv/internal/domain/service/playlist/favorites"
	playlist_meta "github.com/hahaclassic/orpheon/services/content-msv/internal/domain/service/playlist/meta"
	playlist_policy "github.com/hahaclassic/orpheon/services/content-msv/internal/domain/service/playlist/policy"
	playlist_privacy "github.com/hahaclassic/orpheon/services/content-msv/internal/domain/service/playlist/privacy"
	playlist_tracks "github.com/hahaclassic/orpheon/services/content-msv/internal/domain/service/playlist/tracks"
	"github.com/hahaclassic/orpheon/services/content-msv/internal/domain/service/search"
	"github.com/hahaclassic/orpheon/services/content-msv/internal/domain/service/track"
	grpc_user_info "github.com/hahaclassic/orpheon/services/content-msv/internal/providers/user-info/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	userproto "github.com/hahaclassic/orpheon/api/user-msv/v1/proto"
)

// TODO: change implementations to interfaces from /usecase
type Services struct {
	// Aggregator
	aggregator *content_aggregator.ContentAggregator

	// Album
	albumMeta  *album.AlbumMetaService
	albumCover *album.AlbumCoverService
	albumTrack *album.AlbumTrackService

	// Artist
	artistMeta   *artist.ArtistMetaService
	artistAvatar *artist.ArtistCoverService
	artistAssign *artist.ArtistAssignService

	// Track
	trackMeta    *track.TrackMetaService
	trackAudio   *track.AudioFileService
	trackSegment *track.TrackSegmentService
	trackStat    *track.ListeningStatService

	// Playlist
	playlistAggregator *playlist_aggregator.PlaylistAggregator
	playlistMeta       *playlist_meta.PlaylistMetaService
	playlistCover      *playlist_cover.PlaylistCoverService
	playlistPolicy     *playlist_policy.PlaylistPolicyService
	playlistPrivacy    *playlist_privacy.PlaylistPrivacyChanger
	playlistTracks     *playlist_tracks.PlaylistTrackService
	playlistFavorites  *playlist_favorites.PlaylistFavoriteService
	playlistDeleter    *playlist_deleter.PlaylistDeleter

	// Genres
	genreMeta   *genre.GenreService
	genreAssign *genre.GenreAssignService

	// Lisences
	licenseMeta *license.LicenseService

	// Search
	search *search.SearchService

	// UserInfo
	userInfo *grpc_user_info.UserInfoService
}

func InitServices(cfg *config.Config, repos *Repositories) (*Services, error) {
	s := &Services{
		// Album
		albumMeta: album.NewAlbumMetaService(repos.albumMeta),
		//albumCover: album.NewAlbumCoverService(repos.albumCover),
		albumTrack: album.NewAlbumTrackService(repos.albumTrack),

		// Artist
		artistMeta: artist.NewArtistMetaService(repos.artistMeta),
		//artistAvatar: artist.NewArtistCoverService(repos.artistAvatar),
		artistAssign: artist.NewArtistAssignService(repos.artistAssign),

		// Track
		//trackAudio:   track.NewAudioFileService(repos.trackAudio, audioconverter.New()),
		trackSegment: track.NewTrackSegmentService(repos.trackSegment),
		trackStat:    track.NewListeningStatService(repos.trackMeta, repos.trackSegment),

		// Genres
		genreMeta:   genre.NewGenreService(repos.genreMeta),
		genreAssign: genre.NewGenreAssignService(repos.genreAssign),

		// Lisences
		licenseMeta: license.NewLicenseService(repos.licenseMeta),

		// Search
		search: search.NewSearchService(repos.search),
	}

	err := s.initUserInfo(cfg)
	if err != nil {
		return nil, err
	}

	s.initPlaylist(cfg, repos)

	s.aggregator = content_aggregator.NewContentAggregator(s.trackMeta,
		s.artistAssign, s.albumMeta, s.licenseMeta, s.genreMeta)
	s.trackMeta = track.NewTrackMetaService(repos.trackMeta, s.trackSegment)

	return s, nil
}

func (s *Services) initUserInfo(cfg *config.Config) error {
	address := net.JoinHostPort(cfg.UserInfo.Host, cfg.UserInfo.Port)
	userCreatorGRPCClient, err := grpc.NewClient(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	userproto.NewUserServiceClient(userCreatorGRPCClient)

	s.userInfo = grpc_user_info.NewUserInfoService(userCreatorGRPCClient)

	return nil
}

func (s *Services) initPlaylist(cfg *config.Config, repos *Repositories) {
	s.playlistPolicy = playlist_policy.New(repos.playlistAccessWithCache)
	s.playlistPrivacy = playlist_privacy.NewPlaylistPrivacyChanger(
		s.playlistPolicy, s.playlistFavorites, repos.playlistAccessWithCache)
	s.playlistMeta = playlist_meta.NewPlaylistMetaService(repos.playlistMeta,
		s.playlistPolicy, repos.playlistAccessWithCache)
	//s.playlistCover = playlist_cover.New(repos.playlistCover, s.playlistPolicy)

	s.playlistTracks = playlist_tracks.NewPlaylistTrackService(repos.playlistTracks,
		s.playlistPolicy)
	s.playlistFavorites = playlist_favorites.NewPlaylistFavoriteService(
		repos.playlistFavorite, s.playlistPolicy)

	s.playlistDeleter = playlist_deleter.New(
		playlist_deleter.WithMetaDeletion(s.playlistMeta),
		playlist_deleter.WithCoverDeletion(s.playlistCover),
		playlist_deleter.WithTracksDeletion(s.playlistTracks),
		playlist_deleter.WithFavoritesDeletion(s.playlistFavorites),
	)
	s.playlistAggregator = playlist_aggregator.NewPlaylistAggregator(
		s.playlistMeta, s.playlistFavorites, s.playlistTracks, s.userInfo,
	)
}
