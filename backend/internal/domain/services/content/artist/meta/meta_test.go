package meta_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/content/artist/meta"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/artist"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/hahaclassic/orpheon/backend/mocks"
)

func TestGetArtistMeta(t *testing.T) {
	ctx := context.Background()
	artistID := uuid.New()
	expected := &entity.ArtistMeta{ID: artistID}

	t.Run("success", func(t *testing.T) {
		repo := mocks.NewArtistMetaRepository(t)
		repo.On("GetByID", ctx, artistID).Return(expected, nil)

		svc := meta.New(repo)
		got, err := svc.GetArtistMeta(ctx, artistID)

		assert.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := mocks.NewArtistMetaRepository(t)
		repo.On("GetByID", ctx, artistID).Return(nil, errors.New("not found"))

		svc := meta.New(repo)
		got, err := svc.GetArtistMeta(ctx, artistID)

		assert.Nil(t, got)
		assert.Error(t, err)
	})
}

func TestCreateArtistMeta(t *testing.T) {
	ctx := context.Background()
	artist := &entity.ArtistMeta{}
	user := &entity.Claims{AccessLvl: entity.User}
	admin := &entity.Claims{AccessLvl: entity.Admin}

	t.Run("success", func(t *testing.T) {
		repo := mocks.NewArtistMetaRepository(t)
		repo.On("Create", ctx, mock.MatchedBy(func(a *entity.ArtistMeta) bool {
			return a.ID != uuid.Nil
		})).Return(nil)

		svc := meta.New(repo)
		err := svc.CreateArtistMeta(ctx, admin, artist)

		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, artist.ID)
	})

	t.Run("forbidden for user", func(t *testing.T) {
		repo := mocks.NewArtistMetaRepository(t)

		svc := meta.New(repo)
		err := svc.CreateArtistMeta(ctx, user, artist)

		assert.ErrorIs(t, err, commonerr.ErrForbidden)
	})
}

func TestUpdateArtistMeta(t *testing.T) {
	ctx := context.Background()
	artist := &entity.ArtistMeta{}
	user := &entity.Claims{AccessLvl: entity.User}
	admin := &entity.Claims{AccessLvl: entity.Admin}

	t.Run("success", func(t *testing.T) {
		repo := mocks.NewArtistMetaRepository(t)
		repo.On("Update", ctx, artist).Return(nil)

		svc := meta.New(repo)
		err := svc.UpdateArtistMeta(ctx, admin, artist)

		assert.NoError(t, err)
	})

	t.Run("forbidden for user", func(t *testing.T) {
		repo := mocks.NewArtistMetaRepository(t)

		svc := meta.New(repo)
		err := svc.UpdateArtistMeta(ctx, user, artist)

		assert.ErrorIs(t, err, commonerr.ErrForbidden)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := mocks.NewArtistMetaRepository(t)
		repo.On("Update", ctx, artist).Return(errors.New("update failed"))

		svc := meta.New(repo)
		err := svc.UpdateArtistMeta(ctx, admin, artist)

		assert.ErrorIs(t, err, usecase.ErrUpdateArtistMeta)
	})
}

func TestDeleteArtistMeta(t *testing.T) {
	ctx := context.Background()
	artistID := uuid.New()
	admin := &entity.Claims{AccessLvl: entity.Admin}
	user := &entity.Claims{AccessLvl: entity.User}

	t.Run("success", func(t *testing.T) {
		repo := mocks.NewArtistMetaRepository(t)
		repo.On("Delete", ctx, artistID).Return(nil)

		svc := meta.New(repo)
		err := svc.DeleteArtistMeta(ctx, admin, artistID)

		assert.NoError(t, err)
	})

	t.Run("forbidden for user", func(t *testing.T) {
		repo := mocks.NewArtistMetaRepository(t)

		svc := meta.New(repo)
		err := svc.DeleteArtistMeta(ctx, user, artistID)

		assert.ErrorIs(t, err, commonerr.ErrForbidden)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := mocks.NewArtistMetaRepository(t)
		repo.On("Delete", ctx, artistID).Return(errors.New("delete error"))

		svc := meta.New(repo)
		err := svc.DeleteArtistMeta(ctx, admin, artistID)

		assert.ErrorIs(t, err, usecase.ErrDeleteArtistMeta)
	})
}
