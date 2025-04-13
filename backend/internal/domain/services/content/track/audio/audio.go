package audio

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

var (
	ErrGetAudioChunk   = errors.New("get audio chunk error")
	ErrUploadAudioFile = errors.New("upload audio file error")
	ErrDeleteAudioFile = errors.New("delete audio file error")

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

func (a *AudioFileService) GetAudioChunk(ctx context.Context, chunk *entity.AudioChunk) (*entity.AudioChunk, error) {
	if chunk.End <= chunk.Start {
		return nil, errwrap.Wrap(ErrGetAudioChunk, ErrInvalidChunkParams)
	}

	chunk, err := a.repo.GetAudioChunk(ctx, chunk)
	if err != nil {
		return nil, errwrap.Wrap(ErrGetAudioChunk, err)
	}

	return chunk, nil
}

func (a *AudioFileService) UploadAudioFile(ctx context.Context, claims *entity.Claims, chunk *entity.AudioChunk) error {
	switch {
	case claims.AccessLvl != entity.Admin:
		return errwrap.Wrap(ErrUploadAudioFile, ErrForbidden)
	case chunk.End <= chunk.Start || chunk.Start != 0 || chunk.End != uint64(len(chunk.Data)):
		return errwrap.Wrap(ErrUploadAudioFile, ErrInvalidChunkParams)
	}

	err := a.repo.UploadAudioFile(ctx, chunk)

	return errwrap.WrapIfErr(ErrUploadAudioFile, err)
}

func (a *AudioFileService) DeleteAudioFile(ctx context.Context, claims *entity.Claims, trackID uuid.UUID) error {
	if claims.AccessLvl != entity.Admin {
		return errwrap.Wrap(ErrDeleteAudioFile, ErrForbidden)
	}

	err := a.repo.DeleteFile(ctx, trackID)

	return errwrap.WrapIfErr(ErrDeleteAudioFile, err)
}
