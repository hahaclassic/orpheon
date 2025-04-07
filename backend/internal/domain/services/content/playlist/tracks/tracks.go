package tracks

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entities"
)

type PlaylistTrackService interface {
	AddTrack(ctx context.Context, claims *entities.Claims, playlistID uuid.UUID, trackID uuid.UUID) error
	GetAllTracks(ctx context.Context, claims *entities.Claims, playlistID uuid.UUID) error
	DeleteTrack(ctx context.Context, claims *entities.Claims, playlistID uuid.UUID, trackID uuid.UUID) error
	DeleteAllTracks(ctx context.Context, claims *entities.Claims, playlistID uuid.UUID) error
}
