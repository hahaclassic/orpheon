package cover

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/models"
)

type CoverService interface {
	UploadCover(ctx context.Context, id uuid.UUID, cover *models.Cover) error
	GetCover(ctx context.Context, id uuid.UUID) (*models.Cover, error)
	DeleteCover(ctx context.Context, id uuid.UUID) error
}
