package cover

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

type CoverService interface {
	UploadCover(ctx context.Context, claims *entity.Claims, cover *entity.Cover) error
	GetCover(ctx context.Context, claims *entity.Claims, objID uuid.UUID) (*entity.Cover, error)
	DeleteCover(ctx context.Context, claims *entity.Claims, objID uuid.UUID) error
}
