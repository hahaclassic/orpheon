package assign_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/content/artist/assign"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/hahaclassic/orpheon/backend/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestArtistAssignService_AssignArtistToTrack(t *testing.T) {
	adminClaims := &entity.Claims{AccessLvl: entity.Admin}
	userClaims := &entity.Claims{AccessLvl: entity.User}
	artistID := uuid.New()
	trackID := uuid.New()

	tests := []struct {
		name      string
		claims    *entity.Claims
		mockFunc  func(repo *mocks.ArtistAssignRepository)
		wantErr   bool
		targetErr error
	}{
		{
			name:   "success",
			claims: adminClaims,
			mockFunc: func(repo *mocks.ArtistAssignRepository) {
				repo.On("AssignArtistToTrack", mock.Anything, artistID, trackID).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "forbidden",
			claims: userClaims,
			mockFunc: func(repo *mocks.ArtistAssignRepository) {
				// не вызывается
			},
			wantErr:   true,
			targetErr: commonerr.ErrForbidden,
		},
		{
			name:   "repo error",
			claims: adminClaims,
			mockFunc: func(repo *mocks.ArtistAssignRepository) {
				repo.On("AssignArtistToTrack", mock.Anything, artistID, trackID).
					Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewArtistAssignRepository(t)
			if tt.mockFunc != nil {
				tt.mockFunc(repo)
			}

			s := assign.NewArtistAssignService(repo)
			err := s.AssignArtistToTrack(context.Background(), tt.claims, artistID, trackID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.targetErr != nil {
					assert.ErrorIs(t, err, tt.targetErr)
				}
			} else {
				assert.NoError(t, err)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestArtistAssignService_AssignArtistToAlbum(t *testing.T) {
	adminClaims := &entity.Claims{AccessLvl: entity.Admin}
	userClaims := &entity.Claims{AccessLvl: entity.User}
	artistID := uuid.New()
	albumID := uuid.New()

	tests := []struct {
		name      string
		claims    *entity.Claims
		mockFunc  func(repo *mocks.ArtistAssignRepository)
		wantErr   bool
		targetErr error
	}{
		{
			name:   "success",
			claims: adminClaims,
			mockFunc: func(repo *mocks.ArtistAssignRepository) {
				repo.On("AssignArtistToAlbum", mock.Anything, artistID, albumID).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "forbidden",
			claims: userClaims,
			mockFunc: func(repo *mocks.ArtistAssignRepository) {
				// не вызывается
			},
			wantErr:   true,
			targetErr: commonerr.ErrForbidden,
		},
		{
			name:   "repo error",
			claims: adminClaims,
			mockFunc: func(repo *mocks.ArtistAssignRepository) {
				repo.On("AssignArtistToAlbum", mock.Anything, artistID, albumID).
					Return(errors.New("repo error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewArtistAssignRepository(t)
			if tt.mockFunc != nil {
				tt.mockFunc(repo)
			}

			s := assign.NewArtistAssignService(repo)
			err := s.AssignArtistToAlbum(context.Background(), tt.claims, artistID, albumID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.targetErr != nil {
					assert.ErrorIs(t, err, tt.targetErr)
				}
			} else {
				assert.NoError(t, err)
			}
			repo.AssertExpectations(t)
		})
	}
}
