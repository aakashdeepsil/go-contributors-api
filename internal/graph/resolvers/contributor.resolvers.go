package resolvers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aakashdeepsil/go-contributors-api/internal/database/models"
	"github.com/aakashdeepsil/go-contributors-api/internal/graph/generated"
	"github.com/aakashdeepsil/go-contributors-api/internal/graph/model"
)

const (
	cacheTTL        = time.Hour
	cacheCtxTimeout = 5 * time.Second
)

// CreateContributor is the resolver for the createContributor field.
func (r *mutationResolver) CreateContributor(ctx context.Context, input model.ContributorInput) (*model.Contributor, error) {
	contributorInput := &models.ContributorInput{
		Username: input.Username,
		Email:    input.Email,
		Name:     input.Name,
		Projects: input.Projects,
	}

	if input.AvatarURL != nil {
		contributorInput.AvatarURL = *input.AvatarURL
	}

	contributor, err := r.ContributorRepo.Create(ctx, contributorInput)
	if err != nil {
		return nil, err
	}

	// Cache the new contributor
	r.cacheContributor(ctx, contributorCacheKey(contributor.ID.Hex()), contributor)

	return mapDBToGraphQL(contributor), nil
}

// UpdateContributor is the resolver for the updateContributor field.
func (r *mutationResolver) UpdateContributor(ctx context.Context, id string, input model.UpdateContributorInput) (*model.Contributor, error) {
	contributorInput := &models.ContributorInput{
		Username: *input.Username,
		Email:    *input.Email,
		Name:     *input.Name,
		Projects: input.Projects,
	}

	if input.AvatarURL != nil {
		contributorInput.AvatarURL = *input.AvatarURL
	}

	contributor, err := r.ContributorRepo.Update(ctx, id, contributorInput)
	if err != nil {
		return nil, err
	}

	// Invalidate cache
	cacheKey := contributorCacheKey(id)
	go r.CacheRepo.Delete(ctx, cacheKey)

	return mapDBToGraphQL(contributor), nil
}

// DeleteContributor is the resolver for the deleteContributor field.
func (r *mutationResolver) DeleteContributor(ctx context.Context, id string) (bool, error) {
	err := r.ContributorRepo.Delete(ctx, id)
	if err != nil {
		return false, err
	}

	// Invalidate cache
	cacheKey := contributorCacheKey(id)
	go r.CacheRepo.Delete(ctx, cacheKey)

	return true, nil
}

// Contributor is the resolver for the contributor field.
func (r *queryResolver) Contributor(ctx context.Context, id string) (*model.Contributor, error) {
	cacheKey := contributorCacheKey(id)

	// Try cache first
	if cached, err := r.CacheRepo.Get(ctx, cacheKey); err == nil {
		var contributor models.Contributor
		if err := json.Unmarshal([]byte(cached), &contributor); err == nil {
			return mapDBToGraphQL(&contributor), nil
		}
	}

	// Get from database
	contributor, err := r.ContributorRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cache in background
	r.cacheContributor(ctx, cacheKey, contributor)

	return mapDBToGraphQL(contributor), nil
}

// ContributorByUsername is the resolver for the contributorByUsername field.
func (r *queryResolver) ContributorByUsername(ctx context.Context, username string) (*model.Contributor, error) {
	cacheKey := fmt.Sprintf("contributor:username:%s", username)

	// Try cache first
	if cached, err := r.CacheRepo.Get(ctx, cacheKey); err == nil {
		var contributor models.Contributor
		if err := json.Unmarshal([]byte(cached), &contributor); err == nil {
			return mapDBToGraphQL(&contributor), nil
		}
	}

	contributor, err := r.ContributorRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	// Cache in background
	r.cacheContributor(ctx, cacheKey, contributor)

	return mapDBToGraphQL(contributor), nil
}

// Contributors is the resolver for the contributors field.
func (r *queryResolver) Contributors(ctx context.Context, limit *int, offset *int) ([]*model.Contributor, error) {
	defaultLimit := 10
	defaultOffset := 0

	if limit != nil {
		defaultLimit = *limit
	}
	if offset != nil {
		defaultOffset = *offset
	}

	contributors, err := r.ContributorRepo.List(ctx, defaultLimit, defaultOffset)
	if err != nil {
		return nil, err
	}

	result := make([]*model.Contributor, len(contributors))
	for i, contributor := range contributors {
		result[i] = mapDBToGraphQL(contributor)
	}

	return result, nil
}

// ContributorUpdated is the resolver for the contributorUpdated field.
func (r *subscriptionResolver) ContributorUpdated(ctx context.Context) (<-chan *model.Contributor, error) {
	ch := make(chan *model.Contributor, 1)

	// This is a basic implementation. You might want to use a proper pub/sub system
	// like Redis pub/sub for a production environment
	go func() {
		<-ctx.Done()
		close(ch)
	}()

	return ch, nil
}

func mapDBToGraphQL(c *models.Contributor) *model.Contributor {
	if c == nil {
		return nil
	}

	var avatarURL *string
	if c.AvatarURL != "" {
		avatarURL = &c.AvatarURL
	}

	return &model.Contributor{
		ID:        c.ID.Hex(),
		Username:  c.Username,
		Email:     c.Email,
		Name:      c.Name,
		AvatarURL: avatarURL,
		Projects:  c.Projects,
		JoinedAt:  c.JoinedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func contributorCacheKey(id string) string {
	return fmt.Sprintf("contributor:%s", id)
}

func (r *Resolver) cacheContributor(_ context.Context, key string, contributor *models.Contributor) {
	cacheData, err := json.Marshal(contributor)
	if err != nil {
		return // Silent fail for cache operations
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), cacheCtxTimeout)
		defer cancel()
		r.CacheRepo.Set(ctx, key, string(cacheData), cacheTTL)
	}()
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
