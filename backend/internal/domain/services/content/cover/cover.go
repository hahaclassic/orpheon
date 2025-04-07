package cover

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entities"
)

type CoverService interface {
	UploadCover(ctx context.Context, claims *entities.Claims, cover *entities.Cover) error
	GetCover(ctx context.Context, claims *entities.Claims, objID uuid.UUID) (*entities.Cover, error)
	DeleteCover(ctx context.Context, claims *entities.Claims, objID uuid.UUID) error
}
