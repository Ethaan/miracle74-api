package repo

import (
	"context"
	"time"

	"github.com/ethaan/miracle74-api/internal/types"
	"github.com/ethaan/miracle74-api/pkg/cache"
)

const (
	PowerGamersTTL = 5 * time.Minute
)

type PowerGamersRepo struct {
	cache *cache.Client
}

func NewPowerGamersRepo(cacheClient *cache.Client) *PowerGamersRepo {
	return &PowerGamersRepo{
		cache: cacheClient,
	}
}

func (r *PowerGamersRepo) Get(ctx context.Context, includeAll bool, list string, vocation string) ([]types.PowerGamer, error) {
	key := r.BuildKey(includeAll, list, vocation)

	var powerGamers []types.PowerGamer
	if err := r.cache.Get(ctx, key, &powerGamers); err != nil {
		return nil, err
	}

	return powerGamers, nil
}

func (r *PowerGamersRepo) Set(ctx context.Context, powerGamers []types.PowerGamer, includeAll bool, list string, vocation string) error {
	key := r.BuildKey(includeAll, list, vocation)
	return r.cache.SetWithTTL(ctx, key, powerGamers, PowerGamersTTL)
}

func (r *PowerGamersRepo) Delete(ctx context.Context, includeAll bool, list string, vocation string) error {
	key := r.BuildKey(includeAll, list, vocation)
	return r.cache.Delete(ctx, key)
}

func (r *PowerGamersRepo) BuildKey(includeAll bool, list string, vocation string) string {
	scope := "page:1"
	if includeAll {
		scope = "all"
	}

	vocationKey := "all"
	if vocation != "" {
		vocationKey = vocation
	}

	return "powergamers:" + list + ":" + vocationKey + ":" + scope
}
