package audio

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entities"
)

type AudioFileService interface {
	GetAudioChunk(ctx context.Context, chunk *entities.AudioChunk) (*entities.AudioChunk, error)
	// Admin
	UploadAudioFile(ctx context.Context, claims *entities.Claims, chunk *entities.AudioChunk) error
	DeleteFile(ctx context.Context, claims *entities.Claims, trackID uuid.UUID) error
}
