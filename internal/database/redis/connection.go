package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisConnection struct {
	client *redis.Client
}

func NewRedisConnection(ctx context.Context, url, password string, db int) (*RedisConnection, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: password,
		DB:       db,
	})

	// Test the connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisConnection{
		client: client,
	}, nil
}

func (rc *RedisConnection) Close() error {
	return rc.client.Close()
}

func (rc *RedisConnection) Client() *redis.Client {
	return rc.client
}
