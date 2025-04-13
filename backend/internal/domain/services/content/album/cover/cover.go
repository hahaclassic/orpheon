package cover

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/cover"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

var (
	ErrForbidden = errors.New("permission denied")
)

type AlbumCoverRepository interface {
	GetCover(ctx context.Context, albumID uuid.UUID) (*entity.Cover, error)
	SaveCover(ctx context.Context, cover *entity.Cover) error
	DeleteCover(ctx context.Context, albumID uuid.UUID) error
}

type AlbumCoverService struct {
	repo AlbumCoverRepository
}

func New(repo AlbumCoverRepository) *AlbumCoverService {
	return &AlbumCoverService{
		repo: repo,
	}
}

func (c *AlbumCoverService) GetCover(ctx context.Context, claims *entity.Claims, albumID uuid.UUID) (_ *entity.Cover, err error) {
	defer func() {
		if err != nil {
			err = errwrap.Wrap(usecase.ErrGetCover, err)
		}
	}()

	cover, err := c.repo.GetCover(ctx, albumID)
	if err != nil {
		return nil, err
	}

	return cover, nil
}

func (c *AlbumCoverService) UploadCover(ctx context.Context, claims *entity.Claims, cover *entity.Cover) (err error) {
	defer func() {
		if err != nil {
			err = errwrap.Wrap(usecase.ErrUploadCover, err)
		}
	}()

	if claims.AccessLvl != entity.Admin {
		return ErrForbidden
	}

	return c.repo.SaveCover(ctx, cover)
}

func (c *AlbumCoverService) DeleteCover(ctx context.Context, claims *entity.Claims, albumID uuid.UUID) (err error) {
	defer func() {
		if err != nil {
			err = errwrap.Wrap(usecase.ErrDeleteCover, err)
		}
	}()

	if claims.AccessLvl != entity.Admin {
		return ErrForbidden
	}

	return c.repo.DeleteCover(ctx, albumID)
}
