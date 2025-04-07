package audioservice

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entities"
)

// TODO: Ошибки обработать правильно

type audioFileRepository interface {
	GetAudioChunk(ctx context.Context, chunk *entities.AudioChunk) (*entities.AudioChunk, error)
	UploadAudioFile(ctx context.Context, chunk *entities.AudioChunk) error
	DeleteFile(ctx context.Context, trackID uuid.UUID) error
}

type AudioFileService struct {
	repo audioFileRepository
}

func New(repo audioFileRepository) *AudioFileService {
	return &AudioFileService{
		repo: repo,
	}
}

func (a *AudioFileService) GetAudioChunk(ctx context.Context, chunk *entities.AudioChunk) (*entities.AudioChunk, error) {
	if chunk.End <= chunk.Start {
		return nil, errors.New("invalid request")
	}

	return a.repo.GetAudioChunk(ctx, chunk)
}

func (a *AudioFileService) UploadAudioFile(ctx context.Context, claims *entities.Claims, chunk *entities.AudioChunk) error {
	switch {
	case claims.AccessLvl != entities.Admin:
		return errors.New("forbidden")
	case chunk.End <= chunk.Start || chunk.Start != 0 || chunk.End != uint64(len(chunk.Data)):
		return errors.New("invalid request")
	}

	return a.repo.UploadAudioFile(ctx, chunk)
}

func (a *AudioFileService) DeleteFile(ctx context.Context, claims *entities.Claims, trackID uuid.UUID) error {
	switch {
	case claims.AccessLvl != entities.Admin:
		return errors.New("forbidden")
	}

	return a.repo.DeleteFile(ctx, trackID)
}
