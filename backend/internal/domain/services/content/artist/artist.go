package artist

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/artist"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

var (
	ErrForbidden  = errors.New("access forbidden error")
	ErrGenerateID = errors.New("id generation error")
)

type ArtistMetaRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ArtistMeta, error)
	Create(ctx context.Context, artist *entity.ArtistMeta) error
	Update(ctx context.Context, artist *entity.ArtistMeta) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ArtistMetaService struct {
	repo ArtistMetaRepository
}

func New(repo ArtistMetaRepository) *ArtistMetaService {
	return &ArtistMetaService{repo: repo}
}

func (s *ArtistMetaService) GetArtistMeta(ctx context.Context, artistID uuid.UUID) (*entity.ArtistMeta, error) {
	return s.repo.GetByID(ctx, artistID)
}

func (s *ArtistMetaService) CreateArtistMeta(ctx context.Context, claims *entity.Claims, artist *entity.ArtistMeta) (err error) {
	defer func() {
		if err != nil {
			err = errwrap.Wrap(usecase.ErrCreateArtistMeta, err)
		}
	}()

	if claims.AccessLvl == entity.Admin {
		return ErrForbidden
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return ErrGenerateID
	}
	artist.ID = id

	return s.repo.Create(ctx, artist)
}

func (s *ArtistMetaService) UpdateArtistMeta(ctx context.Context, claims *entity.Claims, artist *entity.ArtistMeta) (err error) {
	defer func() {
		if err != nil {
			err = errwrap.Wrap(usecase.ErrUpdateArtistMeta, err)
		}
	}()

	if claims.AccessLvl == entity.Admin {
		return ErrForbidden
	}

	return s.repo.Update(ctx, artist)
}

func (s *ArtistMetaService) DeleteArtistMeta(ctx context.Context, claims *entity.Claims, artistID uuid.UUID) (err error) {
	defer func() {
		if err != nil {
			err = errwrap.Wrap(usecase.ErrDeleteArtistMeta, err)
		}
	}()

	if claims.AccessLvl != entity.Admin {
		return ErrForbidden
	}

	return s.repo.Delete(ctx, artistID)
}
