package album

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entities"
)

type AlbumService interface {
	GetAlbum(ctx context.Context, albumID uuid.UUID) (*entities.AlbumMeta, error)

	// Admin
	CreateAlbum(ctx context.Context, claims *entities.Claims, album *entities.AlbumMeta) error
	UpdateAlbum(ctx context.Context, claims *entities.Claims, album *entities.AlbumMeta) error
	DeleteAlbum(ctx context.Context, claims *entities.Claims, albumID uuid.UUID) error
}
