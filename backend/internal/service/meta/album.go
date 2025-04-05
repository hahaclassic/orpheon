package content

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/models"
)

type AlbumService interface {
	CreateAlbum(ctx context.Context, album *models.AlbumMeta) error
	GetAlbum(ctx context.Context, id uuid.UUID) (*models.AlbumMeta, error)
	Update(ctx context.Context, album *models.AlbumMeta) error
	Delete(ctx context.Context, id uuid.UUID) error
}
