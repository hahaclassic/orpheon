package meta_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/content/album/meta"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/album"
	"github.com/hahaclassic/orpheon/backend/mocks"
)

func TestCreateAlbum(t *testing.T) {
	ctx := context.Background()
	album := &entity.AlbumMeta{}
	adminClaims := &entity.Claims{AccessLvl: entity.Admin}
	userClaims := &entity.Claims{AccessLvl: entity.User}

	t.Run("success", func(t *testing.T) {
		repo := mocks.NewAlbumRepository(t)
		repo.On("CreateAlbum", ctx, mock.AnythingOfType("*entity.AlbumMeta")).Return(nil)
		svc := meta.New(repo)
		err := svc.CreateAlbum(ctx, adminClaims, album)
		assert.NoError(t, err)
	})

	t.Run("forbidden", func(t *testing.T) {
		repo := mocks.NewAlbumRepository(t)
		svc := meta.New(repo)
		err := svc.CreateAlbum(ctx, userClaims, album)
		assert.ErrorIs(t, err, meta.ErrForbidden)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := mocks.NewAlbumRepository(t)
		repo.On("CreateAlbum", ctx, mock.AnythingOfType("*entity.AlbumMeta")).Return(errors.New("db error"))
		svc := meta.New(repo)
		err := svc.CreateAlbum(ctx, adminClaims, album)
		assert.ErrorIs(t, err, usecase.ErrCreateAlbum)
	})
}

func TestGetAlbum(t *testing.T) {
	ctx := context.Background()
	albumID := uuid.New()
	album := &entity.AlbumMeta{ID: albumID}

	t.Run("success", func(t *testing.T) {
		repo := mocks.NewAlbumRepository(t)
		repo.On("GetAlbum", ctx, albumID).Return(album, nil)
		svc := meta.New(repo)
		res, err := svc.GetAlbum(ctx, albumID)
		assert.NoError(t, err)
		assert.Equal(t, album, res)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := mocks.NewAlbumRepository(t)
		repo.On("GetAlbum", ctx, albumID).Return(nil, errors.New("repo error"))
		svc := meta.New(repo)
		_, err := svc.GetAlbum(ctx, albumID)
		assert.ErrorIs(t, err, usecase.ErrGetAlbum)
	})
}

func TestUpdateAlbum(t *testing.T) {
	ctx := context.Background()
	album := &entity.AlbumMeta{}
	adminClaims := &entity.Claims{AccessLvl: entity.Admin}
	userClaims := &entity.Claims{AccessLvl: entity.User}

	t.Run("success", func(t *testing.T) {
		repo := mocks.NewAlbumRepository(t)
		repo.On("UpdateAlbum", ctx, album).Return(nil)
		svc := meta.New(repo)
		err := svc.UpdateAlbum(ctx, adminClaims, album)
		assert.NoError(t, err)
	})

	t.Run("forbidden", func(t *testing.T) {
		repo := mocks.NewAlbumRepository(t)
		svc := meta.New(repo)
		err := svc.UpdateAlbum(ctx, userClaims, album)
		assert.ErrorIs(t, err, meta.ErrForbidden)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := mocks.NewAlbumRepository(t)
		repo.On("UpdateAlbum", ctx, album).Return(errors.New("repo error"))
		svc := meta.New(repo)
		err := svc.UpdateAlbum(ctx, adminClaims, album)
		assert.ErrorIs(t, err, usecase.ErrUpdateAlbum)
	})
}

func TestDeleteAlbum(t *testing.T) {
	ctx := context.Background()
	albumID := uuid.New()
	adminClaims := &entity.Claims{AccessLvl: entity.Admin}
	userClaims := &entity.Claims{AccessLvl: entity.User}

	t.Run("success", func(t *testing.T) {
		repo := mocks.NewAlbumRepository(t)
		repo.On("DeleteAlbum", ctx, albumID).Return(nil)
		svc := meta.New(repo)
		err := svc.DeleteAlbum(ctx, adminClaims, albumID)
		assert.NoError(t, err)
	})

	t.Run("forbidden", func(t *testing.T) {
		repo := mocks.NewAlbumRepository(t)
		svc := meta.New(repo)
		err := svc.DeleteAlbum(ctx, userClaims, albumID)
		assert.ErrorIs(t, err, meta.ErrForbidden)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := mocks.NewAlbumRepository(t)
		repo.On("DeleteAlbum", ctx, albumID).Return(errors.New("repo error"))
		svc := meta.New(repo)
		err := svc.DeleteAlbum(ctx, adminClaims, albumID)
		assert.ErrorIs(t, err, usecase.ErrDeleteAlbum)
	})
}
