package cover

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

var (
	ErrUploadCover = errors.New("upload cover error")
	ErrGetCover    = errors.New("get cover error")
	ErrDeleteCover = errors.New("delete cover error")

	ErrForbidden = errors.New("access forbidden")
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

func (c *AlbumCoverService) GetCover(ctx context.Context, claims *entity.Claims, albumID uuid.UUID) (*entity.Cover, error) {
	cover, err := c.repo.GetCover(ctx, albumID)
	if err != nil {
		return nil, errwrap.Wrap(ErrGetCover, err)
	}

	return cover, nil
}

func (c *AlbumCoverService) UploadCover(ctx context.Context, claims *entity.Claims, cover *entity.Cover) error {
	if claims.AccessLvl != entity.Admin {
		return errwrap.Wrap(ErrUploadCover, ErrForbidden)
	}

	err := c.repo.SaveCover(ctx, cover)
	if err != nil {
		return errwrap.Wrap(ErrUploadCover, err)
	}

	return nil
}

func (c *AlbumCoverService) DeleteCover(ctx context.Context, claims *entity.Claims, albumID uuid.UUID) error {
	if claims.AccessLvl != entity.Admin {
		return errwrap.Wrap(ErrDeleteCover, ErrForbidden)
	}

	err := c.repo.DeleteCover(ctx, albumID)
	if err != nil {
		return errwrap.Wrap(ErrDeleteCover, err)
	}

	return nil
}
