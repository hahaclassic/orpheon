package processor

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/mocks"
)

func TestListeningStatService_UpdateStat(t *testing.T) {
	tests := []struct {
		name          string
		event         *entity.ListeningEvent
		setupMock     func(trackRepo *mocks.TrackStatRepository, segmentRepo *mocks.SegmentStatRepository, event *entity.ListeningEvent)
		expectedError bool
	}{
		{
			name: "successfully update stat with track increment",
			event: &entity.ListeningEvent{
				TrackID: uuid.New(),
				UserID:  uuid.New(),
				Ranges: []*entity.Range{
					{Start: 0, End: 35}, // > 30 sec
				},
			},
			setupMock: func(trackRepo *mocks.TrackStatRepository, segmentRepo *mocks.SegmentStatRepository, event *entity.ListeningEvent) {
				segments := []*entity.Segment{
					{Range: &entity.Range{Start: 0, End: 10}},
					{Range: &entity.Range{Start: 10, End: 20}},
					{Range: &entity.Range{Start: 20, End: 30}},
					{Range: &entity.Range{Start: 30, End: 40}},
				}
				segmentRepo.On("GetSegments", mock.Anything, event.TrackID).Return(segments, nil).Once()
				segmentRepo.On("IncrementTotalStreams", mock.Anything, event.TrackID, mock.Anything).Return(nil).Once()
				trackRepo.On("IncrementTrackTotalStreams", mock.Anything, event.TrackID).Return(nil).Once()
			},
			expectedError: false,
		},
		{
			name: "successfully update stat without track increment",
			event: &entity.ListeningEvent{
				TrackID: uuid.New(),
				UserID:  uuid.New(),
				Ranges: []*entity.Range{
					{Start: 0, End: 20}, // < 30 sec
				},
			},
			setupMock: func(trackRepo *mocks.TrackStatRepository, segmentRepo *mocks.SegmentStatRepository, event *entity.ListeningEvent) {
				segments := []*entity.Segment{
					{Range: &entity.Range{Start: 0, End: 10}},
					{Range: &entity.Range{Start: 10, End: 20}},
					{Range: &entity.Range{Start: 20, End: 30}},
				}
				segmentRepo.On("GetSegments", mock.Anything, event.TrackID).Return(segments, nil).Once()
				segmentRepo.On("IncrementTotalStreams", mock.Anything, event.TrackID, mock.Anything).Return(nil).Once()
			},
			expectedError: false,
		},
		{
			name: "error on getting segments",
			event: &entity.ListeningEvent{
				TrackID: uuid.New(),
				UserID:  uuid.New(),
			},
			setupMock: func(trackRepo *mocks.TrackStatRepository, segmentRepo *mocks.SegmentStatRepository, event *entity.ListeningEvent) {
				segmentRepo.On("GetSegments", mock.Anything, event.TrackID).Return(nil, errors.New("db error")).Once()
			},
			expectedError: true,
		},
		{
			name: "error on incrementing segments",
			event: &entity.ListeningEvent{
				TrackID: uuid.New(),
				UserID:  uuid.New(),
				Ranges: []*entity.Range{
					{Start: 0, End: 35},
				},
			},
			setupMock: func(trackRepo *mocks.TrackStatRepository, segmentRepo *mocks.SegmentStatRepository, event *entity.ListeningEvent) {
				segments := []*entity.Segment{
					{Range: &entity.Range{Start: 0, End: 10}},
					{Range: &entity.Range{Start: 10, End: 20}},
					{Range: &entity.Range{Start: 20, End: 30}},
					{Range: &entity.Range{Start: 30, End: 40}},
				}
				segmentRepo.On("GetSegments", mock.Anything, event.TrackID).Return(segments, nil).Once()
				segmentRepo.On("IncrementTotalStreams", mock.Anything, event.TrackID, mock.Anything).Return(errors.New("increment error")).Once()
			},
			expectedError: true,
		},
		{
			name: "error on incrementing track plays",
			event: &entity.ListeningEvent{
				TrackID: uuid.New(),
				UserID:  uuid.New(),
				Ranges: []*entity.Range{
					{Start: 0, End: 40},
				},
			},
			setupMock: func(trackRepo *mocks.TrackStatRepository, segmentRepo *mocks.SegmentStatRepository, event *entity.ListeningEvent) {
				segments := []*entity.Segment{
					{Range: &entity.Range{Start: 0, End: 10}},
					{Range: &entity.Range{Start: 10, End: 20}},
					{Range: &entity.Range{Start: 20, End: 30}},
					{Range: &entity.Range{Start: 30, End: 40}},
				}
				segmentRepo.On("GetSegments", mock.Anything, event.TrackID).Return(segments, nil).Once()
				segmentRepo.On("IncrementTotalStreams", mock.Anything, event.TrackID, mock.Anything).Return(nil).Once()
				trackRepo.On("IncrementTrackTotalStreams", mock.Anything, event.TrackID).Return(errors.New("increment track plays error")).Once()
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trackRepo := mocks.NewTrackStatRepository(t)
			segmentRepo := mocks.NewSegmentStatRepository(t)

			if tt.setupMock != nil {
				tt.setupMock(trackRepo, segmentRepo, tt.event)
			}

			service := NewListeningStatService(trackRepo, segmentRepo)
			err := service.UpdateStat(context.Background(), tt.event)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
