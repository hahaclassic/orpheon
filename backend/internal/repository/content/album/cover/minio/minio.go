package album_cover_minio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type AlbumCoverRepository struct {
	client     *minio.Client
	bucketName string
}

func NewAlbumCoverRepository(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*AlbumCoverRepository, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("create minio client: %w", err)
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("check bucket exists: %w", err)
	}
	if !exists {
		err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("create bucket: %w", err)
		}
	}

	return &AlbumCoverRepository{
		client:     client,
		bucketName: bucket,
	}, nil
}

func (r *AlbumCoverRepository) SaveCover(ctx context.Context, cover *entity.Cover) error {
	objectName := fmt.Sprintf("album-covers/%s", cover.ObjectID.String())
	_, err := r.client.PutObject(ctx, r.bucketName, objectName,
		bytes.NewReader(cover.Data),
		int64(len(cover.Data)),
		minio.PutObjectOptions{ContentType: "application/octet-stream"},
	)
	if err != nil {
		return fmt.Errorf("upload cover: %w", err)
	}
	return nil
}

func (r *AlbumCoverRepository) GetCover(ctx context.Context, albumID uuid.UUID) (*entity.Cover, error) {
	objectName := fmt.Sprintf("album-covers/%s", albumID.String())

	obj, err := r.client.GetObject(ctx, r.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("get cover object: %w", err)
	}
	defer func() {
		err := obj.Close()
		if err != nil {
			slog.Error("err", "object close error", err)
		}
	}()

	data, err := io.ReadAll(obj)
	if err != nil {
		return nil, fmt.Errorf("read cover data: %w", err)
	}

	return &entity.Cover{
		ObjectID: albumID,
		Data:     data,
	}, nil
}

func (r *AlbumCoverRepository) DeleteCover(ctx context.Context, albumID uuid.UUID) error {
	objectName := fmt.Sprintf("album-covers/%s", albumID.String())
	err := r.client.RemoveObject(ctx, r.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("delete cover: %w", err)
	}
	return nil
}
