package repo

import (
	"context"
	"time"

	"github.com/ethaan/miracle74-api/internal/types"
	"github.com/ethaan/miracle74-api/pkg/cache"
)

const (
	PowerGamersTTL = 1 * time.Minute
)

type PowerGamersRepo struct {
	cache *cache.Client
}

func NewPowerGamersRepo(cacheClient *cache.Client) *PowerGamersRepo {
	return &PowerGamersRepo{
		cache: cacheClient,
	}
}

func (r *PowerGamersRepo) Get(ctx context.Context, includeAll bool) ([]types.PowerGamer, error) {
	key := r.buildKey(includeAll)

	var powerGamers []types.PowerGamer
	if err := r.cache.Get(ctx, key, &powerGamers); err != nil {
		return nil, err
	}

	return powerGamers, nil
}

func (r *PowerGamersRepo) Set(ctx context.Context, powerGamers []types.PowerGamer, includeAll bool) error {
	key := r.buildKey(includeAll)
	return r.cache.SetWithTTL(ctx, key, powerGamers, PowerGamersTTL)
}

func (r *PowerGamersRepo) Delete(ctx context.Context, includeAll bool) error {
	key := r.buildKey(includeAll)
	return r.cache.Delete(ctx, key)
}

func (r *PowerGamersRepo) buildKey(includeAll bool) string {
	if includeAll {
		return "powergamers:today:all"
	}
	return "powergamers:today:page:1"
}
