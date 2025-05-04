package audio_minio

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/minio/minio-go/v7"
)

type AudioFileRepository struct {
	minioClient *minio.Client
	bucketName  string
}

func NewAudioFileRepository(client *minio.Client, bucket string) *AudioFileRepository {
	return &AudioFileRepository{
		minioClient: client,
		bucketName:  bucket,
	}
}

func (r *AudioFileRepository) UploadAudioFile(ctx context.Context, chunk *entity.AudioChunk) error {
	objectName := chunk.TrackID.String()

	_, err := r.minioClient.PutObject(ctx, r.bucketName, objectName,
		bytesToReader(chunk.Data), int64(len(chunk.Data)),
		minio.PutObjectOptions{ContentType: "audio/mpeg"})
	if err != nil {
		return fmt.Errorf("failed to upload audio file: %w", err)
	}

	return nil
}

func (r *AudioFileRepository) GetAudioChunk(ctx context.Context, chunk *entity.AudioChunk) (*entity.AudioChunk, error) {
	objectName := chunk.TrackID.String()

	opts := minio.GetObjectOptions{}
	if err := opts.SetRange(int64(chunk.Start), int64(chunk.End-1)); err != nil {
		return nil, fmt.Errorf("invalid range: %w", err)
	}

	obj, err := r.minioClient.GetObject(ctx, r.bucketName, objectName, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get audio chunk: %w", err)
	}
	defer func() {
		err := obj.Close()
		if err != nil {
			slog.Error("err", "object close error", err)
		}
	}()

	data, err := io.ReadAll(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to read chunk: %w", err)
	}

	return &entity.AudioChunk{
		Data:    data,
		TrackID: chunk.TrackID,
		Start:   chunk.Start,
		End:     chunk.End,
	}, nil
}

func (r *AudioFileRepository) DeleteFile(ctx context.Context, trackID uuid.UUID) error {
	objectName := trackID.String()
	err := r.minioClient.RemoveObject(ctx, r.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// вспомогательная функция
func bytesToReader(b []byte) io.ReadSeeker {
	return &byteReader{b, 0}
}

type byteReader struct {
	data []byte
	pos  int
}

func (r *byteReader) Read(p []byte) (int, error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n := copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}

func (r *byteReader) Seek(offset int64, whence int) (int64, error) {
	var abs int
	switch whence {
	case io.SeekStart:
		abs = int(offset)
	case io.SeekCurrent:
		abs = r.pos + int(offset)
	case io.SeekEnd:
		abs = len(r.data) + int(offset)
	default:
		return 0, fmt.Errorf("invalid whence")
	}
	if abs < 0 || abs > len(r.data) {
		return 0, fmt.Errorf("invalid seek position")
	}
	r.pos = abs
	return int64(abs), nil
}
