package audio_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/content/track/audio"
	"github.com/hahaclassic/orpheon/backend/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAudioFileService_GetAudioChunk(t *testing.T) {
	tests := []struct {
		name    string
		chunk   *entity.AudioChunk
		mock    func(repo *mocks.AudioFileRepository)
		wantErr bool
	}{
		{
			name:  "valid chunk",
			chunk: &entity.AudioChunk{TrackID: uuid.New(), Start: 0, End: 10},
			mock: func(repo *mocks.AudioFileRepository) {
				repo.On("GetAudioChunk", mock.Anything, mock.Anything).Return(&entity.AudioChunk{Data: []byte("data")}, nil)
			},
		},
		{
			name:    "invalid chunk params",
			chunk:   &entity.AudioChunk{Start: 10, End: 5},
			mock:    func(repo *mocks.AudioFileRepository) {},
			wantErr: true,
		},
		{
			name:  "repo error",
			chunk: &entity.AudioChunk{Start: 0, End: 10},
			mock: func(repo *mocks.AudioFileRepository) {
				repo.On("GetAudioChunk", mock.Anything, mock.Anything).Return(nil, errors.New("repo error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewAudioFileRepository(t)
			converter := mocks.NewAudioConverter(t)
			tt.mock(repo)

			service := audio.New(repo, converter)

			_, err := service.GetAudioChunk(context.Background(), tt.chunk)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAudioFileService_UploadAudioFile(t *testing.T) {
	admin := &entity.Claims{AccessLvl: entity.Admin}
	user := &entity.Claims{AccessLvl: entity.User}

	tests := []struct {
		name    string
		claims  *entity.Claims
		chunk   *entity.AudioChunk
		mock    func(repo *mocks.AudioFileRepository, conv *mocks.AudioConverter)
		wantErr bool
	}{
		{
			name:   "valid upload",
			claims: admin,
			chunk:  &entity.AudioChunk{TrackID: uuid.New(), Data: []byte("testdata"), Start: 0, End: 8},
			mock: func(repo *mocks.AudioFileRepository, conv *mocks.AudioConverter) {
				conv.On("ChangeBitrate", mock.Anything, mock.Anything).Return(&entity.AudioChunk{Data: []byte("converted")}, nil)
				repo.On("UploadAudioFile", mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			name:    "not admin",
			claims:  user,
			chunk:   &entity.AudioChunk{TrackID: uuid.New(), Start: 0, End: 5},
			mock:    func(repo *mocks.AudioFileRepository, conv *mocks.AudioConverter) {},
			wantErr: true,
		},
		{
			name:    "invalid chunk params",
			claims:  admin,
			chunk:   &entity.AudioChunk{TrackID: uuid.New(), Start: 1, End: 5},
			mock:    func(repo *mocks.AudioFileRepository, conv *mocks.AudioConverter) {},
			wantErr: true,
		},
		{
			name:   "converter error",
			claims: admin,
			chunk:  &entity.AudioChunk{TrackID: uuid.New(), Data: []byte("testdata"), Start: 0, End: 8},
			mock: func(repo *mocks.AudioFileRepository, conv *mocks.AudioConverter) {
				conv.On("ChangeBitrate", mock.Anything, mock.Anything).Return(nil, errors.New("convert error"))
			},
			wantErr: true,
		},
		{
			name:   "repo upload error",
			claims: admin,
			chunk:  &entity.AudioChunk{TrackID: uuid.New(), Data: []byte("testdata"), Start: 0, End: 8},
			mock: func(repo *mocks.AudioFileRepository, conv *mocks.AudioConverter) {
				conv.On("ChangeBitrate", mock.Anything, mock.Anything).Return(&entity.AudioChunk{}, nil)
				repo.On("UploadAudioFile", mock.Anything, mock.Anything).Return(errors.New("upload error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewAudioFileRepository(t)
			converter := mocks.NewAudioConverter(t)
			tt.mock(repo, converter)

			service := audio.New(repo, converter)

			err := service.UploadAudioFile(context.Background(), tt.claims, tt.chunk)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAudioFileService_DeleteAudioFile(t *testing.T) {
	admin := &entity.Claims{AccessLvl: entity.Admin}
	user := &entity.Claims{AccessLvl: entity.User}
	trackID := uuid.New()

	tests := []struct {
		name    string
		claims  *entity.Claims
		mock    func(repo *mocks.AudioFileRepository)
		wantErr bool
	}{
		{
			name:   "valid delete",
			claims: admin,
			mock: func(repo *mocks.AudioFileRepository) {
				repo.On("DeleteFile", mock.Anything, trackID).Return(nil)
			},
		},
		{
			name:    "not admin",
			claims:  user,
			mock:    func(repo *mocks.AudioFileRepository) {},
			wantErr: true,
		},
		{
			name:   "repo error",
			claims: admin,
			mock: func(repo *mocks.AudioFileRepository) {
				repo.On("DeleteFile", mock.Anything, trackID).Return(errors.New("repo error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewAudioFileRepository(t)
			converter := mocks.NewAudioConverter(t)
			tt.mock(repo)

			service := audio.New(repo, converter)

			err := service.DeleteAudioFile(context.Background(), tt.claims, trackID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
