package album

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

type AlbumService interface {
	GetAlbum(ctx context.Context, albumID uuid.UUID) (*entity.AlbumMeta, error)

	// Admin
	CreateAlbum(ctx context.Context, claims *entity.Claims, album *entity.AlbumMeta) error
	UpdateAlbum(ctx context.Context, claims *entity.Claims, album *entity.AlbumMeta) error
	DeleteAlbum(ctx context.Context, claims *entity.Claims, albumID uuid.UUID) error
}
