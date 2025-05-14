package album

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

var (
	ErrGetAlbum     = errors.New("get album error")
	ErrCreateAlbum  = errors.New("create album error")
	ErrUpdateAlbum  = errors.New("update album error")
	ErrDeleteAlbum  = errors.New("delete album error")
	ErrGetAllAlbums = errors.New("get all albums error")
)

type AlbumService interface {
	GetAlbum(ctx context.Context, albumID uuid.UUID) (*entity.AlbumMeta, error)
	GetAllAlbums(ctx context.Context) ([]*entity.AlbumMeta, error)

	// Admin
	CreateAlbum(ctx context.Context, claims *entity.Claims, album *entity.AlbumMeta) error
	UpdateAlbum(ctx context.Context, claims *entity.Claims, album *entity.AlbumMeta) error
	DeleteAlbum(ctx context.Context, claims *entity.Claims, albumID uuid.UUID) error
}
