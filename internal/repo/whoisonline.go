package repo

import (
	"context"
	"time"

	"github.com/ethaan/miracle74-api/internal/types"
	"github.com/ethaan/miracle74-api/pkg/cache"
)

const (
	WhoIsOnlineTTL = 15 * time.Second
)

type WhoIsOnlineRepo struct {
	cache *cache.Client
}

func NewWhoIsOnlineRepo(cacheClient *cache.Client) *WhoIsOnlineRepo {
	return &WhoIsOnlineRepo{
		cache: cacheClient,
	}
}

func (r *WhoIsOnlineRepo) Get(ctx context.Context, order string) ([]types.OnlinePlayer, error) {
	key := r.BuildKey(order)

	var onlinePlayers []types.OnlinePlayer
	if err := r.cache.Get(ctx, key, &onlinePlayers); err != nil {
		return nil, err
	}

	return onlinePlayers, nil
}

func (r *WhoIsOnlineRepo) Set(ctx context.Context, onlinePlayers []types.OnlinePlayer, order string) error {
	key := r.BuildKey(order)
	return r.cache.SetWithTTL(ctx, key, onlinePlayers, WhoIsOnlineTTL)
}

func (r *WhoIsOnlineRepo) Delete(ctx context.Context, order string) error {
	key := r.BuildKey(order)
	return r.cache.Delete(ctx, key)
}

func (r *WhoIsOnlineRepo) BuildKey(order string) string {
	orderKey := "name"
	if order != "" {
		orderKey = order
	}

	return "whoisonline:" + orderKey
}
