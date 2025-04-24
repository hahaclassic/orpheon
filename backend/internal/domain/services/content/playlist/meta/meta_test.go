package meta_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/content/playlist/meta"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
	"github.com/hahaclassic/orpheon/backend/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateMeta(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	tests := []struct {
		name     string
		playlist *entity.PlaylistMeta
		setup    func(repo *mocks.PlaylistMetaRepository)
		wantErr  error
	}{
		{
			name:     "success",
			playlist: &entity.PlaylistMeta{Name: "My Playlist"},
			setup: func(repo *mocks.PlaylistMetaRepository) {
				repo.On("Create", ctx, mock.MatchedBy(func(p *entity.PlaylistMeta) bool {
					return p.Name == "My Playlist" && p.OwnerID == userID
				})).Return(nil)
			},
			wantErr: nil,
		},
		{
			name:     "empty name",
			playlist: &entity.PlaylistMeta{Name: ""},
			setup:    func(repo *mocks.PlaylistMetaRepository) {},
			wantErr:  meta.ErrEmptyPlaylistName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewPlaylistMetaRepository(t)
			mockPolicy := mocks.NewPlaylistPolicyService(t)
			tt.setup(mockRepo)

			svc := meta.NewPlaylistMetaService(mockRepo, mockPolicy)
			claims := &entity.Claims{UserID: userID}
			err := svc.CreateMeta(ctx, claims, tt.playlist)

			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestGetMeta(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	playlistID := uuid.New()
	playlistMeta := &entity.PlaylistMeta{ID: playlistID, Name: "Meta"}

	tests := []struct {
		name    string
		setup   func(repo *mocks.PlaylistMetaRepository, policy *mocks.PlaylistPolicyService)
		want    *entity.PlaylistMeta
		wantErr error
	}{
		{
			name: "success",
			setup: func(repo *mocks.PlaylistMetaRepository, policy *mocks.PlaylistPolicyService) {
				policy.On("CanView", ctx, mock.Anything, playlistID).Return(nil)
				repo.On("GetByID", ctx, playlistID).Return(playlistMeta, nil)
			},
			want:    playlistMeta,
			wantErr: nil,
		},
		{
			name: "forbidden",
			setup: func(repo *mocks.PlaylistMetaRepository, policy *mocks.PlaylistPolicyService) {
				policy.On("CanView", ctx, mock.Anything, playlistID).Return(playlist.ErrForbidden)
			},
			want:    nil,
			wantErr: playlist.ErrForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewPlaylistMetaRepository(t)
			mockPolicy := mocks.NewPlaylistPolicyService(t)
			tt.setup(mockRepo, mockPolicy)

			svc := meta.NewPlaylistMetaService(mockRepo, mockPolicy)
			claims := &entity.Claims{UserID: userID}
			got, err := svc.GetMeta(ctx, claims, playlistID)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetUserAllPlaylistsMeta(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherID := uuid.New()

	playlists := []*entity.PlaylistMeta{
		{ID: uuid.New(), OwnerID: userID, Name: "Public1", IsPrivate: false},
		{ID: uuid.New(), OwnerID: userID, Name: "Private", IsPrivate: true},
		{ID: uuid.New(), OwnerID: userID, Name: "Public2", IsPrivate: false},
	}

	tests := []struct {
		name      string
		claims    *entity.Claims
		queryUser uuid.UUID
		wantCount int
	}{
		{
			name:      "owner sees all",
			claims:    &entity.Claims{UserID: userID},
			queryUser: userID,
			wantCount: 3,
		},
		{
			name:      "not owner sees only public",
			claims:    &entity.Claims{UserID: otherID},
			queryUser: userID,
			wantCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewPlaylistMetaRepository(t)
			mockPolicy := mocks.NewPlaylistPolicyService(t)
			mockRepo.On("GetByUser", ctx, tt.queryUser).Return(playlists, nil)

			svc := meta.NewPlaylistMetaService(mockRepo, mockPolicy)
			got, err := svc.GetUserAllPlaylistsMeta(ctx, tt.claims, tt.queryUser)

			assert.NoError(t, err)
			assert.Len(t, got, tt.wantCount)
		})
	}
}

func TestUpdateMeta(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	playlistID := uuid.New()
	metaToUpdate := &entity.PlaylistMeta{ID: playlistID, Name: "Update"}

	tests := []struct {
		name    string
		setup   func(repo *mocks.PlaylistMetaRepository, policy *mocks.PlaylistPolicyService)
		wantErr error
	}{
		{
			name: "success",
			setup: func(repo *mocks.PlaylistMetaRepository, policy *mocks.PlaylistPolicyService) {
				policy.On("CanEdit", ctx, mock.Anything, playlistID).Return(nil)
				repo.On("Update", ctx, metaToUpdate).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "forbidden",
			setup: func(repo *mocks.PlaylistMetaRepository, policy *mocks.PlaylistPolicyService) {
				policy.On("CanEdit", ctx, mock.Anything, playlistID).Return(playlist.ErrForbidden)
			},
			wantErr: playlist.ErrForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewPlaylistMetaRepository(t)
			mockPolicy := mocks.NewPlaylistPolicyService(t)
			tt.setup(mockRepo, mockPolicy)

			svc := meta.NewPlaylistMetaService(mockRepo, mockPolicy)
			claims := &entity.Claims{UserID: userID}
			err := svc.UpdateMeta(ctx, claims, metaToUpdate)

			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestDeleteMeta(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	playlistID := uuid.New()

	tests := []struct {
		name    string
		setup   func(repo *mocks.PlaylistMetaRepository, policy *mocks.PlaylistPolicyService)
		wantErr error
	}{
		{
			name: "success",
			setup: func(repo *mocks.PlaylistMetaRepository, policy *mocks.PlaylistPolicyService) {
				policy.On("CanDelete", ctx, mock.Anything, playlistID).Return(nil)
				repo.On("Delete", ctx, playlistID).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "forbidden",
			setup: func(repo *mocks.PlaylistMetaRepository, policy *mocks.PlaylistPolicyService) {
				policy.On("CanDelete", ctx, mock.Anything, playlistID).Return(playlist.ErrForbidden)
			},
			wantErr: playlist.ErrForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewPlaylistMetaRepository(t)
			mockPolicy := mocks.NewPlaylistPolicyService(t)
			tt.setup(mockRepo, mockPolicy)

			svc := meta.NewPlaylistMetaService(mockRepo, mockPolicy)
			claims := &entity.Claims{UserID: userID}
			err := svc.DeleteMeta(ctx, claims, playlistID)

			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
