package access_cache_redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	access_meta "github.com/hahaclassic/orpheon/backend/internal/repository/content/playlist/access-meta/with-cache"
	"github.com/redis/go-redis/v9"
)

type AccessCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewAccessCache(client *redis.Client, ttl time.Duration) *AccessCache {
	return &AccessCache{
		client: client,
		ttl:    ttl,
	}
}

func (a *AccessCache) key(id uuid.UUID) string {
	return "playlist_access:" + id.String()
}

func (a *AccessCache) Set(ctx context.Context, playlistID uuid.UUID, meta *entity.PlaylistAccessMeta) error {
	data, err := json.Marshal(meta)
	if err != nil {
		return fmt.Errorf("failed to marshal playlist access meta: %w", err)
	}

	if err := a.client.Set(ctx, a.key(playlistID), data, a.ttl).Err(); err != nil {
		return fmt.Errorf("failed to set playlist access meta in redis: %w", err)
	}

	return nil
}

func (a *AccessCache) Get(ctx context.Context, playlistID uuid.UUID) (*entity.PlaylistAccessMeta, error) {
	data, err := a.client.Get(ctx, a.key(playlistID)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, access_meta.ErrCacheMiss
		}
		return nil, fmt.Errorf("failed to get playlist access meta from redis: %w", err)
	}

	var meta entity.PlaylistAccessMeta
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, fmt.Errorf("failed to unmarshal playlist access meta: %w", err)
	}

	return &meta, nil
}

func (a *AccessCache) Delete(ctx context.Context, playlistID uuid.UUID) error {
	if err := a.client.Del(ctx, a.key(playlistID)).Err(); err != nil {
		return fmt.Errorf("failed to delete playlist access meta from redis: %w", err)
	}
	return nil
}
