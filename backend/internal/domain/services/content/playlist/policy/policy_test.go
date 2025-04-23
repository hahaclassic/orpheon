package policy

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
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
			repoErr: nil,
			wantErr: nil,
		},
		{
			name: "not owner, but playlist is public",
			meta: &entity.PlaylistAccessMeta{
				OwnerID:   uuid.New(),
				IsPrivate: false,
			},
			claims:  &entity.Claims{UserID: userID},
			repoErr: nil,
			wantErr: nil,
		},
		{
			name: "not owner and private",
			meta: &entity.PlaylistAccessMeta{
				OwnerID:   uuid.New(),
				IsPrivate: true,
			},
			claims:  &entity.Claims{UserID: userID},
			repoErr: nil,
			wantErr: playlist.ErrForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewPlaylistAccessRepository(t)
			mockRepo.On("GetAccessMeta", ctx, playlistID).Return(tt.meta, tt.repoErr)

			svc := New(mockRepo)
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

	t.Run("owner cannot edit", func(t *testing.T) {
		meta := &entity.PlaylistAccessMeta{
			OwnerID: userID,
		}
		mockRepo := new(mocks.PlaylistAccessRepository)
		mockRepo.On("GetAccessMeta", ctx, playlistID).Return(meta, nil)

		svc := New(mockRepo)
		err := svc.CanEdit(ctx, &entity.Claims{UserID: userID}, playlistID)

		assert.ErrorIs(t, err, playlist.ErrForbidden)
	})

	t.Run("non-owner can edit", func(t *testing.T) {
		meta := &entity.PlaylistAccessMeta{
			OwnerID: uuid.New(),
		}
		mockRepo := mocks.NewPlaylistAccessRepository(t)
		mockRepo.On("GetAccessMeta", ctx, playlistID).Return(meta, nil)

		svc := New(mockRepo)
		err := svc.CanEdit(ctx, &entity.Claims{UserID: userID}, playlistID)

		assert.NoError(t, err)
	})
}

func TestCanDelete(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	playlistID := uuid.New()

	t.Run("owner can delete", func(t *testing.T) {
		meta := &entity.PlaylistAccessMeta{
			OwnerID:   userID,
			IsPrivate: true,
		}
		mockRepo := mocks.NewPlaylistAccessRepository(t)
		mockRepo.On("GetAccessMeta", ctx, playlistID).Return(meta, nil)

		svc := New(mockRepo)
		err := svc.CanDelete(ctx, &entity.Claims{UserID: userID}, playlistID)

		assert.ErrorIs(t, err, nil)
	})

	t.Run("admin can delete public", func(t *testing.T) {
		meta := &entity.PlaylistAccessMeta{
			OwnerID:   uuid.New(),
			IsPrivate: false,
		}
		mockRepo := mocks.NewPlaylistAccessRepository(t)
		mockRepo.On("GetAccessMeta", ctx, playlistID).Return(meta, nil)

		svc := New(mockRepo)
		err := svc.CanDelete(ctx, &entity.Claims{
			UserID:    userID,
			AccessLvl: entity.Admin,
		}, playlistID)

		assert.ErrorIs(t, err, playlist.ErrForbidden)
	})

	t.Run("non-admin non-owner can't delete", func(t *testing.T) {
		meta := &entity.PlaylistAccessMeta{
			OwnerID:   uuid.New(),
			IsPrivate: true,
		}
		mockRepo := mocks.NewPlaylistAccessRepository(t)
		mockRepo.On("GetAccessMeta", ctx, playlistID).Return(meta, nil)

		svc := New(mockRepo)
		err := svc.CanDelete(ctx, &entity.Claims{
			UserID:    userID,
			AccessLvl: entity.User,
		}, playlistID)

		assert.NoError(t, err)
	})
}
