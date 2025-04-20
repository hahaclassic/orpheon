package track

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

var (
	ErrDeleteTrack = errors.New("failed to delete track")
)

type TrackDeleter interface {
	DeleteTrack(ctx context.Context, claims *entity.Claims, trackID uuid.UUID) error
}
