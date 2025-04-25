package assign

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/artist"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

var (
	ErrForbidden = errors.New("permission denied error")
)

type ArtistAssignRepository interface {
	AssignArtistToTrack(ctx context.Context, artistID uuid.UUID, trackID uuid.UUID) error
	AssignArtistToAlbum(ctx context.Context, artistID uuid.UUID, albumID uuid.UUID) error
}

type ArtistAssignService struct {
	repo ArtistAssignRepository
}

func NewArtistAssignService(repo ArtistAssignRepository) *ArtistAssignService {
	return &ArtistAssignService{
		repo: repo,
	}
}

func (a *ArtistAssignService) AssignArtistToTrack(ctx context.Context, claims *entity.Claims, artistID uuid.UUID, trackID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrAssignArtistOnTrack, err)
	}()

	if claims.AccessLvl != entity.Admin {
		return ErrForbidden
	}

	return a.repo.AssignArtistToTrack(ctx, artistID, trackID)
}

func (a *ArtistAssignService) AssignArtistToAlbum(ctx context.Context, claims *entity.Claims, artistID uuid.UUID, albumID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrAssignArtistOnAlbum, err)
	}()

	if claims.AccessLvl != entity.Admin {
		return ErrForbidden
	}

	return a.repo.AssignArtistToAlbum(ctx, artistID, albumID)
}
