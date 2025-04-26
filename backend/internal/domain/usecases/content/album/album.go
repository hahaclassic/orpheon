package album

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

var (
	ErrGetAlbum    = errors.New("get album error")
	ErrCreateAlbum = errors.New("get album error")
	ErrUpdateAlbum = errors.New("get album error")
	ErrDeleteAlbum = errors.New("get album error")
)

type AlbumService interface {
	GetAlbum(ctx context.Context, albumID uuid.UUID) (*entity.AlbumMeta, error)

	// Admin
	CreateAlbum(ctx context.Context, claims *entity.Claims, album *entity.AlbumMeta) error
	UpdateAlbum(ctx context.Context, claims *entity.Claims, album *entity.AlbumMeta) error
	DeleteAlbum(ctx context.Context, claims *entity.Claims, albumID uuid.UUID) error
}
