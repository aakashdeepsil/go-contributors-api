package database

import (
	"context"

	"github.com/aakashdeepsil/go-contributors-api/internal/config"
	"github.com/aakashdeepsil/go-contributors-api/internal/database/mongodb"
	"github.com/aakashdeepsil/go-contributors-api/internal/database/redis"
	"github.com/aakashdeepsil/go-contributors-api/internal/database/repository"
)

type Repositories struct {
	Contributor repository.ContributorRepository
	Cache       repository.CacheRepository
}

func NewRepositories(ctx context.Context, cfg *config.Config) (*Repositories, error) {
	// Initialize MongoDB connection
	mongoConn, err := mongodb.NewMongoDBConnection(ctx, cfg.MongoDB.URI, cfg.MongoDB.Database)
	if err != nil {
		return nil, err
	}

	// Initialize Redis connection
	redisConn, err := redis.NewRedisConnection(ctx, cfg.Redis.URL, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		return nil, err
	}

	return &Repositories{
		Contributor: mongodb.NewContributorRepository(mongoConn.Database()),
		Cache:       redis.NewCacheRepository(redisConn.Client()),
	}, nil
}
