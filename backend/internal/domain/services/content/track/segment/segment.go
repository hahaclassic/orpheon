package tracksegment

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

const (
	defaultSegmentCount = 10
)

var (
	ErrInvalidTrackID  = errors.New("invalid track ID")
	ErrInvalidDuration = errors.New("track duration must be positive")

	ErrCreateSegments   = errors.New("failed to create segments")
	ErrDeleteSegments   = errors.New("failed to delete segments")
	ErrGetSegments      = errors.New("failed to get segments")
	ErrIncrementStreams = errors.New("failed to increment total streams")
)

type TrackSegmentRepository interface {
	GetByTrackID(ctx context.Context, trackID uuid.UUID) ([]*entity.Segment, error)
	Create(ctx context.Context, segments []*entity.Segment) error
	DeleteByTrackID(ctx context.Context, trackID uuid.UUID) error
	IncrementTotalStreams(ctx context.Context, trackID uuid.UUID, segmentIdxs []int) error
}

type Service struct {
	repo TrackSegmentRepository
}

func NewTrackSegmentService(repo TrackSegmentRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetSegments(ctx context.Context, trackID uuid.UUID) ([]*entity.Segment, error) {
	if trackID == uuid.Nil {
		return nil, ErrInvalidTrackID
	}
	segments, err := s.repo.GetByTrackID(ctx, trackID)
	if err != nil {
		return nil, ErrGetSegments
	}
	return segments, nil
}

func (s *Service) CreateSegments(ctx context.Context, trackID uuid.UUID, trackDuration int64) error {
	if trackID == uuid.Nil {
		return ErrInvalidTrackID
	}
	if trackDuration <= 0 {
		return ErrInvalidDuration
	}

	segmentDuration := trackDuration / defaultSegmentCount
	segments := make([]*entity.Segment, 0, defaultSegmentCount)

	for i := 0; i < defaultSegmentCount; i++ {
		start := int64(i) * segmentDuration
		end := start + segmentDuration
		if i == defaultSegmentCount-1 {
			end = trackDuration // последний сегмент — до конца трека
		}
		segments = append(segments, &entity.Segment{
			TrackID:      trackID,
			Index:        i,
			StartMillis:  start,
			EndMillis:    end,
			TotalStreams: 0,
		})
	}

	if err := s.repo.Create(ctx, segments); err != nil {
		return ErrCreateSegments
	}
	return nil
}

func (s *Service) DeleteSegments(ctx context.Context, trackID uuid.UUID) error {
	if trackID == uuid.Nil {
		return ErrInvalidTrackID
	}
	if err := s.repo.DeleteByTrackID(ctx, trackID); err != nil {
		return ErrDeleteSegments
	}
	return nil
}

func (s *Service) IncrementSegmentsTotalStreams(ctx context.Context, trackID uuid.UUID, segmentIdxs []int) error {
	if trackID == uuid.Nil {
		return ErrInvalidTrackID
	}
	if err := s.repo.IncrementTotalStreams(ctx, trackID, segmentIdxs); err != nil {
		return ErrIncrementStreams
	}
	return nil
}
