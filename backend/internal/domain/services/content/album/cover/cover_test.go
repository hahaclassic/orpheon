package cover_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/content/album/cover"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/album"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/hahaclassic/orpheon/backend/mocks"
)

func TestGetCover(t *testing.T) {
	ctx := context.Background()
	albumID := uuid.New()
	mockCover := &entity.Cover{ObjectID: albumID}

	tests := []struct {
		name    string
		repoRes *entity.Cover
		repoErr error
		wantErr error
	}{
		{"success", mockCover, nil, nil},
		{"repo error", nil, errors.New("db error"), usecase.ErrGetCover},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewAlbumCoverRepository(t)
			repo.On("GetCover", ctx, albumID).Return(tt.repoRes, tt.repoErr)
			svc := cover.New(repo)

			res, err := svc.GetCover(ctx, albumID)
			assert.Equal(t, tt.repoRes, res)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestUploadCover(t *testing.T) {
	ctx := context.Background()
	admin := &entity.Claims{AccessLvl: entity.Admin}
	nonAdmin := &entity.Claims{AccessLvl: entity.User}
	albumID := uuid.New()
	coverData := &entity.Cover{ObjectID: albumID}

	tests := []struct {
		name    string
		claims  *entity.Claims
		repoErr error
		wantErr error
	}{
		{"admin success", admin, nil, nil},
		{"admin repo error", admin, errors.New("db error"), usecase.ErrUploadCover},
		{"non-admin forbidden", nonAdmin, nil, commonerr.ErrForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewAlbumCoverRepository(t)
			if tt.claims.AccessLvl == entity.Admin {
				repo.On("SaveCover", ctx, coverData).Return(tt.repoErr)
			}
			svc := cover.New(repo)
			err := svc.UploadCover(ctx, tt.claims, coverData)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestDeleteCover(t *testing.T) {
	ctx := context.Background()
	admin := &entity.Claims{AccessLvl: entity.Admin}
	nonAdmin := &entity.Claims{AccessLvl: entity.User}
	albumID := uuid.New()

	tests := []struct {
		name    string
		claims  *entity.Claims
		repoErr error
		wantErr error
	}{
		{"admin success", admin, nil, nil},
		{"admin repo error", admin, errors.New("db error"), usecase.ErrDeleteCover},
		{"non-admin forbidden", nonAdmin, nil, commonerr.ErrForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewAlbumCoverRepository(t)
			if tt.claims.AccessLvl == entity.Admin {
				repo.On("DeleteCover", ctx, albumID).Return(tt.repoErr)
			}
			svc := cover.New(repo)
			err := svc.DeleteCover(ctx, tt.claims, albumID)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
