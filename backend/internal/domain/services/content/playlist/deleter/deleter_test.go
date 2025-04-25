package deleter_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/content/playlist/deleter"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
	"github.com/hahaclassic/orpheon/backend/mocks"
)

func TestPlaylistDeleter_DeletePlaylist(t *testing.T) {
	ctx := context.Background()
	playlistID := uuid.New()
	claims := &entity.Claims{UserID: uuid.New()}
	trackIDs := []uuid.UUID{uuid.New()}
	userIDs := []uuid.UUID{uuid.New()}
	cover := &entity.Cover{ObjectID: uuid.New()}

	tests := []struct {
		name       string
		setupMocks func(
			*mock.Mock, *mock.Mock, *mock.Mock, *mock.Mock,
		)
		wantErr error
	}{
		{
			name: "success",
			setupMocks: func(meta, track, fav, coverMock *mock.Mock) {
				meta.On("DeleteMeta", ctx, claims, playlistID).Return(nil)
				track.On("GetAllTracks", ctx, claims, playlistID).Return(trackIDs, nil)
				track.On("DeleteAllTracks", ctx, claims, playlistID).Return(nil)
				fav.On("GetUsersWithFavoritePlaylist", ctx, claims, playlistID).Return(userIDs, nil)
				fav.On("DeletePlaylistFromAllFavorites", ctx, claims, playlistID).Return(nil)
				coverMock.On("GetCover", ctx, claims, playlistID).Return(cover, nil)
				coverMock.On("DeleteCover", ctx, claims, playlistID).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "error on delete favorites, triggers rollback nothing",
			setupMocks: func(meta, track, fav, coverMock *mock.Mock) {
				fav.On("GetUsersWithFavoritePlaylist", ctx, claims, playlistID).Return(nil, errors.New("fail"))
			},
			wantErr: usecase.ErrDeletePlaylist,
		},
		{
			name: "error on delete cover, triggers rollback of favorites",
			setupMocks: func(meta, track, fav, coverMock *mock.Mock) {
				fav.On("GetUsersWithFavoritePlaylist", ctx, claims, playlistID).Return(userIDs, nil)
				fav.On("DeletePlaylistFromAllFavorites", ctx, claims, playlistID).Return(nil)
				fav.On("AddPlaylistToAllFavorites", ctx, claims, userIDs, playlistID).Return(nil)

				coverMock.On("GetCover", ctx, claims, playlistID).Return(nil, errors.New("cover error"))
			},
			wantErr: usecase.ErrDeletePlaylist,
		},
		{
			name: "error on delete tracks, triggers rollback of fav and cover",
			setupMocks: func(meta, track, fav, coverMock *mock.Mock) {
				fav.On("GetUsersWithFavoritePlaylist", ctx, claims, playlistID).Return(userIDs, nil)
				fav.On("DeletePlaylistFromAllFavorites", ctx, claims, playlistID).Return(nil)
				fav.On("AddPlaylistToAllFavorites", ctx, claims, userIDs, playlistID).Return(nil)

				coverMock.On("GetCover", ctx, claims, playlistID).Return(cover, nil)
				coverMock.On("DeleteCover", ctx, claims, playlistID).Return(nil)
				coverMock.On("SaveCover", ctx, claims, cover).Return(nil)

				track.On("GetAllTracks", ctx, claims, playlistID).Return(nil, errors.New("track error"))
			},
			wantErr: usecase.ErrDeletePlaylist,
		},
		{
			name: "error on delete meta, triggers rollback of all",
			setupMocks: func(meta, track, fav, coverMock *mock.Mock) {
				fav.On("GetUsersWithFavoritePlaylist", ctx, claims, playlistID).Return(userIDs, nil)
				fav.On("DeletePlaylistFromAllFavorites", ctx, claims, playlistID).Return(nil)
				fav.On("AddPlaylistToAllFavorites", ctx, claims, userIDs, playlistID).Return(nil)

				coverMock.On("GetCover", ctx, claims, playlistID).Return(cover, nil)
				coverMock.On("DeleteCover", ctx, claims, playlistID).Return(nil)
				coverMock.On("SaveCover", ctx, claims, cover).Return(nil)

				track.On("GetAllTracks", ctx, claims, playlistID).Return(trackIDs, nil)
				track.On("DeleteAllTracks", ctx, claims, playlistID).Return(nil)
				track.On("RestoreAllTracks", ctx, claims, playlistID, trackIDs).Return(nil)

				meta.On("DeleteMeta", ctx, claims, playlistID).Return(errors.New("meta error"))
			},
			wantErr: usecase.ErrDeletePlaylist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta := mocks.NewMetaDeletionService(t)
			tracks := mocks.NewTrackDeletionService(t)
			fav := mocks.NewFavoritesDeletionService(t)
			coverSvc := mocks.NewPlaylistCoverDeletionService(t)

			tt.setupMocks(&meta.Mock, &tracks.Mock, &fav.Mock, &coverSvc.Mock)

			svc := deleter.New(
				deleter.WithMetaDeletion(meta),
				deleter.WithTracksDeletion(tracks),
				deleter.WithFavoritesDeletion(fav),
				deleter.WIthCoverDeletion(coverSvc),
			)

			err := svc.DeletePlaylist(ctx, claims, playlistID)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
