package playlist

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

var (
	ErrAddToUserFavorites           = errors.New("failed to add playlist to user favorites")
	ErrGetUserFavorites             = errors.New("failed to get user favorites")
	ErrGetUsersWithFavoritePlaylist = errors.New("failed to get users with favorite playlist")
	ErrDeleteFromUserFavorites      = errors.New("failed to delete playlist from user favorites")
	ErrDeleteFromAllFavorites       = errors.New("failed to delete favorite playlist for all users")
	ErrAddPlaylistToAllFavorites    = errors.New("failed to add playlist to all favorites")
)

type PlaylistFavoriteService interface {
	AddToFavorites(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) error
	GetUserFavorites(ctx context.Context, claims *entity.Claims) ([]uuid.UUID, error)
	DeleteFromFavorites(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) error

	AddPlaylistToAllFavorites(ctx context.Context, claims *entity.Claims, userIDs []uuid.UUID, playlistID uuid.UUID) error
	GetUsersWithFavoritePlaylist(ctx context.Context, playlistID uuid.UUID) ([]uuid.UUID, error)
	DeletePlaylistFromAllFavorites(ctx context.Context, claims *entity.Claims, userIDs []uuid.UUID, playlistID uuid.UUID) error
}
