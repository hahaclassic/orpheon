package favorites

import (
	"context"

	"github.com/google/uuid"
)

type PlaylistFavoriteService interface {
	AddToFavorites(ctx context.Context, userID uuid.UUID, playlistID uuid.UUID) error
	GetUserFavorites(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	GetUsersWithFavoritePlaylist(ctx context.Context, playlistID uuid.UUID) ([]uuid.UUID, error)
	DeleteFromFavorites(ctx context.Context, userID uuid.UUID, playlistID uuid.UUID) error
	DeletePlaylistFromAllFavorites(ctx context.Context, playlistID uuid.UUID) error
}
