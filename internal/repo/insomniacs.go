package repo

import (
	"context"
	"time"

	"github.com/ethaan/miracle74-api/internal/types"
	"github.com/ethaan/miracle74-api/pkg/cache"
)

const (
	InsomniacsTTL = 5 * time.Minute
)

type InsomniacsRepo struct {
	cache *cache.Client
}

func NewInsomniacsRepo(cacheClient *cache.Client) *InsomniacsRepo {
	return &InsomniacsRepo{
		cache: cacheClient,
	}
}

func (r *InsomniacsRepo) Get(ctx context.Context, includeAll bool) ([]types.Insomniac, error) {
	key := r.buildKey(includeAll)

	var insomniacs []types.Insomniac
	if err := r.cache.Get(ctx, key, &insomniacs); err != nil {
		return nil, err
	}

	return insomniacs, nil
}

func (r *InsomniacsRepo) Set(ctx context.Context, insomniacs []types.Insomniac, includeAll bool) error {
	key := r.buildKey(includeAll)
	return r.cache.SetWithTTL(ctx, key, insomniacs, InsomniacsTTL)
}

func (r *InsomniacsRepo) Delete(ctx context.Context, includeAll bool) error {
	key := r.buildKey(includeAll)
	return r.cache.Delete(ctx, key)
}

func (r *InsomniacsRepo) buildKey(includeAll bool) string {
	if includeAll {
		return "insomniacs:all"
	}
	return "insomniacs:page:1"
}
