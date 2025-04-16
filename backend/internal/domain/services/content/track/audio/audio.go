package audio

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/track"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

var (
	ErrForbidden          = errors.New("permission denied")
	ErrInvalidChunkParams = errors.New("invalid chunk parameters")
)

type audioFileRepository interface {
	GetAudioChunk(ctx context.Context, chunk *entity.AudioChunk) (*entity.AudioChunk, error)
	UploadAudioFile(ctx context.Context, chunk *entity.AudioChunk) error
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

func (a *AudioFileService) GetAudioChunk(ctx context.Context, chunk *entity.AudioChunk) (result *entity.AudioChunk, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetAudioChunk, err)
	}()

	if chunk.End <= chunk.Start {
		return nil, ErrInvalidChunkParams
	}

	return a.repo.GetAudioChunk(ctx, chunk)
}

func (a *AudioFileService) UploadAudioFile(ctx context.Context, claims *entity.Claims, chunk *entity.AudioChunk) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrUploadAudioFile, err)
	}()

	switch {
	case claims.AccessLvl != entity.Admin:
		return ErrForbidden
	case chunk.End <= chunk.Start || chunk.Start != 0 || chunk.End != uint64(len(chunk.Data)):
		return ErrInvalidChunkParams
	}

	return a.repo.UploadAudioFile(ctx, chunk)
}

func (a *AudioFileService) DeleteAudioFile(ctx context.Context, claims *entity.Claims, trackID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrDeleteAudioFile, err)
	}()

	if claims.AccessLvl != entity.Admin {
		return ErrForbidden
	}

	return a.repo.DeleteFile(ctx, trackID)
}
