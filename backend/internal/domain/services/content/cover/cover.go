package cover

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/entities"
)

type CoverService interface {
	UploadCover(ctx context.Context, id uuid.UUID, cover *entities.Cover) error
	GetCover(ctx context.Context, id uuid.UUID) (*entities.Cover, error)
	DeleteCover(ctx context.Context, id uuid.UUID) error
}
