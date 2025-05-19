package genre_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/content/genre"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/genre"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/hahaclassic/orpheon/backend/mocks"
)

func TestCreateGenre(t *testing.T) {
	ctx := context.Background()
	adminClaims := &entity.Claims{AccessLvl: entity.Admin}
	userClaims := &entity.Claims{AccessLvl: entity.User}
	validGenre := &entity.Genre{ID: uuid.New()}

	tests := []struct {
		name      string
		claims    *entity.Claims
		genre     *entity.Genre
		mockFn    func(r *mocks.GenreRepository)
		wantError error
	}{
		{
			name:   "success",
			claims: adminClaims,
			genre:  validGenre,
			mockFn: func(r *mocks.GenreRepository) {
				r.On("Create", ctx, validGenre).Return(nil)
			},
		},
		{
			name:      "forbidden",
			claims:    userClaims,
			genre:     validGenre,
			mockFn:    func(r *mocks.GenreRepository) {},
			wantError: commonerr.ErrForbidden,
		},
		{
			name:   "repo error",
			claims: adminClaims,
			genre:  validGenre,
			mockFn: func(r *mocks.GenreRepository) {
				r.On("Create", ctx, validGenre).Return(errors.New("db error"))
			},
			wantError: usecase.ErrCreateGenre,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewGenreRepository(t)
			tt.mockFn(repo)
			svc := genre.NewGenreService(repo)
			err := svc.CreateGenre(ctx, tt.claims, tt.genre)
			if tt.wantError != nil {
				assert.ErrorIs(t, err, tt.wantError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetGenre(t *testing.T) {
	ctx := context.Background()
	genreID := uuid.New()
	genreEntity := &entity.Genre{ID: genreID}

	t.Run("success", func(t *testing.T) {
		repo := mocks.NewGenreRepository(t)
		repo.On("GetByID", ctx, genreID).Return(genreEntity, nil)
		svc := genre.NewGenreService(repo)

		res, err := svc.GetGenreByID(ctx, genreID)
		assert.NoError(t, err)
		assert.Equal(t, genreEntity, res)
	})

	t.Run("invalid ID", func(t *testing.T) {
		repo := mocks.NewGenreRepository(t)
		svc := genre.NewGenreService(repo)

		_, err := svc.GetGenreByID(ctx, uuid.Nil)
		assert.ErrorIs(t, err, genre.ErrInvalidGenreID)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := mocks.NewGenreRepository(t)
		repo.On("GetByID", ctx, genreID).Return(nil, errors.New("repo error"))
		svc := genre.NewGenreService(repo)

		_, err := svc.GetGenreByID(ctx, genreID)
		assert.ErrorIs(t, err, usecase.ErrGetGenre)
	})
}

func TestUpdateGenre(t *testing.T) {
	ctx := context.Background()
	adminClaims := &entity.Claims{AccessLvl: entity.Admin}
	userClaims := &entity.Claims{AccessLvl: entity.User}
	validGenre := &entity.Genre{ID: uuid.New()}

	tests := []struct {
		name      string
		claims    *entity.Claims
		genre     *entity.Genre
		mockFn    func(r *mocks.GenreRepository)
		wantError error
	}{
		{
			name:   "success",
			claims: adminClaims,
			genre:  validGenre,
			mockFn: func(r *mocks.GenreRepository) {
				r.On("Update", ctx, validGenre).Return(nil)
			},
		},
		{
			name:      "forbidden",
			claims:    userClaims,
			genre:     validGenre,
			mockFn:    func(r *mocks.GenreRepository) {},
			wantError: commonerr.ErrForbidden,
		},
		{
			name:      "invalid ID",
			claims:    adminClaims,
			genre:     &entity.Genre{ID: uuid.Nil},
			mockFn:    func(r *mocks.GenreRepository) {},
			wantError: genre.ErrInvalidGenreID,
		},
		{
			name:   "repo error",
			claims: adminClaims,
			genre:  validGenre,
			mockFn: func(r *mocks.GenreRepository) {
				r.On("Update", ctx, validGenre).Return(errors.New("db error"))
			},
			wantError: usecase.ErrUpdateGenre,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewGenreRepository(t)
			tt.mockFn(repo)
			svc := genre.NewGenreService(repo)
			err := svc.UpdateGenre(ctx, tt.claims, tt.genre)
			if tt.wantError != nil {
				assert.ErrorIs(t, err, tt.wantError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteGenre(t *testing.T) {
	ctx := context.Background()
	genreID := uuid.New()
	adminClaims := &entity.Claims{AccessLvl: entity.Admin}
	userClaims := &entity.Claims{AccessLvl: entity.User}

	tests := []struct {
		name      string
		claims    *entity.Claims
		genreID   uuid.UUID
		mockFn    func(r *mocks.GenreRepository)
		wantError error
	}{
		{
			name:    "success",
			claims:  adminClaims,
			genreID: genreID,
			mockFn: func(r *mocks.GenreRepository) {
				r.On("Delete", ctx, genreID).Return(nil)
			},
		},
		{
			name:      "forbidden",
			claims:    userClaims,
			genreID:   genreID,
			mockFn:    func(r *mocks.GenreRepository) {},
			wantError: commonerr.ErrForbidden,
		},
		{
			name:      "invalid ID",
			claims:    adminClaims,
			genreID:   uuid.Nil,
			mockFn:    func(r *mocks.GenreRepository) {},
			wantError: genre.ErrInvalidGenreID,
		},
		{
			name:    "repo error",
			claims:  adminClaims,
			genreID: genreID,
			mockFn: func(r *mocks.GenreRepository) {
				r.On("Delete", ctx, genreID).Return(errors.New("repo error"))
			},
			wantError: usecase.ErrDeleteGenre,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewGenreRepository(t)
			tt.mockFn(repo)
			svc := genre.NewGenreService(repo)
			err := svc.DeleteGenre(ctx, tt.claims, tt.genreID)
			if tt.wantError != nil {
				assert.ErrorIs(t, err, tt.wantError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
