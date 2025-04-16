package favorites

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

type playlistFavoriteRepository interface {
	AddToFavorites(ctx context.Context, userID uuid.UUID, playlistID uuid.UUID) error
	RemoveFromFavorites(ctx context.Context, userID uuid.UUID, playlistID uuid.UUID) error
	GetFavoritePlaylists(ctx context.Context, userID uuid.UUID) ([]*entity.Playlist, error)
	GetUsersWithFavoritePlaylist(ctx context.Context, playlistID uuid.UUID) ([]uuid.UUID, error)
	RemovePlaylistFromAllFavorites(ctx context.Context, playlistID uuid.UUID) error
}

type PlaylistFavoriteService struct {
	favoriteRepo  playlistFavoriteRepository
	policyService usecase.PlaylistPolicyService
}

func NewPlaylistFavoriteService(favoriteRepo playlistFavoriteRepository,
	policyService usecase.PlaylistPolicyService) *PlaylistFavoriteService {
	return &PlaylistFavoriteService{
		favoriteRepo:  favoriteRepo,
		policyService: policyService,
	}
}

func (s *PlaylistFavoriteService) AddToFavorites(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrAddToFavorites, err)
	}()

	err = s.policyService.CanView(ctx, claims, playlistID)
	if err != nil {
		return err
	}

	return s.favoriteRepo.AddToFavorites(ctx, claims.UserID, playlistID)
}

// Only user can view his favorite playlists
func (s *PlaylistFavoriteService) GetUserFavorites(ctx context.Context, claims *entity.Claims) (_ []uuid.UUID, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetUserFavorites, err)
	}()

	playlists, err := s.favoriteRepo.GetFavoritePlaylists(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	var favoritePlaylistIDs []uuid.UUID
	for _, playlist := range playlists {
		favoritePlaylistIDs = append(favoritePlaylistIDs, playlist.ID)
	}

	return favoritePlaylistIDs, nil
}

// TODO
func (s *PlaylistFavoriteService) DeleteFromFavorites(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrDeleteFromFavorites, err)
	}()

	// err = s.policyService.CanView(ctx, claims, playlistID)
	// if err != nil {
	// 	return err
	// }

	return s.favoriteRepo.RemoveFromFavorites(ctx, claims.UserID, playlistID)
}

// IDK (mb for recoverable transaction)
func (s *PlaylistFavoriteService) GetUsersWithFavoritePlaylist(ctx context.Context, playlistID uuid.UUID) (_ []uuid.UUID, err error) {
	usersWithFavorite, err := s.favoriteRepo.GetUsersWithFavoritePlaylist(ctx, playlistID)
	if err != nil {
		return nil, err
	}
	return usersWithFavorite, nil
}

// for deleter
func (s *PlaylistFavoriteService) DeletePlaylistFromAllFavorites(ctx context.Context, playlistID uuid.UUID) error {
	err := s.favoriteRepo.RemovePlaylistFromAllFavorites(ctx, playlistID)
	if err != nil {
		return err
	}
	return nil
}
