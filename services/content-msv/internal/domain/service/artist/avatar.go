package artist

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/pkg/commonerr"
	"github.com/hahaclassic/orpheon/pkg/errwrap"
	"github.com/hahaclassic/orpheon/services/content-msv/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/services/content-msv/internal/domain/usecase/artist"
)

type ArtistAvatarRepository interface {
	SaveCover(ctx context.Context, cover *entity.Cover) error
	GetCover(ctx context.Context, artistID uuid.UUID) (*entity.Cover, error)
	DeleteCover(ctx context.Context, artistID uuid.UUID) error
}

type ArtistCoverService struct {
	repo ArtistAvatarRepository
}

func NewArtistCoverService(repo ArtistAvatarRepository) *ArtistCoverService {
	return &ArtistCoverService{
		repo: repo,
	}
}

func (s *ArtistCoverService) GetCover(ctx context.Context, artistID uuid.UUID) (_ *entity.Cover, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetAvatar, err)
	}()

	return s.repo.GetCover(ctx, artistID)
}

func (s *ArtistCoverService) UploadCover(ctx context.Context, claims *entity.Claims, cover *entity.Cover) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrUploadAvatar, err)
	}()

	if claims == nil || claims.AccessLvl != entity.AdminLvl {
		return commonerr.ErrForbidden
	}

	return s.repo.SaveCover(ctx, cover)
}

func (s *ArtistCoverService) DeleteCover(ctx context.Context, claims *entity.Claims, artistID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrDeleteAvatar, err)
	}()

	if claims == nil || claims.AccessLvl != entity.AdminLvl {
		return commonerr.ErrForbidden
	}

	return s.repo.DeleteCover(ctx, artistID)
}
