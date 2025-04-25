package favorites

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

type PlaylistFavoriteRepository interface {
	AddToFavorites(ctx context.Context, userID uuid.UUID, playlistID uuid.UUID) error
	GetUserFavorites(ctx context.Context, userID uuid.UUID) ([]*entity.PlaylistMeta, error)
	DeleteFromUserFavorites(ctx context.Context, userID uuid.UUID, trackID uuid.UUID) error

	GetUsersWithFavoritePlaylist(ctx context.Context, playlistID uuid.UUID) ([]uuid.UUID, error)
	DeleteFromAllFavorites(ctx context.Context, playlistID uuid.UUID) error
	RestoreAllFavorites(ctx context.Context, userIDs []uuid.UUID, playlistID uuid.UUID) error
}

type PlaylistFavoriteService struct {
	favoriteRepo  PlaylistFavoriteRepository
	policyService usecase.PlaylistPolicyService
}

func NewPlaylistFavoriteService(favoriteRepo PlaylistFavoriteRepository,
	policyService usecase.PlaylistPolicyService) *PlaylistFavoriteService {
	return &PlaylistFavoriteService{
		favoriteRepo:  favoriteRepo,
		policyService: policyService,
	}
}

func (s *PlaylistFavoriteService) AddToUserFavorites(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrAddToUserFavorites, err)
	}()

	err = s.policyService.CanView(ctx, claims, playlistID)
	if err != nil {
		return err
	}

	return s.favoriteRepo.AddToFavorites(ctx, claims.UserID, playlistID)
}

// Only user can view his favorite playlists
func (s *PlaylistFavoriteService) GetUserFavorites(ctx context.Context, claims *entity.Claims) (_ []*entity.PlaylistMeta, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetUserFavorites, err)
	}()

	return s.favoriteRepo.GetUserFavorites(ctx, claims.UserID)
}

func (s *PlaylistFavoriteService) DeleteFromUserFavorites(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrDeleteFromAllFavorites, err)
	}()

	if err = s.policyService.CanDelete(ctx, claims, playlistID); err != nil {
		return err
	}

	return s.favoriteRepo.DeleteFromUserFavorites(ctx, claims.UserID, playlistID)
}

func (s *PlaylistFavoriteService) GetUsersWithFavoritePlaylist(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (_ []uuid.UUID, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetUsersWithFavoritePlaylist, err)
	}()

	if err = s.policyService.CanView(ctx, claims, playlistID); err != nil {
		return nil, err
	}

	return s.favoriteRepo.GetUsersWithFavoritePlaylist(ctx, playlistID)
}

func (s *PlaylistFavoriteService) DeleteFromAllFavorites(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrDeleteFromAllFavorites, err)
	}()

	if err = s.policyService.CanDelete(ctx, claims, playlistID); err != nil {
		return err
	}

	return s.favoriteRepo.DeleteFromAllFavorites(ctx, playlistID)
}

func (s *PlaylistFavoriteService) AddPlaylistToAllFavorites(ctx context.Context, claims *entity.Claims, userIDs []uuid.UUID, playlistID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrAddPlaylistToAllFavorites, err)
	}()

	if err = s.policyService.CanDelete(ctx, claims, playlistID); err != nil {
		return err
	}

	return s.favoriteRepo.RestoreAllFavorites(ctx, userIDs, playlistID)
}
