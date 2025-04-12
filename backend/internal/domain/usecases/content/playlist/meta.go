package playlist

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

type PlaylistMetaService interface {
	CreatePlaylist(ctx context.Context, claims *entity.Claims, playlist *entity.Playlist) error
	GetPlaylist(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (*entity.Playlist, error)
	GetUserPlaylists(ctx context.Context, claims *entity.Claims, userID uuid.UUID) ([]*entity.Playlist, error)
	UpdatePlaylist(ctx context.Context, claims *entity.Claims, playlist *entity.Playlist) error
	DeletePlaylist(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) error
}
