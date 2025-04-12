package meta

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

type AlbumRepository interface {
	CreateAlbum(ctx context.Context, album *entity.AlbumMeta) error
	GetAlbum(ctx context.Context, id uuid.UUID) (*entity.AlbumMeta, error)
	UpdateAlbum(ctx context.Context, album *entity.AlbumMeta) error
	DeleteAlbum(ctx context.Context, id uuid.UUID) error
}

type AlbumService struct {
	repo AlbumRepository
}

func New(repo AlbumRepository) *AlbumService {
	return &AlbumService{
		repo: repo,
	}
}

func (a *AlbumService) CreateAlbum(ctx context.Context, claims *entity.Claims, album *entity.AlbumMeta) error {
	if claims.AccessLvl != entity.Admin {
		return errors.New("forbidden")
	}

	var err error
	album.ID, err = uuid.NewRandom()
	if err != nil {
		return err
	}

	return a.repo.CreateAlbum(ctx, album)
}

func (a *AlbumService) GetAlbum(ctx context.Context, albumID uuid.UUID) (*entity.AlbumMeta, error) {
	return a.repo.GetAlbum(ctx, albumID)
}

func (a *AlbumService) UpdateAlbum(ctx context.Context, claims *entity.Claims, album *entity.AlbumMeta) error {
	if claims.AccessLvl != entity.Admin {
		return errors.New("forbidden")
	}

	return a.repo.UpdateAlbum(ctx, album)
}

func (a *AlbumService) DeleteAlbum(ctx context.Context, claims *entity.Claims, albumID uuid.UUID) error {
	if claims.AccessLvl != entity.Admin {
		return errors.New("forbidden")
	}

	return a.repo.DeleteAlbum(ctx, albumID)
}
