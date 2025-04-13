package avatar

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

var ()

type ArtistCoverRepository interface {
	SaveCover(ctx context.Context, cover *entity.Cover) error
	GetCover(ctx context.Context, artistID uuid.UUID) (*entity.Cover, error)
	DeleteCover(ctx context.Context, artistID uuid.UUID) error
}

type ArtistCoverService struct {
	repo ArtistCoverRepository
}

func NewArtistCoverService(repo ArtistCoverRepository) *ArtistCoverService {
	return &ArtistCoverService{
		repo: repo,
	}
}

func (s *ArtistCoverService) UploadCover(ctx context.Context, claims *entity.Claims, cover *entity.Cover) error {
	if claims.AccessLvl != entity.Admin {
		return errors.New("permission denied")
	}

	return s.repo.SaveCover(ctx, cover)
}

func (s *ArtistCoverService) GetCover(ctx context.Context, claims *entity.Claims, artistID uuid.UUID) (*entity.Cover, error) {
	return s.repo.GetCover(ctx, artistID)
}

func (s *ArtistCoverService) DeleteCover(ctx context.Context, claims *entity.Claims, artistID uuid.UUID) error {
	if claims.AccessLvl != entity.Admin {
		return errors.New("permission denied")
	}

	return s.repo.DeleteCover(ctx, artistID)
}
