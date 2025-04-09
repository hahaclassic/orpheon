package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

type coverRepository interface {
	SaveCover(ctx context.Context, cover *entity.Cover) error
	GetCover(ctx context.Context, objectID uuid.UUID) (*entity.Cover, error)
	DeleteCover(ctx context.Context, objectID uuid.UUID) error
}

type playlistPolicyService interface {
	CanDelete(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (bool, error)
	CanEdit(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (bool, error)
	CanView(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (bool, error)
}

type CoverService struct {
	policy playlistPolicyService
	repo   coverRepository
}

func New(repo coverRepository) *CoverService {
	return &CoverService{
		repo: repo,
	}
}

func (c *CoverService) UploadCover(ctx context.Context, claims *entity.Claims, cover *entity.Cover) error {
	if cover.ObjectType == entity.CoverAlbum && claims.AccessLvl != entity.Admin {
		return errors.New("forbidden")
	}

	if cover.ObjectType == entity.CoverPlaylist {
		canUpdate, err := c.policy.CanEdit(ctx, claims, cover.ObjectID)
		if err != nil {
			return err
		}
		if !canUpdate {
			return errors.New("forbidden")
		}
	}

	return c.repo.SaveCover(ctx, cover)
}

func (c *CoverService) GetCover(ctx context.Context, claims *entity.Claims, objID uuid.UUID, objType entity.CoverObjectType) (*entity.Cover, error) {
	if objType == entity.CoverPlaylist {
		canUpdate, err := c.policy.CanEdit(ctx, claims, objID)
		if err != nil {
			return nil, err
		}
		if !canUpdate {
			return nil, errors.New("forbidden")
		}
	}

	return c.repo.GetCover(ctx, objID)
}

func (c *CoverService) DeleteCover(ctx context.Context, claims *entity.Claims, objID uuid.UUID, objType entity.CoverObjectType) error {
	if objType == entity.CoverPlaylist {
		canUpdate, err := c.policy.CanDelete(ctx, claims, objID)
		if err != nil {
			return err
		}
		if !canUpdate {
			return errors.New("forbidden")
		}
	}

	return c.repo.DeleteCover(ctx, objID)
}
