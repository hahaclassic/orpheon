package stats

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

type ListeningStatService interface {
	UpdateStat(ctx context.Context, event *entity.ListeningEvent) error
	GetTrackSegments(ctx context.Context, trackID uuid.UUID) ([]*entity.Segment, error)
}
