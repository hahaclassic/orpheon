package processor

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/stat"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

const (
	MinSeconds = 30 // the minimum number of listening seconds to count
)

type ListeningStatService struct {
	trackRepo   TrackStatRepository
	segmentRepo SegmentStatRepository
}

type TrackStatRepository interface {
	IncrementTrackStreamCount(ctx context.Context, trackID uuid.UUID, userID uuid.UUID) error
}

type SegmentStatRepository interface {
	GetTrackSegments(ctx context.Context, trackID uuid.UUID) ([]*entity.Segment, error)
	IncrementSegmentStreamCount(ctx context.Context, trackID uuid.UUID, segmentsIdxs []int) error
}

func NewListeningStatService(trackRepo TrackStatRepository, segmentRepo SegmentStatRepository) *ListeningStatService {
	return &ListeningStatService{trackRepo: trackRepo, segmentRepo: segmentRepo}
}

func (s *ListeningStatService) UpdateStat(ctx context.Context, event *entity.ListeningEvent) (err error) {
	defer func() {
		if err != nil {
			err = errwrap.Wrap(usecase.ErrUpdateStat, err)
		}
	}()

	segments, err := s.segmentRepo.GetTrackSegments(ctx, event.TrackID)
	if err != nil {
		return err
	}

	affectedSegIdx, totalDuration := s.proccessListeningEvent(segments, event)

	if err = s.segmentRepo.IncrementSegmentStreamCount(ctx, event.TrackID, affectedSegIdx); err != nil {
		return err
	}

	if totalDuration > MinSeconds {
		if err = s.trackRepo.IncrementTrackStreamCount(ctx, event.TrackID, event.UserID); err != nil {
			return err
		}
	}

	return nil
}

func (ListeningStatService) proccessListeningEvent(segments []*entity.Segment, event *entity.ListeningEvent) ([]int, int) {
	totalDuration := 0
	segLength := segments[0].Range.Len()
	affectedSegIdx := make([]int, 0, len(segments))

	incrementStreamCount := func(segIdx int, lisRange *entity.Range) {
		if segIdx < 0 || segIdx >= len(segments) {
			return
		}

		intersec := intersection(segments[segIdx].Range, lisRange)
		if float64(intersec.Len()) >= float64(segments[segIdx].Range.Len())/2 {
			segments[segIdx].StreamCount++
			affectedSegIdx = append(affectedSegIdx, segIdx)
		}
	}

	for _, listenedRange := range event.Ranges {
		totalDuration += listenedRange.End - listenedRange.Start
		segStartIdx, segEndIdx := listenedRange.Start/segLength, listenedRange.End/segLength

		incrementStreamCount(segStartIdx, listenedRange)
		incrementStreamCount(segEndIdx, listenedRange)

		for idx := segStartIdx + 1; idx < segEndIdx; idx++ {
			if idx < 0 || idx >= len(segments) {
				continue
			}
			segments[idx].StreamCount++
			affectedSegIdx = append(affectedSegIdx, idx)
		}
	}

	return affectedSegIdx, totalDuration
}

func intersection(r1 *entity.Range, r2 *entity.Range) *entity.Range {
	return &entity.Range{
		Start: max(r1.Start, r2.Start),
		End:   min(r1.End, r2.End),
	}
}
