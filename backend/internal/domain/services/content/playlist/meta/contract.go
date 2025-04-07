package meta

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entities"
)

type PlaylistMetaService interface {
	CreatePlaylist(ctx context.Context, claims *entities.Claims, playlist *entities.Playlist) error
	GetPlaylist(ctx context.Context, claims *entities.Claims, playlistID uuid.UUID) (*entities.Playlist, error)
	GetUserPlaylists(ctx context.Context, claims *entities.Claims, userID uuid.UUID) ([]*entities.Playlist, error)
	UpdatePlaylist(ctx context.Context, claims *entities.Claims, playlist *entities.Playlist) error
	DeletePlaylist(ctx context.Context, claims *entities.Claims, playlistID uuid.UUID) error
}
