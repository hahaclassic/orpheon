package favorites

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/hahaclassic/orpheon/backend/mocks"
)

func TestAddToUserFavorites(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	playlistID := uuid.New()
	claims := &entity.Claims{UserID: userID}

	tests := []struct {
		name      string
		policyErr error
		repoErr   error
		wantErr   error
	}{
		{"success", nil, nil, nil},
		{"policy error", commonerr.ErrForbidden, nil, commonerr.ErrForbidden},
		{"repo error", nil, errors.New("repo error"), usecase.ErrAddToUserFavorites},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy := mocks.NewPlaylistPolicyService(t)
			repo := mocks.NewPlaylistFavoriteRepository(t)
			policy.On("CanView", ctx, claims, playlistID).Return(tt.policyErr)
			if tt.policyErr == nil {
				repo.On("AddToFavorites", ctx, userID, playlistID).Return(tt.repoErr)
			}
			svc := NewPlaylistFavoriteService(repo, policy)
			err := svc.AddToUserFavorites(ctx, claims, playlistID)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestGetUserFavorites(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	claims := &entity.Claims{UserID: userID}
	mockResult := []*entity.PlaylistMeta{{ID: uuid.New()}}

	tests := []struct {
		name    string
		repoRes []*entity.PlaylistMeta
		repoErr error
		wantErr error
	}{
		{"success", mockResult, nil, nil},
		{"repo error", nil, errors.New("repo error"), usecase.ErrGetUserFavorites},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewPlaylistFavoriteRepository(t)
			policy := mocks.NewPlaylistPolicyService(t)
			repo.On("GetUserFavorites", ctx, userID).Return(tt.repoRes, tt.repoErr)
			svc := NewPlaylistFavoriteService(repo, policy)
			res, err := svc.GetUserFavorites(ctx, claims)
			assert.Equal(t, tt.repoRes, res)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestDeleteFromUserFavorites(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	playlistID := uuid.New()
	claims := &entity.Claims{UserID: userID}

	tests := []struct {
		name    string
		repoErr error
		wantErr error
	}{
		{"success", nil, nil},
		{"repo error", errors.New("repo error"), usecase.ErrDeleteFromUserFavorites},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy := mocks.NewPlaylistPolicyService(t)
			repo := mocks.NewPlaylistFavoriteRepository(t)

			repo.On("DeleteFromUserFavorites", ctx, userID, playlistID).Return(tt.repoErr)

			svc := NewPlaylistFavoriteService(repo, policy)
			err := svc.DeleteFromUserFavorites(ctx, claims, playlistID)
			assert.ErrorIs(t, err, tt.wantErr)

			repo.AssertExpectations(t)
		})
	}
}

func TestGetUsersWithFavoritePlaylist(t *testing.T) {
	ctx := context.Background()
	playlistID := uuid.New()
	claims := &entity.Claims{UserID: uuid.New()}
	userList := []uuid.UUID{uuid.New(), uuid.New()}

	tests := []struct {
		name      string
		policyErr error
		repoRes   []uuid.UUID
		repoErr   error
		wantErr   error
	}{
		{"success", nil, userList, nil, nil},
		{"policy error", commonerr.ErrForbidden, nil, nil, commonerr.ErrForbidden},
		{"repo error", nil, nil, errors.New("repo error"), usecase.ErrGetUsersWithFavoritePlaylist},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy := mocks.NewPlaylistPolicyService(t)
			repo := mocks.NewPlaylistFavoriteRepository(t)
			policy.On("CanView", ctx, claims, playlistID).Return(tt.policyErr)
			if tt.policyErr == nil {
				repo.On("GetUsersWithFavoritePlaylist", ctx, playlistID, false).Return(tt.repoRes, tt.repoErr)
			}
			svc := NewPlaylistFavoriteService(repo, policy)
			res, err := svc.GetUsersWithFavoritePlaylist(ctx, claims, playlistID, true)
			assert.Equal(t, tt.repoRes, res)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestDeleteFromAllFavorites(t *testing.T) {
	ctx := context.Background()
	playlistID := uuid.New()
	claims := &entity.Claims{UserID: uuid.New()}

	tests := []struct {
		name      string
		policyErr error
		repoErr   error
		wantErr   error
	}{
		{"success", nil, nil, nil},
		{"policy error", commonerr.ErrForbidden, nil, commonerr.ErrForbidden},
		{"repo error", nil, errors.New("repo error"), usecase.ErrDeleteFromAllFavorites},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy := mocks.NewPlaylistPolicyService(t)
			repo := mocks.NewPlaylistFavoriteRepository(t)
			policy.On("CanDelete", ctx, claims, playlistID).Return(tt.policyErr)
			if tt.policyErr == nil {
				repo.On("DeleteFromAllFavorites", ctx, playlistID, true).Return(tt.repoErr)
			}
			svc := NewPlaylistFavoriteService(repo, policy)
			err := svc.DeleteFromAllFavorites(ctx, claims, playlistID, true)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestAddPlaylistToAllFavorites(t *testing.T) {
	ctx := context.Background()
	playlistID := uuid.New()
	claims := &entity.Claims{UserID: uuid.New()}
	userIDs := []uuid.UUID{uuid.New(), uuid.New()}

	tests := []struct {
		name      string
		policyErr error
		repoErr   error
		wantErr   error
	}{
		{"success", nil, nil, nil},
		{"policy error", commonerr.ErrForbidden, nil, commonerr.ErrForbidden},
		{"repo error", nil, errors.New("repo error"), usecase.ErrAddPlaylistToAllFavorites},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy := mocks.NewPlaylistPolicyService(t)
			repo := mocks.NewPlaylistFavoriteRepository(t)
			policy.On("CanDelete", ctx, claims, playlistID).Return(tt.policyErr)
			if tt.policyErr == nil {
				repo.On("RestoreAllFavorites", ctx, userIDs, playlistID).Return(tt.repoErr)
			}
			svc := NewPlaylistFavoriteService(repo, policy)
			err := svc.AddPlaylistToAllFavorites(ctx, claims, userIDs, playlistID)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
