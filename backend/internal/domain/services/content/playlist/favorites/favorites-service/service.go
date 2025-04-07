package favoritesservice

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entities"
)

type playlistPolicyService interface {
	CanView(ctx context.Context, playlist uuid.UUID, userID uuid.UUID) (bool, error)
}

type playlistFavoriteRepository interface {
	AddToFavorites(ctx context.Context, userID uuid.UUID, playlistID uuid.UUID) error
	RemoveFromFavorites(ctx context.Context, userID uuid.UUID, playlistID uuid.UUID) error
	GetFavoritePlaylists(ctx context.Context, userID uuid.UUID) ([]*entities.Playlist, error)
	GetUsersWithFavoritePlaylist(ctx context.Context, playlistID uuid.UUID) ([]uuid.UUID, error)
	RemovePlaylistFromAllFavorites(ctx context.Context, playlistID uuid.UUID) error
}

type PlaylistFavoriteService struct {
	favoriteRepo  playlistFavoriteRepository
	policyService playlistPolicyService
}

func NewPlaylistFavoriteService(favoriteRepo playlistFavoriteRepository,
	policyService playlistPolicyService) *PlaylistFavoriteService {
	return &PlaylistFavoriteService{
		favoriteRepo:  favoriteRepo,
		policyService: policyService,
	}
}

func (s *PlaylistFavoriteService) AddToFavorites(ctx context.Context, userID uuid.UUID, playlistID uuid.UUID) error {
	canView, err := s.policyService.CanView(ctx, playlistID, userID)
	if err != nil {
		return err
	}
	if !canView {
		return errors.New("user does not have permission to view the playlist")
	}

	err = s.favoriteRepo.AddToFavorites(ctx, userID, playlistID)
	if err != nil {
		return err
	}
	return nil
}

func (s *PlaylistFavoriteService) GetUserFavorites(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	playlists, err := s.favoriteRepo.GetFavoritePlaylists(ctx, userID)
	if err != nil {
		return nil, err
	}

	var favoritePlaylistIDs []uuid.UUID
	for _, playlist := range playlists {
		favoritePlaylistIDs = append(favoritePlaylistIDs, playlist.ID)
	}
	return favoritePlaylistIDs, nil
}

func (s *PlaylistFavoriteService) GetUsersWithFavoritePlaylist(ctx context.Context, playlistID uuid.UUID) ([]uuid.UUID, error) {
	usersWithFavorite, err := s.favoriteRepo.GetUsersWithFavoritePlaylist(ctx, playlistID)
	if err != nil {
		return nil, err
	}
	return usersWithFavorite, nil
}

func (s *PlaylistFavoriteService) DeleteFromFavorites(ctx context.Context, userID uuid.UUID, playlistID uuid.UUID) error {
	err := s.favoriteRepo.RemoveFromFavorites(ctx, userID, playlistID)
	if err != nil {
		return err
	}
	return nil
}

func (s *PlaylistFavoriteService) DeletePlaylistFromAllFavorites(ctx context.Context, playlistID uuid.UUID) error {
	err := s.favoriteRepo.RemovePlaylistFromAllFavorites(ctx, playlistID)
	if err != nil {
		return err
	}
	return nil
}
