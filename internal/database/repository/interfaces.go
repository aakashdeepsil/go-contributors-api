package repository

import (
	"context"
	"time"

	"github.com/aakashdeepsil/go-contributors-api/internal/database/models"
)

type ContributorRepository interface {
	Create(ctx context.Context, contributor *models.ContributorInput) (*models.Contributor, error)
	GetByID(ctx context.Context, id string) (*models.Contributor, error)
	GetByUsername(ctx context.Context, username string) (*models.Contributor, error)
	Update(ctx context.Context, id string, contributor *models.ContributorInput) (*models.Contributor, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*models.Contributor, error)
}

type CacheRepository interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	IncrementCounter(ctx context.Context, key string) (int64, error)
	ResetCounter(ctx context.Context, key string) error
}
