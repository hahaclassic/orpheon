package meta_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/content/track/meta"
	"github.com/hahaclassic/orpheon/backend/mocks"
)

func TestTrackMetaService_CreateTrackMeta(t *testing.T) {
	type args struct {
		claims *entity.Claims
		track  *entity.TrackMeta
	}
	tests := []struct {
		name       string
		args       args
		mockSetup  func(repo *mocks.TrackMetaRepository)
		expectsErr bool
	}{
		{
			name: "success",
			args: args{
				claims: &entity.Claims{AccessLvl: entity.Admin},
				track:  &entity.TrackMeta{},
			},
			mockSetup: func(repo *mocks.TrackMetaRepository) {
				repo.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			expectsErr: false,
		},
		{
			name: "permission denied",
			args: args{
				claims: &entity.Claims{AccessLvl: entity.User},
				track:  &entity.TrackMeta{},
			},
			mockSetup:  func(repo *mocks.TrackMetaRepository) {},
			expectsErr: true,
		},
		{
			name: "repo error",
			args: args{
				claims: &entity.Claims{AccessLvl: entity.Admin},
				track:  &entity.TrackMeta{},
			},
			mockSetup: func(repo *mocks.TrackMetaRepository) {
				repo.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			expectsErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewTrackMetaRepository(t)
			service := meta.NewTrackMetaService(repo)
			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}
			_, err := service.CreateTrackMeta(context.Background(), tt.args.claims, tt.args.track)
			if tt.expectsErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				repo.AssertCalled(t, "Create", mock.Anything, mock.Anything)
			}
		})
	}
}

func TestTrackMetaService_GetTrackMeta(t *testing.T) {
	repo := mocks.NewTrackMetaRepository(t)
	service := meta.NewTrackMetaService(repo)

	trackID := uuid.New()
	expected := &entity.TrackMeta{ID: trackID, Name: "Test Track"}
	repo.On("GetByID", mock.Anything, trackID).Return(expected, nil)

	track, err := service.GetTrackMeta(context.Background(), trackID)
	require.NoError(t, err)
	assert.Equal(t, expected, track)
}

func TestTrackMetaService_UpdateTrackMeta(t *testing.T) {
	tests := []struct {
		name       string
		claims     *entity.Claims
		track      *entity.TrackMeta
		mockSetup  func(repo *mocks.TrackMetaRepository)
		expectsErr bool
	}{
		{
			name:   "success",
			claims: &entity.Claims{AccessLvl: entity.Admin},
			track:  &entity.TrackMeta{ID: uuid.New()},
			mockSetup: func(repo *mocks.TrackMetaRepository) {
				repo.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			expectsErr: false,
		},
		{
			name:       "permission denied",
			claims:     &entity.Claims{AccessLvl: entity.User},
			track:      &entity.TrackMeta{ID: uuid.New()},
			mockSetup:  func(repo *mocks.TrackMetaRepository) {},
			expectsErr: true,
		},
		{
			name:   "repo error",
			claims: &entity.Claims{AccessLvl: entity.Admin},
			track:  &entity.TrackMeta{ID: uuid.New()},
			mockSetup: func(repo *mocks.TrackMetaRepository) {
				repo.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			expectsErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewTrackMetaRepository(t)
			service := meta.NewTrackMetaService(repo)
			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}
			err := service.UpdateTrackMeta(context.Background(), tt.claims, tt.track)
			if tt.expectsErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				repo.AssertCalled(t, "Update", mock.Anything, mock.Anything)
			}
		})
	}
}

func TestTrackMetaService_DeleteTrackMeta(t *testing.T) {
	tests := []struct {
		name       string
		claims     *entity.Claims
		trackID    uuid.UUID
		mockSetup  func(repo *mocks.TrackMetaRepository)
		expectsErr bool
	}{
		{
			name:    "success",
			claims:  &entity.Claims{AccessLvl: entity.Admin},
			trackID: uuid.New(),
			mockSetup: func(repo *mocks.TrackMetaRepository) {
				repo.On("Delete", mock.Anything, mock.Anything).Return(nil)
			},
			expectsErr: false,
		},
		{
			name:       "permission denied",
			claims:     &entity.Claims{AccessLvl: entity.User},
			trackID:    uuid.New(),
			mockSetup:  func(repo *mocks.TrackMetaRepository) {},
			expectsErr: true,
		},
		{
			name:    "repo error",
			claims:  &entity.Claims{AccessLvl: entity.Admin},
			trackID: uuid.New(),
			mockSetup: func(repo *mocks.TrackMetaRepository) {
				repo.On("Delete", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			expectsErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewTrackMetaRepository(t)
			service := meta.NewTrackMetaService(repo)
			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}
			err := service.DeleteTrackMeta(context.Background(), tt.claims, tt.trackID)
			if tt.expectsErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				repo.AssertCalled(t, "Delete", mock.Anything, mock.Anything)
			}
		})
	}
}
