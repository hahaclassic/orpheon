package cover

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
	"github.com/hahaclassic/orpheon/backend/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetCover(t *testing.T) {
	ctx := context.Background()
	playlistID := uuid.New()
	claims := &entity.Claims{UserID: uuid.New()}
	cover := &entity.Cover{ObjectID: playlistID, Data: []byte("some image")}

	tests := []struct {
		name      string
		policyErr error
		repoCover *entity.Cover
		repoErr   error
		wantErr   error
	}{
		{
			name:      "success",
			policyErr: nil,
			repoCover: cover,
			wantErr:   nil,
		},
		{
			name:      "policy denied",
			policyErr: playlist.ErrForbidden,
			wantErr:   playlist.ErrGetCover,
		},
		{
			name:      "repo error",
			policyErr: nil,
			repoErr:   errors.New("repo fail"),
			wantErr:   playlist.ErrGetCover,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy := mocks.NewPlaylistPolicyService(t)
			repo := mocks.NewPlaylistCoverRepository(t)

			policy.On("CanView", ctx, claims, playlistID).Return(tt.policyErr)

			if tt.policyErr == nil {
				repo.On("GetCover", ctx, playlistID).Return(tt.repoCover, tt.repoErr)
			}

			svc := New(repo, policy)

			got, err := svc.GetCover(ctx, claims, playlistID)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.repoCover, got)
			}
		})
	}
}

func TestUploadCover(t *testing.T) {
	ctx := context.Background()
	playlistID := uuid.New()
	claims := &entity.Claims{UserID: uuid.New()}
	cover := &entity.Cover{ObjectID: playlistID, Data: []byte("img")}

	tests := []struct {
		name      string
		policyErr error
		repoErr   error
		wantErr   error
	}{
		{
			name:    "success",
			wantErr: nil,
		},
		{
			name:      "policy denied",
			policyErr: playlist.ErrForbidden,
			wantErr:   playlist.ErrUploadCover,
		},
		{
			name:    "repo error",
			repoErr: errors.New("save fail"),
			wantErr: playlist.ErrUploadCover,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy := mocks.NewPlaylistPolicyService(t)
			repo := mocks.NewPlaylistCoverRepository(t)

			policy.On("CanEdit", ctx, claims, playlistID).Return(tt.policyErr)

			if tt.policyErr == nil {
				repo.On("SaveCover", ctx, cover).Return(tt.repoErr)
			}

			svc := New(repo, policy)

			err := svc.UploadCover(ctx, claims, cover)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteCover(t *testing.T) {
	ctx := context.Background()
	playlistID := uuid.New()
	claims := &entity.Claims{UserID: uuid.New()}

	tests := []struct {
		name      string
		policyErr error
		repoErr   error
		wantErr   error
	}{
		{
			name:    "success",
			wantErr: nil,
		},
		{
			name:      "policy denied",
			policyErr: playlist.ErrForbidden,
			wantErr:   playlist.ErrDeleteCover,
		},
		{
			name:    "repo error",
			repoErr: errors.New("delete fail"),
			wantErr: playlist.ErrDeleteCover,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy := mocks.NewPlaylistPolicyService(t)
			repo := mocks.NewPlaylistCoverRepository(t)

			policy.On("CanDelete", ctx, claims, playlistID).Return(tt.policyErr)

			if tt.policyErr == nil {
				repo.On("DeleteCover", ctx, playlistID).Return(tt.repoErr)
			}

			svc := New(repo, policy)

			err := svc.DeleteCover(ctx, claims, playlistID)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
