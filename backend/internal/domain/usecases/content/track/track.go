package track

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

type TrackMetaService interface {
	GetTrackMeta(ctx context.Context, trackID uuid.UUID) (*entity.TrackMeta, error)
	CreateTrackMeta(ctx context.Context, claims *entity.Claims, track *entity.TrackMeta) (uuid.UUID, error)
	UpdateTrackMeta(ctx context.Context, claims *entity.Claims, track *entity.TrackMeta) error
	DeleteTrackMeta(ctx context.Context, claims *entity.Claims, trackID uuid.UUID) error
}
