package album

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entities"
)

type AlbumService interface {
	CreateAlbum(ctx context.Context, album *entities.AlbumMeta) error
	GetAlbum(ctx context.Context, albumID uuid.UUID) (*entities.AlbumMeta, error)
	UpdateAlbum(ctx context.Context, album *entities.AlbumMeta) error
	DeleteAlbum(ctx context.Context, albumID uuid.UUID) error
}
