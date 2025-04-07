package coverservice

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entities"
)

type coverRepository interface {
	SaveCover(ctx context.Context, cover *entities.Cover) error
	GetCover(ctx context.Context, objectID uuid.UUID) (*entities.Cover, error)
	DeleteCover(ctx context.Context, objectID uuid.UUID) error
}

type playlistPolicyService interface {
	CanDelete(ctx context.Context, claims *entities.Claims, playlistID uuid.UUID) (bool, error)
	CanEdit(ctx context.Context, claims *entities.Claims, playlistID uuid.UUID) (bool, error)
	CanView(ctx context.Context, claims *entities.Claims, playlistID uuid.UUID) (bool, error)
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

func (c *CoverService) UploadCover(ctx context.Context, claims *entities.Claims, cover *entities.Cover) error {
	if cover.ObjectType == entities.CoverAlbum && claims.AccessLvl != entities.Admin {
		return errors.New("forbidden")
	}

	if cover.ObjectType == entities.CoverPlaylist {
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

func (c *CoverService) GetCover(ctx context.Context, claims *entities.Claims, objID uuid.UUID, objType entities.CoverObjectType) (*entities.Cover, error) {
	if objType == entities.CoverPlaylist {
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

func (c *CoverService) DeleteCover(ctx context.Context, claims *entities.Claims, objID uuid.UUID, objType entities.CoverObjectType) error {
	if objType == entities.CoverPlaylist {
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
