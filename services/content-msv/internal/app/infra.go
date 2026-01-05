package app

import (
	"context"
	"log/slog"

	"github.com/hahaclassic/orpheon/pkg/infrastructure/postgres"
	redis_client "github.com/hahaclassic/orpheon/pkg/infrastructure/redis"
	"github.com/hahaclassic/orpheon/services/content-msv/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Infra struct {
	postgres *pgxpool.Pool
	redis    *redis.Client
	//minIO    *minio.Client
}

func (i *Infra) Init(ctx context.Context, cfg *config.Config) error {
	i.postgres = postgres.NewPostgresPool(ctx, postgres.PostgresConfig(cfg.Postgres))

	redisClient, err := redis_client.NewRedisClient(redis_client.RedisConfig(cfg.Redis))
	if err != nil {
		closeErr := i.Close(ctx)
		slog.Error("[CONTENT-MSV] failed to close infra connection", "err", closeErr)
		return err
	}
	i.redis = redisClient

	// minioClient, err := minio_client.NewMinioClient(minio_client.MinIOConfig(cfg.MinIO))
	// if err != nil {
	// 	closeErr := i.Close(ctx)
	// 	slog.Error("[CONTENT-MSV] failed to close infra connection", "err", closeErr)
	// 	return err
	// }
	// i.minIO = minioClient

	return nil
}

func (i *Infra) Close(ctx context.Context) error {
	var err error

	done := make(chan struct{})

	go func() {
		defer close(done)

		if i.postgres != nil {
			i.postgres.Close()
		}
		if i.redis != nil {
			err = i.redis.Close()
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return err
	}
}
