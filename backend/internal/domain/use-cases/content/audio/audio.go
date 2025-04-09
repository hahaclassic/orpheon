package audio

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

type AudioFileService interface {
	GetAudioChunk(ctx context.Context, chunk *entity.AudioChunk) (*entity.AudioChunk, error)
	// Admin
	UploadAudioFile(ctx context.Context, claims *entity.Claims, chunk *entity.AudioChunk) error
	DeleteFile(ctx context.Context, claims *entity.Claims, trackID uuid.UUID) error
}
