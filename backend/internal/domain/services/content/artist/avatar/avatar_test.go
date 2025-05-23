package avatar_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/content/artist/avatar"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/artist"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/hahaclassic/orpheon/backend/mocks"
)

func TestGetCover(t *testing.T) {
	ctx := context.Background()
	artistID := uuid.New()
	expected := &entity.Cover{}

	t.Run("success", func(t *testing.T) {
		repo := mocks.NewArtistAvatarRepository(t)
		repo.On("GetCover", ctx, artistID).Return(expected, nil)

		svc := avatar.NewArtistCoverService(repo)
		got, err := svc.GetCover(ctx, artistID)

		assert.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := mocks.NewArtistAvatarRepository(t)
		repo.On("GetCover", ctx, artistID).Return(nil, errors.New("repo error"))

		svc := avatar.NewArtistCoverService(repo)
		_, err := svc.GetCover(ctx, artistID)

		assert.ErrorIs(t, err, usecase.ErrGetAvatar)
	})
}

func TestUploadCover(t *testing.T) {
	ctx := context.Background()
	cover := &entity.Cover{}
	admin := &entity.Claims{AccessLvl: entity.Admin}
	user := &entity.Claims{AccessLvl: entity.User}

	t.Run("success", func(t *testing.T) {
		repo := mocks.NewArtistAvatarRepository(t)
		repo.On("SaveCover", ctx, cover).Return(nil)

		svc := avatar.NewArtistCoverService(repo)
		err := svc.UploadCover(ctx, admin, cover)

		assert.NoError(t, err)
	})

	t.Run("forbidden", func(t *testing.T) {
		repo := mocks.NewArtistAvatarRepository(t)

		svc := avatar.NewArtistCoverService(repo)
		err := svc.UploadCover(ctx, user, cover)

		assert.ErrorIs(t, err, commonerr.ErrForbidden)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := mocks.NewArtistAvatarRepository(t)
		repo.On("SaveCover", ctx, cover).Return(errors.New("db error"))

		svc := avatar.NewArtistCoverService(repo)
		err := svc.UploadCover(ctx, admin, cover)

		assert.ErrorIs(t, err, usecase.ErrUploadAvatar)
	})
}

func TestDeleteCover(t *testing.T) {
	ctx := context.Background()
	artistID := uuid.New()
	admin := &entity.Claims{AccessLvl: entity.Admin}
	user := &entity.Claims{AccessLvl: entity.User}

	t.Run("success", func(t *testing.T) {
		repo := mocks.NewArtistAvatarRepository(t)
		repo.On("DeleteCover", ctx, artistID).Return(nil)

		svc := avatar.NewArtistCoverService(repo)
		err := svc.DeleteCover(ctx, admin, artistID)

		assert.NoError(t, err)
	})

	t.Run("forbidden", func(t *testing.T) {
		repo := mocks.NewArtistAvatarRepository(t)

		svc := avatar.NewArtistCoverService(repo)
		err := svc.DeleteCover(ctx, user, artistID)

		assert.ErrorIs(t, err, commonerr.ErrForbidden)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := mocks.NewArtistAvatarRepository(t)
		repo.On("DeleteCover", ctx, artistID).Return(errors.New("repo error"))

		svc := avatar.NewArtistCoverService(repo)
		err := svc.DeleteCover(ctx, admin, artistID)

		assert.ErrorIs(t, err, usecase.ErrDeleteAvatar)
	})
}
