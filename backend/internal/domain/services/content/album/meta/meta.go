package meta

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/album"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

var (
	ErrGenerateID = errors.New("id generation error")
)

type AlbumRepository interface {
	CreateAlbum(ctx context.Context, album *entity.AlbumMeta) error
	GetAlbum(ctx context.Context, id uuid.UUID) (*entity.AlbumMeta, error)
	UpdateAlbum(ctx context.Context, album *entity.AlbumMeta) error
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

func (a *AlbumService) CreateAlbum(ctx context.Context, claims *entity.Claims, album *entity.AlbumMeta) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrCreateAlbum, err)
	}()
	if claims.AccessLvl != entity.Admin {
		return commonerr.ErrForbidden
	}

	album.ID, err = uuid.NewRandom()
	if err != nil {
		return ErrGenerateID
	}

	return a.repo.CreateAlbum(ctx, album)
}

func (a *AlbumService) GetAlbum(ctx context.Context, albumID uuid.UUID) (_ *entity.AlbumMeta, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetAlbum, err)
	}()

	return a.repo.GetAlbum(ctx, albumID)
}

func (a *AlbumService) UpdateAlbum(ctx context.Context, claims *entity.Claims, album *entity.AlbumMeta) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrUpdateAlbum, err)
	}()

	if claims.AccessLvl != entity.Admin {
		return commonerr.ErrForbidden
	}

	return a.repo.UpdateAlbum(ctx, album)
}

func (a *AlbumService) DeleteAlbum(ctx context.Context, claims *entity.Claims, albumID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrDeleteAlbum, err)
	}()

	if claims.AccessLvl != entity.Admin {
		return commonerr.ErrForbidden
	}

	return a.repo.DeleteAlbum(ctx, albumID)
}
