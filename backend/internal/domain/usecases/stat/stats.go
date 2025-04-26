package stats

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

var (
	ErrUpdateStat       = errors.New("failed to update listening stat")
	ErrGetTrackSegments = errors.New("failed to get track segments")
)

type ListeningStatService interface {
	UpdateStat(ctx context.Context, event *entity.ListeningEvent) error
	GetTrackSegments(ctx context.Context, trackID uuid.UUID) ([]*entity.Segment, error)
}
