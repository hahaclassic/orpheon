package album

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/pkg/commonerr"
	"github.com/hahaclassic/orpheon/pkg/errwrap"
	"github.com/hahaclassic/orpheon/services/content-msv/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/services/content-msv/internal/domain/usecase/album"
)

var (
	ErrGenerateID = errors.New("id generation error")
)

type AlbumRepository interface {
	CreateAlbum(ctx context.Context, album *entity.AlbumMeta) error
	GetAlbum(ctx context.Context, id uuid.UUID) (*entity.AlbumMeta, error)
	UpdateAlbum(ctx context.Context, album *entity.AlbumMeta) error
	DeleteAlbum(ctx context.Context, id uuid.UUID) error
	GetAllAlbums(ctx context.Context) ([]*entity.AlbumMeta, error)
	GetAlbumArtists(ctx context.Context, albumID uuid.UUID) ([]*entity.ArtistMeta, error)
	GetAlbumGenres(ctx context.Context, albumID uuid.UUID) ([]*entity.Genre, error)
}

type AlbumMetaService struct {
	repo AlbumRepository
}

func NewAlbumMetaService(repo AlbumRepository) *AlbumMetaService {
	return &AlbumMetaService{
		repo: repo,
	}
}

func (a *AlbumMetaService) CreateAlbum(ctx context.Context, claims *entity.Claims, album *entity.AlbumMeta) (id uuid.UUID, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrCreateAlbum, err)
	}()
	if claims == nil || claims.AccessLvl != entity.AdminLvl {
		return uuid.Nil, commonerr.ErrForbidden
	}

	album.ID, err = uuid.NewRandom()
	if err != nil {
		return uuid.Nil, ErrGenerateID
	}

	err = a.repo.CreateAlbum(ctx, album)
	if err != nil {
		return uuid.Nil, err
	}

	return album.ID, nil
}

func (a *AlbumMetaService) GetAlbum(ctx context.Context, albumID uuid.UUID) (_ *entity.AlbumMeta, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetAlbum, err)
	}()

	return a.repo.GetAlbum(ctx, albumID)
}

func (a *AlbumMetaService) GetAllAlbums(ctx context.Context) (_ []*entity.AlbumMeta, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetAllAlbums, err)
	}()

	return a.repo.GetAllAlbums(ctx)
}

func (a *AlbumMetaService) UpdateAlbum(ctx context.Context, claims *entity.Claims, album *entity.AlbumMeta) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrUpdateAlbum, err)
	}()

	if claims == nil || claims.AccessLvl != entity.AdminLvl {
		return commonerr.ErrForbidden
	}

	return a.repo.UpdateAlbum(ctx, album)
}

func (a *AlbumMetaService) DeleteAlbum(ctx context.Context, claims *entity.Claims, albumID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrDeleteAlbum, err)
	}()

	if claims == nil || claims.AccessLvl != entity.AdminLvl {
		return commonerr.ErrForbidden
	}

	return a.repo.DeleteAlbum(ctx, albumID)
}
