package albumservice

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entities"
)

type AlbumRepository interface {
	CreateAlbum(ctx context.Context, album *entities.AlbumMeta) error
	GetAlbum(ctx context.Context, id uuid.UUID) (*entities.AlbumMeta, error)
	UpdateAlbum(ctx context.Context, album *entities.AlbumMeta) error
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

func (a *AlbumService) CreateAlbum(ctx context.Context, claims *entities.Claims, album *entities.AlbumMeta) error {
	if claims.AccessLvl != entities.Admin {
		return errors.New("forbidden")
	}

	var err error
	album.ID, err = uuid.NewRandom()
	if err != nil {
		return err
	}

	return a.repo.CreateAlbum(ctx, album)
}

func (a *AlbumService) GetAlbum(ctx context.Context, albumID uuid.UUID) (*entities.AlbumMeta, error) {
	return a.repo.GetAlbum(ctx, albumID)
}

func (a *AlbumService) UpdateAlbum(ctx context.Context, claims *entities.Claims, album *entities.AlbumMeta) error {
	if claims.AccessLvl != entities.Admin {
		return errors.New("forbidden")
	}

	return a.repo.UpdateAlbum(ctx, album)
}

func (a *AlbumService) DeleteAlbum(ctx context.Context, claims *entities.Claims, albumID uuid.UUID) error {
	if claims.AccessLvl != entities.Admin {
		return errors.New("forbidden")
	}

	return a.repo.DeleteAlbum(ctx, albumID)
}
