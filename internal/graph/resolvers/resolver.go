package resolvers

import (
	"github.com/aakashdeepsil/go-contributors-api/internal/database/repository"
)

type Resolver struct {
	ContributorRepo repository.ContributorRepository
	CacheRepo       repository.CacheRepository
}

func NewResolver(contributorRepo repository.ContributorRepository, cacheRepo repository.CacheRepository) *Resolver {
	return &Resolver{
		ContributorRepo: contributorRepo,
		CacheRepo:       cacheRepo,
	}
}
