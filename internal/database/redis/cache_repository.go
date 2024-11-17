package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheRepository struct {
	client *redis.Client
}

func NewCacheRepository(client *redis.Client) *CacheRepository {
	return &CacheRepository{
		client: client,
	}
}

func (r *CacheRepository) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *CacheRepository) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *CacheRepository) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *CacheRepository) IncrementCounter(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

func (r *CacheRepository) ResetCounter(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
