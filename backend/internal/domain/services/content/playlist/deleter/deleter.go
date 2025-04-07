package facade

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entities"
)

type PlaylistDeletionService interface {
	DeletePlaylist(ctx context.Context, claims *entities.Claims, playlistID uuid.UUID) error
}
