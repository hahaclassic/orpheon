package playlist_cover_minio

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type PlaylistCoverRepository struct {
	minio  *minio.Client
	bucket string
}

func NewPlaylistCoverRepository(minioEndpoint, minioAccessKey, minioSecretKey, bucketName string) (*PlaylistCoverRepository, error) {
	client, err := minio.New(minioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioAccessKey, minioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MinIO: %w", err)
	}

	return &PlaylistCoverRepository{
		minio:  client,
		bucket: bucketName,
	}, nil
}

func (r *PlaylistCoverRepository) SaveCover(ctx context.Context, cover *entity.Cover) error {
	objectName := fmt.Sprintf("playlist_covers/%s", cover.ObjectID)

	reader := bytes.NewReader(cover.Data)

	_, err := r.minio.PutObject(ctx, r.bucket, objectName, reader, int64(len(cover.Data)), minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})
	if err != nil {
		return fmt.Errorf("failed to upload cover to MinIO: %w", err)
	}

	return nil
}

func (r *PlaylistCoverRepository) GetCover(ctx context.Context, objectID uuid.UUID) (*entity.Cover, error) {
	objectName := fmt.Sprintf("playlist_covers/%s", objectID)

	object, err := r.minio.GetObject(ctx, r.bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get cover from MinIO: %w", err)
	}
	defer object.Close()

	data, err := io.ReadAll(object)
	if err != nil {
		return nil, fmt.Errorf("failed to read cover data: %w", err)
	}

	return &entity.Cover{
		ObjectID: objectID,
		Data:     data,
	}, nil
}

func (r *PlaylistCoverRepository) DeleteCover(ctx context.Context, objectID uuid.UUID) error {
	objectName := fmt.Sprintf("playlist_covers/%s", objectID)

	err := r.minio.RemoveObject(ctx, r.bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete cover from MinIO: %w", err)
	}

	return nil
}
