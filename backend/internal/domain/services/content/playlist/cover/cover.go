package service

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

	ErrForbidden = errors.New("permission denied")
)

type PlaylistCoverRepository interface {
	SaveCover(ctx context.Context, cover *entity.Cover) error
	GetCover(ctx context.Context, playlistID uuid.UUID) (*entity.Cover, error)
	DeleteCover(ctx context.Context, playlistID uuid.UUID) error
}

type PlaylistPolicyService interface {
	CanDelete(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (bool, error)
	CanEdit(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (bool, error)
	CanView(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (bool, error)
}

type PlaylistCoverService struct {
	policy PlaylistPolicyService
	repo   PlaylistCoverRepository
}

func New(repo PlaylistCoverRepository) *PlaylistCoverService {
	return &PlaylistCoverService{
		repo: repo,
	}
}

func (c *PlaylistCoverService) GetCover(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (*entity.Cover, error) {
	canUpdate, err := c.policy.CanEdit(ctx, claims, playlistID)
	if err != nil {
		return nil, errwrap.Wrap(ErrGetCover, err)
	}
	if !canUpdate {
		return nil, errwrap.Wrap(ErrGetCover, ErrForbidden)
	}

	cover, err := c.repo.GetCover(ctx, playlistID)
	if err != nil {
		return nil, errwrap.Wrap(ErrGetCover, err)
	}

	return cover, nil
}

func (c *PlaylistCoverService) UploadCover(ctx context.Context, claims *entity.Claims, cover *entity.Cover) error {
	canUpdate, err := c.policy.CanEdit(ctx, claims, cover.ObjectID)
	if err != nil {
		return errwrap.Wrap(ErrUploadCover, err)
	}
	if !canUpdate {
		return errwrap.Wrap(ErrUploadCover, ErrForbidden)
	}

	err = c.repo.SaveCover(ctx, cover)

	return errwrap.WrapIfErr(ErrUploadCover, err)
}

func (c *PlaylistCoverService) DeleteCover(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) error {
	canUpdate, err := c.policy.CanDelete(ctx, claims, playlistID)
	if err != nil {
		return errwrap.Wrap(ErrDeleteCover, err)
	}
	if !canUpdate {
		return errwrap.Wrap(ErrDeleteCover, ErrForbidden)
	}

	err = c.repo.DeleteCover(ctx, playlistID)

	return errwrap.WrapIfErr(ErrDeleteCover, err)
}
