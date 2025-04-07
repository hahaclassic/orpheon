package playlistservice

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entities"
	"github.com/hahaclassic/orpheon/backend/internal/domain/mocks"
	"github.com/stretchr/testify/assert"
)

func TestPlaylistMetaService_CreatePlaylist(t *testing.T) {
	ctx := context.Background()
	mockRepo := mocks.NewPlaylistMetaRepository(t)
	mockPolicy := mocks.NewPlaylistPolicyService(t)

	service := NewPlaylistMetaService(mockRepo, mockPolicy)

	playlist := &entities.Playlist{
		ID:      uuid.New(),
		OwnerID: uuid.New(),
		Name:    "Test Playlist",
	}

	// OK
	mockRepo.On("Create", ctx, playlist).Return(nil).Once()
	err := service.CreatePlaylist(ctx, &entities.Claims{UserID: playlist.OwnerID}, playlist)
	assert.NoError(t, err)

	// ошибка от репозитория
	mockRepo.On("Create", ctx, playlist).Return(assert.AnError).Once()
	err = service.CreatePlaylist(ctx, &entities.Claims{UserID: playlist.OwnerID}, playlist)
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestPlaylistMetaService_GetPlaylist(t *testing.T) {
	ctx := context.Background()
	mockRepo := mocks.NewPlaylistMetaRepository(t)
	mockPolicy := mocks.NewPlaylistPolicyService(t)

	service := NewPlaylistMetaService(mockRepo, mockPolicy)

	playlistID := uuid.New()
	userID := uuid.New()
	expected := &entities.Playlist{ID: playlistID, OwnerID: userID}

	// OK
	mockPolicy.On("CanView", ctx, &entities.Claims{UserID: userID}, playlistID).Return(true, nil).Once()
	mockRepo.On("GetByID", ctx, playlistID).Return(expected, nil).Once()

	res, err := service.GetPlaylist(ctx, &entities.Claims{UserID: userID}, playlistID)
	assert.NoError(t, err)
	assert.Equal(t, expected, res)

	// нет прав
	mockPolicy.On("CanView", ctx, &entities.Claims{UserID: userID}, playlistID).Return(false, nil).Once()
	res, err = service.GetPlaylist(ctx, &entities.Claims{UserID: userID}, playlistID)
	assert.Error(t, err)
	assert.Nil(t, res)

	// ошибка репозитория
	mockPolicy.On("CanView", ctx, &entities.Claims{UserID: userID}, playlistID).Return(true, nil).Once()
	mockRepo.On("GetByID", ctx, playlistID).Return(nil, assert.AnError).Once()

	res, err = service.GetPlaylist(ctx, &entities.Claims{UserID: userID}, playlistID)
	assert.Error(t, err)
	assert.Nil(t, res)

	mockPolicy.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestPlaylistMetaService_GetUserPlaylists(t *testing.T) {
	ctx := context.Background()
	mockRepo := mocks.NewPlaylistMetaRepository(t)
	mockPolicy := mocks.NewPlaylistPolicyService(t)

	service := NewPlaylistMetaService(mockRepo, mockPolicy)

	userID := uuid.New()
	playlists := []*entities.Playlist{
		{ID: uuid.New(), OwnerID: userID, Name: "Public", IsPrivate: false},
		{ID: uuid.New(), OwnerID: userID, Name: "Private", IsPrivate: true},
	}

	// OK
	mockRepo.On("GetByUser", ctx, userID).Return(playlists, nil).Once()
	res, err := service.GetUserPlaylists(ctx, &entities.Claims{UserID: userID}, userID)
	assert.NoError(t, err)
	assert.Equal(t, playlists, res)

	// ошибка из репозитория
	mockRepo.On("GetByUser", ctx, userID).Return(nil, assert.AnError).Once()
	res, err = service.GetUserPlaylists(ctx, &entities.Claims{UserID: userID}, userID)
	assert.Error(t, err)
	assert.Nil(t, res)

	mockRepo.AssertExpectations(t)
}

func TestPlaylistMetaService_UpdatePlaylist(t *testing.T) {
	ctx := context.Background()
	mockRepo := mocks.NewPlaylistMetaRepository(t)
	mockPolicy := mocks.NewPlaylistPolicyService(t)

	service := NewPlaylistMetaService(mockRepo, mockPolicy)

	playlist := &entities.Playlist{
		ID:      uuid.New(),
		OwnerID: uuid.New(),
		Name:    "Updated",
	}

	// OK
	mockPolicy.On("CanEdit", ctx, &entities.Claims{UserID: playlist.OwnerID}, playlist.ID).Return(true, nil).Once()
	mockRepo.On("Update", ctx, playlist).Return(nil).Once()

	err := service.UpdatePlaylist(ctx, &entities.Claims{UserID: playlist.OwnerID}, playlist)
	assert.NoError(t, err)

	// нет прав
	mockPolicy.On("CanEdit", ctx, &entities.Claims{UserID: playlist.OwnerID}, playlist.ID).Return(false, nil).Once()
	err = service.UpdatePlaylist(ctx, &entities.Claims{UserID: playlist.OwnerID}, playlist)
	assert.Error(t, err)

	// ошибка из репозитория
	mockPolicy.On("CanEdit", ctx, &entities.Claims{UserID: playlist.OwnerID}, playlist.ID).Return(true, nil).Once()
	mockRepo.On("Update", ctx, playlist).Return(assert.AnError).Once()
	err = service.UpdatePlaylist(ctx, &entities.Claims{UserID: playlist.OwnerID}, playlist)
	assert.Error(t, err)

	mockPolicy.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestPlaylistMetaService_DeletePlaylist(t *testing.T) {
	ctx := context.Background()
	mockRepo := mocks.NewPlaylistMetaRepository(t)
	mockPolicy := mocks.NewPlaylistPolicyService(t)

	service := NewPlaylistMetaService(mockRepo, mockPolicy)

	playlistID := uuid.New()
	userID := uuid.New()

	// OK
	mockPolicy.On("CanDelete", ctx, &entities.Claims{UserID: userID}, playlistID).Return(true, nil).Once()
	mockRepo.On("Delete", ctx, playlistID).Return(nil).Once()

	err := service.DeletePlaylist(ctx, &entities.Claims{UserID: userID}, playlistID)
	assert.NoError(t, err)

	// нет прав
	mockPolicy.On("CanDelete", ctx, &entities.Claims{UserID: userID}, playlistID).Return(false, nil).Once()
	err = service.DeletePlaylist(ctx, &entities.Claims{UserID: userID}, playlistID)
	assert.Error(t, err)

	// ошибка из репозитория
	mockPolicy.On("CanDelete", ctx, &entities.Claims{UserID: userID}, playlistID).Return(true, nil).Once()
	mockRepo.On("Delete", ctx, playlistID).Return(assert.AnError).Once()
	err = service.DeletePlaylist(ctx, &entities.Claims{UserID: userID}, playlistID)
	assert.Error(t, err)

	mockPolicy.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}
