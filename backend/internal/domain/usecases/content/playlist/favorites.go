package playlist

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

var (
	ErrAddToFavorites                  = errors.New("failed to add playlist to favorites")
	ErrGetUserFavorites                = errors.New("failed to get user favorites")
	ErrGetUsersWithFavoritePlaylist    = errors.New("failed to get users with favorite playlist")
	ErrDeleteFromFavorites             = errors.New("failed to delete playlist from favorites")
	ErrsDeletePlaylistFromAllFavorites = errors.New("failed to delete favorite playlist for all users")
)

type PlaylistFavoriteService interface {
	AddToFavorites(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) error
	GetUserFavorites(ctx context.Context, claims *entity.Claims) ([]uuid.UUID, error)
	GetUsersWithFavoritePlaylist(ctx context.Context, playlistID uuid.UUID) ([]uuid.UUID, error)
	DeleteFromFavorites(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) error
	DeletePlaylistFromAllFavorites(ctx context.Context, playlistID uuid.UUID) error
}
