package playlist_aggregator

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/pkg/errwrap"
	"github.com/hahaclassic/orpheon/services/content-msv/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/services/content-msv/internal/domain/usecase/playlist"
)

type UserService interface {
	GetUser(ctx context.Context, userID uuid.UUID) (*entity.User, error)
}

type PlaylistAggregator struct {
	playlistMetaService     usecase.PlaylistMetaService
	playlistFavoriteService usecase.PlaylistFavoriteService
	playlistTrackService    usecase.PlaylistTrackService
	userService             UserService
}

func NewPlaylistAggregator(playlistMetaService usecase.PlaylistMetaService,
	playlistFavoriteService usecase.PlaylistFavoriteService,
	playlistTrackService usecase.PlaylistTrackService,
	userService UserService,
) *PlaylistAggregator {
	return &PlaylistAggregator{
		playlistMetaService,
		playlistFavoriteService,
		playlistTrackService,
		userService,
	}
}

func (a *PlaylistAggregator) GetPlaylistsByIDs(ctx context.Context, claims *entity.Claims, playlistIDs ...uuid.UUID) (_ []*entity.PlaylistMetaAggregated, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetPlaylistsAggregated, err)
	}()

	playlists := make([]*entity.PlaylistMetaAggregated, len(playlistIDs))

	for i, id := range playlistIDs {
		playlistMeta, err := a.playlistMetaService.GetMeta(ctx, claims, id)
		if err != nil {
			return nil, err
		}

		owner, err := a.userService.GetUser(ctx, playlistMeta.OwnerID)
		if err != nil {
			return nil, err
		}

		isFavorite, err := a.playlistFavoriteService.IsFavorite(ctx, claims, id)
		if err != nil {
			return nil, err
		}

		tracks, err := a.playlistTrackService.GetAllTracks(ctx, claims, id)
		if err != nil {
			return nil, err
		}

		playlists[i] = &entity.PlaylistMetaAggregated{
			ID:          playlistMeta.ID,
			Owner:       owner,
			Name:        playlistMeta.Name,
			Description: playlistMeta.Description,
			IsPrivate:   playlistMeta.IsPrivate,
			CreatedAt:   playlistMeta.CreatedAt,
			UpdatedAt:   playlistMeta.UpdatedAt,
			Rating:      playlistMeta.Rating,
			IsFavorite:  isFavorite,
			TracksCount: len(tracks),
			Tracks:      tracks,
		}
	}

	return playlists, nil
}

func (a *PlaylistAggregator) GetPlaylists(ctx context.Context, claims *entity.Claims, playlists ...*entity.PlaylistMeta) (_ []*entity.PlaylistMetaAggregated, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetPlaylistsAggregated, err)
	}()

	aggregated := make([]*entity.PlaylistMetaAggregated, len(playlists))

	for i, playlist := range playlists {
		owner, err := a.userService.GetUser(ctx, playlist.OwnerID)
		if err != nil {
			return nil, err
		}

		isFavorite, err := a.playlistFavoriteService.IsFavorite(ctx, claims, playlist.ID)
		if err != nil {
			return nil, err
		}

		tracks, err := a.playlistTrackService.GetAllTracks(ctx, claims, playlist.ID)
		if err != nil {
			return nil, err
		}

		aggregated[i] = &entity.PlaylistMetaAggregated{
			ID:          playlist.ID,
			Owner:       owner,
			Name:        playlist.Name,
			Description: playlist.Description,
			IsPrivate:   playlist.IsPrivate,
			CreatedAt:   playlist.CreatedAt,
			UpdatedAt:   playlist.UpdatedAt,
			IsFavorite:  isFavorite,
			Rating:      playlist.Rating,
			TracksCount: len(tracks),
			Tracks:      tracks,
		}
	}

	return aggregated, nil
}
