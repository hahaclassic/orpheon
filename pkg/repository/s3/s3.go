package s3repo

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/minio/minio-go/v7"
)

var (
	ErrNoBacket = errors.New("s3: the bucket does not exist")
)

type S3Object interface {
	Key() string
	Reader() (io.Reader, int64)
}

type S3Repository[T S3Object] struct {
	client     *minio.Client
	bucketName string
	putOpts    minio.PutObjectOptions
	getOpts    minio.GetObjectOptions
	deleteOpts minio.RemoveObjectOptions
}

func NewS3Repository[T S3Object](client *minio.Client, bucketName string,
	putOpts minio.PutObjectOptions, getOpts minio.GetObjectOptions, deleteOpts minio.RemoveObjectOptions) *S3Repository[T] {

	return &S3Repository[T]{
		client:     client,
		bucketName: bucketName,
		putOpts:    putOpts,
		getOpts:    getOpts,
		deleteOpts: deleteOpts,
	}
}

func (r *S3Repository[T]) Put(ctx context.Context, obj T) error {
	reader, length := obj.Reader()

	_, err := r.client.PutObject(ctx, r.bucketName, obj.Key(),
		reader, length, r.putOpts)
	if err != nil {
		return fmt.Errorf("failed to upload audio file: %w", err)
	}

	return nil
}

func (r *S3Repository[T]) Get(ctx context.Context, key string) ([]byte, error) {
	obj, err := r.client.GetObject(ctx, r.bucketName, key, r.getOpts)
	if err != nil {
		return nil, fmt.Errorf("get object: %w", err)
	}
	defer func() {
		err := obj.Close()
		if err != nil {
			slog.Error("err", "object close error", err)
		}
	}()
	if _, err := obj.Stat(); err != nil {
		return nil, fmt.Errorf("object not found: %w", err)
	}

	data, err := io.ReadAll(obj)
	if err != nil {
		return nil, fmt.Errorf("read data: %w", err)
	}
	return data, nil
}

func (r *S3Repository[T]) GetChunk(ctx context.Context, key string, start int64, end int64) ([]byte, error) {
	opts := r.getOpts
	if err := opts.SetRange(start, end-1); err != nil {
		return nil, fmt.Errorf("invalid range: %w", err)
	}

	obj, err := r.client.GetObject(ctx, r.bucketName, key, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get object with key = %s: %w", key, err)
	}
	defer func() {
		err := obj.Close()
		if err != nil {
			slog.Error("err", "object close error", err)
		}
	}()
	if _, err := obj.Stat(); err != nil {
		return nil, fmt.Errorf("object not found: %w", err)
	}

	data, err := io.ReadAll(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to read chunk: %w", err)
	}

	return data, nil
}

func (r *S3Repository[T]) Delete(ctx context.Context, key string) error {
	err := r.client.RemoveObject(ctx, r.bucketName, key, r.deleteOpts)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}
