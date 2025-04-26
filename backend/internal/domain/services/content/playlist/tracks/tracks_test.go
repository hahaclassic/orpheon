package tracks_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/content/playlist/tracks"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
	"github.com/hahaclassic/orpheon/backend/mocks"
	"github.com/stretchr/testify/assert"
)

func TestPlaylistTrackService_AddTrack(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	playlistID := uuid.New()
	trackID := uuid.New()
	claims := &entity.Claims{UserID: userID}

	cases := []struct {
		name      string
		setup     func(*mocks.PlaylistTracksRepository, *mocks.PlaylistPolicyService)
		expectErr error
	}{
		{
			name: "success",
			setup: func(repo *mocks.PlaylistTracksRepository, policy *mocks.PlaylistPolicyService) {
				policy.On("CanEdit", ctx, claims, playlistID).Return(nil)
				repo.On("AddTrackToPlaylist", ctx, playlistID, trackID).Return(nil)
			},
			expectErr: nil,
		},
		{
			name: "forbidden",
			setup: func(_ *mocks.PlaylistTracksRepository, policy *mocks.PlaylistPolicyService) {
				policy.On("CanEdit", ctx, claims, playlistID).Return(playlist.ErrForbidden)
			},
			expectErr: playlist.ErrForbidden,
		},
		{
			name: "repo error",
			setup: func(repo *mocks.PlaylistTracksRepository, policy *mocks.PlaylistPolicyService) {
				policy.On("CanEdit", ctx, claims, playlistID).Return(nil)
				repo.On("AddTrackToPlaylist", ctx, playlistID, trackID).Return(errors.New("db error"))
			},
			expectErr: errors.New("db error"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewPlaylistTracksRepository(t)
			policy := mocks.NewPlaylistPolicyService(t)
			tc.setup(repo, policy)

			svc := tracks.NewPlaylistTrackService(repo, policy)
			err := svc.AddTrack(ctx, claims, playlistID, trackID)

			if tc.expectErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tc.expectErr.Error())
			}
		})
	}
}

func TestPlaylistTrackService_GetAllTracks(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	playlistID := uuid.New()
	claims := &entity.Claims{UserID: userID}
	expected := []*entity.TrackMeta{{ID: uuid.New()}}

	cases := []struct {
		name         string
		setup        func(*mocks.PlaylistTracksRepository, *mocks.PlaylistPolicyService)
		expectResult []*entity.TrackMeta
		expectErr    error
	}{
		{
			name: "success",
			setup: func(repo *mocks.PlaylistTracksRepository, policy *mocks.PlaylistPolicyService) {
				policy.On("CanView", ctx, claims, playlistID).Return(nil)
				repo.On("GetAllPlaylistTracks", ctx, playlistID).Return(expected, nil)
			},
			expectResult: expected,
			expectErr:    nil,
		},
		{
			name: "forbidden",
			setup: func(_ *mocks.PlaylistTracksRepository, policy *mocks.PlaylistPolicyService) {
				policy.On("CanView", ctx, claims, playlistID).Return(playlist.ErrForbidden)
			},
			expectResult: nil,
			expectErr:    playlist.ErrForbidden,
		},
		{
			name: "repo error",
			setup: func(repo *mocks.PlaylistTracksRepository, policy *mocks.PlaylistPolicyService) {
				policy.On("CanView", ctx, claims, playlistID).Return(nil)
				repo.On("GetAllPlaylistTracks", ctx, playlistID).Return(nil, errors.New("db error"))
			},
			expectResult: nil,
			expectErr:    errors.New("db error"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewPlaylistTracksRepository(t)
			policy := mocks.NewPlaylistPolicyService(t)
			tc.setup(repo, policy)

			svc := tracks.NewPlaylistTrackService(repo, policy)
			res, err := svc.GetAllTracks(ctx, claims, playlistID)

			assert.Equal(t, tc.expectResult, res)
			if tc.expectErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tc.expectErr.Error())
			}
		})
	}
}

func TestPlaylistTrackService_DeleteTrack(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	playlistID := uuid.New()
	trackID := uuid.New()
	claims := &entity.Claims{UserID: userID}

	cases := []struct {
		name      string
		setup     func(*mocks.PlaylistTracksRepository, *mocks.PlaylistPolicyService)
		expectErr error
	}{
		{
			name: "success",
			setup: func(repo *mocks.PlaylistTracksRepository, policy *mocks.PlaylistPolicyService) {
				policy.On("CanEdit", ctx, claims, playlistID).Return(nil)
				repo.On("DeleteTrackFromPlaylist", ctx, playlistID, trackID).Return(nil)
			},
			expectErr: nil,
		},
		{
			name: "forbidden",
			setup: func(_ *mocks.PlaylistTracksRepository, policy *mocks.PlaylistPolicyService) {
				policy.On("CanEdit", ctx, claims, playlistID).Return(playlist.ErrForbidden)
			},
			expectErr: playlist.ErrForbidden,
		},
		{
			name: "repo error",
			setup: func(repo *mocks.PlaylistTracksRepository, policy *mocks.PlaylistPolicyService) {
				policy.On("CanEdit", ctx, claims, playlistID).Return(nil)
				repo.On("DeleteTrackFromPlaylist", ctx, playlistID, trackID).Return(errors.New("db error"))
			},
			expectErr: errors.New("db error"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewPlaylistTracksRepository(t)
			policy := mocks.NewPlaylistPolicyService(t)
			tc.setup(repo, policy)

			svc := tracks.NewPlaylistTrackService(repo, policy)
			err := svc.DeleteTrack(ctx, claims, playlistID, trackID)

			if tc.expectErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tc.expectErr.Error())
			}
		})
	}
}

func TestPlaylistTrackService_DeleteAllTracks(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	playlistID := uuid.New()
	claims := &entity.Claims{UserID: userID}

	cases := []struct {
		name      string
		setup     func(*mocks.PlaylistTracksRepository, *mocks.PlaylistPolicyService)
		expectErr error
	}{
		{
			name: "success",
			setup: func(repo *mocks.PlaylistTracksRepository, policy *mocks.PlaylistPolicyService) {
				policy.On("CanEdit", ctx, claims, playlistID).Return(nil)
				repo.On("DeleteAllTracksFromPlaylist", ctx, playlistID).Return(nil)
			},
			expectErr: nil,
		},
		{
			name: "forbidden",
			setup: func(_ *mocks.PlaylistTracksRepository, policy *mocks.PlaylistPolicyService) {
				policy.On("CanEdit", ctx, claims, playlistID).Return(playlist.ErrForbidden)
			},
			expectErr: playlist.ErrForbidden,
		},
		{
			name: "repo error",
			setup: func(repo *mocks.PlaylistTracksRepository, policy *mocks.PlaylistPolicyService) {
				policy.On("CanEdit", ctx, claims, playlistID).Return(nil)
				repo.On("DeleteAllTracksFromPlaylist", ctx, playlistID).Return(errors.New("db error"))
			},
			expectErr: errors.New("db error"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewPlaylistTracksRepository(t)
			policy := mocks.NewPlaylistPolicyService(t)
			tc.setup(repo, policy)

			svc := tracks.NewPlaylistTrackService(repo, policy)
			err := svc.DeleteAllTracks(ctx, claims, playlistID)

			if tc.expectErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tc.expectErr.Error())
			}
		})
	}
}
