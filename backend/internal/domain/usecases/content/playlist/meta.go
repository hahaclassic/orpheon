package playlist

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

var (
	ErrCreatePlaylist   = errors.New("failed to create playlist")
	ErrGetPlaylist      = errors.New("failed to get playlist")
	ErrGetUserPlaylists = errors.New("failed to get user playlists")
	ErrUpdatePlaylist   = errors.New("failed to update playlist")
	ErrDeletePlaylist   = errors.New("failed to delete playlist")
)

type PlaylistMetaService interface {
	CreatePlaylist(ctx context.Context, claims *entity.Claims, playlist *entity.Playlist) error
	GetPlaylist(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (*entity.Playlist, error)
	GetUserPlaylists(ctx context.Context, claims *entity.Claims, userID uuid.UUID) ([]*entity.Playlist, error)
	UpdatePlaylist(ctx context.Context, claims *entity.Claims, playlist *entity.Playlist) error
	DeletePlaylist(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) error
}
