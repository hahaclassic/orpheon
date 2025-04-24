package policy_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/content/playlist/policy"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
	"github.com/hahaclassic/orpheon/backend/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCanView(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	playlistID := uuid.New()

	tests := []struct {
		name    string
		meta    *entity.PlaylistAccessMeta
		claims  *entity.Claims
		wantErr error
		repoErr error
	}{
		{
			name: "owner can view private",
			meta: &entity.PlaylistAccessMeta{
				OwnerID:   userID,
				IsPrivate: true,
			},
			claims:  &entity.Claims{UserID: userID},
			wantErr: nil,
		},
		{
			name: "not owner, playlist is public",
			meta: &entity.PlaylistAccessMeta{
				OwnerID:   uuid.New(),
				IsPrivate: false,
			},
			claims:  &entity.Claims{UserID: userID},
			wantErr: nil,
		},
		{
			name: "not owner and private",
			meta: &entity.PlaylistAccessMeta{
				OwnerID:   uuid.New(),
				IsPrivate: true,
			},
			claims:  &entity.Claims{UserID: userID},
			wantErr: playlist.ErrForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewPlaylistAccessRepository(t)
			mockRepo.On("GetAccessMeta", ctx, playlistID).Return(tt.meta, tt.repoErr)

			svc := policy.New(mockRepo)
			err := svc.CanView(ctx, tt.claims, playlistID)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCanEdit(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	playlistID := uuid.New()

	tests := []struct {
		name    string
		meta    *entity.PlaylistAccessMeta
		claims  *entity.Claims
		wantErr error
	}{
		{
			name:    "owner can edit",
			meta:    &entity.PlaylistAccessMeta{OwnerID: userID},
			claims:  &entity.Claims{UserID: userID},
			wantErr: nil,
		},
		{
			name:    "not owner cannot edit",
			meta:    &entity.PlaylistAccessMeta{OwnerID: uuid.New()},
			claims:  &entity.Claims{UserID: userID},
			wantErr: playlist.ErrForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewPlaylistAccessRepository(t)
			mockRepo.On("GetAccessMeta", ctx, playlistID).Return(tt.meta, nil)

			svc := policy.New(mockRepo)
			err := svc.CanEdit(ctx, tt.claims, playlistID)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCanDelete(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	playlistID := uuid.New()

	tests := []struct {
		name    string
		meta    *entity.PlaylistAccessMeta
		claims  *entity.Claims
		wantErr error
	}{
		{
			name: "owner can delete",
			meta: &entity.PlaylistAccessMeta{
				OwnerID:   userID,
				IsPrivate: true,
			},
			claims:  &entity.Claims{UserID: userID},
			wantErr: nil,
		},
		{
			name: "admin can delete public playlist",
			meta: &entity.PlaylistAccessMeta{
				OwnerID:   uuid.New(),
				IsPrivate: false,
			},
			claims:  &entity.Claims{UserID: uuid.New(), AccessLvl: entity.Admin},
			wantErr: nil,
		},
		{
			name: "admin cannot delete private playlist",
			meta: &entity.PlaylistAccessMeta{
				OwnerID:   uuid.New(),
				IsPrivate: true,
			},
			claims:  &entity.Claims{UserID: uuid.New(), AccessLvl: entity.Admin},
			wantErr: playlist.ErrForbidden,
		},
		{
			name: "user cannot delete others' playlist",
			meta: &entity.PlaylistAccessMeta{
				OwnerID:   uuid.New(),
				IsPrivate: false,
			},
			claims:  &entity.Claims{UserID: userID},
			wantErr: playlist.ErrForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewPlaylistAccessRepository(t)
			mockRepo.On("GetAccessMeta", ctx, playlistID).Return(tt.meta, nil)

			svc := policy.New(mockRepo)
			err := svc.CanDelete(ctx, tt.claims, playlistID)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
