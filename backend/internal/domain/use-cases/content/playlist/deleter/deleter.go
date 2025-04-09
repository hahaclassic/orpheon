package facade

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

type PlaylistDeletionService interface {
	DeletePlaylist(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) error
}
