package track

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

type TracksegmentService interface {
	GetSegments(ctx context.Context, trackID uuid.UUID) ([]*entity.Segment, error)
	CreateSegments(ctx context.Context, tracksID uuid.UUID) error
	DeleteSegments(ctx context.Context) error
}
