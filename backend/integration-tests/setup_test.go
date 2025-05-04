package integration_test

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pressly/goose/v3"
	"github.com/redis/go-redis/v9"
)

const (
	migrationsDir = "../db/migrations"
)

var (
	pgxPool     *pgxpool.Pool
	redisClient *redis.Client
	minioClient *minio.Client
)

func TestMain(m *testing.M) {
	var (
		err  error
		code int

		dockerPool     *dockertest.Pool
		dockerPostgres *dockertest.Resource
		dockerMinIO    *dockertest.Resource
		dockerRedis    *dockertest.Resource
	)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}

		teardown(dockerPool, []*dockertest.Resource{dockerPostgres, dockerMinIO, dockerRedis})
		os.Exit(code)
	}()

	dockerPool, err = dockertest.NewPool("")
	if err != nil {
		panic(fmt.Sprintf("failed to start docker: %v", err))
	}

	dockerPostgres, err = setupPostgres(dockerPool)
	if err != nil {
		panic(fmt.Sprintf("failed to start postgres: %v", err))
	}

	dockerMinIO, err = setupMinIO(dockerPool)
	if err != nil {
		panic(fmt.Sprintf("failed to start minio: %v", err))
	}

	dockerRedis, err = setupRedis(dockerPool)
	if err != nil {
		panic(fmt.Sprintf("failed to start redis: %v", err))
	}

	code = m.Run()
}

func runMigrationsUp(dbpool *pgxpool.Pool) error {
	db, err := sql.Open("pgx", dbpool.Config().ConnString())
	if err != nil {
		return err
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(db, migrationsDir)
}

func runMigrationsDown(dbpool *pgxpool.Pool) error {
	db, err := sql.Open("pgx", dbpool.Config().ConnString())
	if err != nil {
		return err
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Reset(db, migrationsDir)
}

func setupPostgres(dockerPool *dockertest.Pool) (*dockertest.Resource, error) {
	resource, err := dockerPool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "17",
		Env: []string{
			"POSTGRES_USER=testuser",
			"POSTGRES_PASSWORD=password",
			"POSTGRES_DB=testdb",
		},
		ExposedPorts: []string{"5432/tcp"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432/tcp": {{HostIP: "0.0.0.0", HostPort: "5432"}},
		},
	})
	if err != nil {
		return nil, err
	}

	port := resource.GetPort("5432/tcp")
	connString := fmt.Sprintf("postgres://testuser:password@localhost:%s/testdb?sslmode=disable", port)

	if err := dockerPool.Retry(func() error {
		var err error
		pgxPool, err = pgxpool.New(context.Background(), connString)
		if err != nil {
			return err
		}
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		return pgxPool.Ping(ctx)
	}); err != nil {
		return nil, err
	}

	return resource, nil
}

func setupRedis(dockerPool *dockertest.Pool) (*dockertest.Resource, error) {
	resource, err := dockerPool.RunWithOptions(&dockertest.RunOptions{
		Repository:   "redis",
		Tag:          "7", // "latest"
		ExposedPorts: []string{"6379/tcp"},
	})
	if err != nil {
		return nil, err
	}

	port := resource.GetPort("6379/tcp")
	addr := fmt.Sprintf("localhost:%s", port)

	if err := dockerPool.Retry(func() error {
		redisClient = redis.NewClient(&redis.Options{
			Addr: addr,
		})
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		return redisClient.Ping(ctx).Err()
	}); err != nil {
		return nil, err
	}

	return resource, nil
}

func setupMinIO(dockerPool *dockertest.Pool) (*dockertest.Resource, error) {
	resource, err := dockerPool.RunWithOptions(&dockertest.RunOptions{
		Repository: "minio/minio",
		Tag:        "latest",
		Env: []string{
			"MINIO_ROOT_USER=minioadmin",
			"MINIO_ROOT_PASSWORD=minioadmin",
		},
		Cmd:          []string{"server", "/data"},
		ExposedPorts: []string{"9000/tcp"},
	})
	if err != nil {
		return nil, err
	}

	port := resource.GetPort("9000/tcp")
	endpoint := fmt.Sprintf("localhost:%s", port)

	if err := dockerPool.Retry(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		minioClient, err = minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
			Secure: false,
		})
		if err != nil {
			return err
		}
		_, err = minioClient.ListBuckets(ctx)
		return err

	}); err != nil {
		return nil, err
	}

	return resource, nil
}

func teardown(pool *dockertest.Pool, resources []*dockertest.Resource) {
	if pgxPool != nil {
		pgxPool.Close()
	}

	if redisClient != nil {
		redisClient.Close()
	}

	for i := range resources {
		if resources[i] == nil {
			continue
		}
		if err := pool.Purge(resources[i]); err != nil {
			slog.Error("failed to purge docker resource", "err", err)
		}
	}
}
