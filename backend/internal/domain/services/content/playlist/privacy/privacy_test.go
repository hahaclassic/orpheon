package privacy

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/hahaclassic/orpheon/backend/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestChangePrivacy(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	playlistID := uuid.New()

	tests := []struct {
		name      string
		isPrivate bool
		setup     func(policy *mocks.PlaylistPolicyService, favorites *mocks.FavoritesDeletionService, accessRepo *mocks.PlaylistPrivacyRepository)
		wantErr   error
	}{
		{
			name:      "success make private",
			isPrivate: true,
			setup: func(policy *mocks.PlaylistPolicyService, favorites *mocks.FavoritesDeletionService, accessRepo *mocks.PlaylistPrivacyRepository) {
				policy.On("CanEdit", ctx, mock.Anything, playlistID).Return(nil)
				favorites.On("GetUsersWithFavoritePlaylist", ctx, mock.Anything, playlistID, false).Return([]uuid.UUID{}, nil)
				favorites.On("DeleteFromAllFavorites", ctx, mock.Anything, playlistID, false).Return(nil)
				accessRepo.On("UpdatePrivacy", ctx, playlistID, true).Return(nil)
			},
			wantErr: nil,
		},
		{
			name:      "success make public",
			isPrivate: false,
			setup: func(policy *mocks.PlaylistPolicyService, favorites *mocks.FavoritesDeletionService, accessRepo *mocks.PlaylistPrivacyRepository) {
				policy.On("CanEdit", ctx, mock.Anything, playlistID).Return(nil)
				accessRepo.On("UpdatePrivacy", ctx, playlistID, false).Return(nil)
			},
			wantErr: nil,
		},
		{
			name:      "forbidden",
			isPrivate: true,
			setup: func(policy *mocks.PlaylistPolicyService, favorites *mocks.FavoritesDeletionService, accessRepo *mocks.PlaylistPrivacyRepository) {
				policy.On("CanEdit", ctx, mock.Anything, playlistID).Return(commonerr.ErrForbidden)
			},
			wantErr: commonerr.ErrForbidden,
		},
		{
			name:      "get users error",
			isPrivate: true,
			setup: func(policy *mocks.PlaylistPolicyService, favorites *mocks.FavoritesDeletionService, accessRepo *mocks.PlaylistPrivacyRepository) {
				policy.On("CanEdit", ctx, mock.Anything, playlistID).Return(nil)
				favorites.On("GetUsersWithFavoritePlaylist", ctx, mock.Anything, playlistID, false).Return(nil, assert.AnError)
			},
			wantErr: playlist.ErrChangePrivacy,
		},
		{
			name:      "delete favorites error",
			isPrivate: true,
			setup: func(policy *mocks.PlaylistPolicyService, favorites *mocks.FavoritesDeletionService, accessRepo *mocks.PlaylistPrivacyRepository) {
				policy.On("CanEdit", ctx, mock.Anything, playlistID).Return(nil)
				favorites.On("GetUsersWithFavoritePlaylist", ctx, mock.Anything, playlistID, false).Return([]uuid.UUID{}, nil)
				favorites.On("DeleteFromAllFavorites", ctx, mock.Anything, playlistID, false).Return(assert.AnError)
			},
			wantErr: playlist.ErrChangePrivacy,
		},
		{
			name:      "update privacy error",
			isPrivate: true,
			setup: func(policy *mocks.PlaylistPolicyService, favorites *mocks.FavoritesDeletionService, accessRepo *mocks.PlaylistPrivacyRepository) {
				policy.On("CanEdit", ctx, mock.Anything, playlistID).Return(nil)
				favorites.On("GetUsersWithFavoritePlaylist", ctx, mock.Anything, playlistID, false).Return([]uuid.UUID{}, nil)
				favorites.On("DeleteFromAllFavorites", ctx, mock.Anything, playlistID, false).Return(nil)
				favorites.On("AddPlaylistToAllFavorites", ctx, mock.Anything, []uuid.UUID{}, playlistID).Return(nil)
				accessRepo.On("UpdatePrivacy", ctx, playlistID, true).Return(assert.AnError)
			},
			wantErr: playlist.ErrChangePrivacy,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPolicy := mocks.NewPlaylistPolicyService(t)
			mockFavorites := mocks.NewFavoritesDeletionService(t)
			mockAccessRepo := mocks.NewPlaylistPrivacyRepository(t)
			tt.setup(mockPolicy, mockFavorites, mockAccessRepo)

			svc := NewPlaylistPrivacyChanger(mockPolicy, mockFavorites, mockAccessRepo)
			claims := &entity.Claims{UserID: userID}
			err := svc.ChangePrivacy(ctx, claims, playlistID, tt.isPrivate)

			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
