package repo

import (
	"context"
	"fmt"

	"github.com/ethaan/miracle74-api/internal/types"
	"github.com/ethaan/miracle74-api/pkg/cache"
)

type GuildRepo struct {
	cache *cache.Client
}

func NewGuildRepo(cacheClient *cache.Client) *GuildRepo {
	return &GuildRepo{
		cache: cacheClient,
	}
}

func (r *GuildRepo) Get(ctx context.Context, guildID int) (*types.Guild, error) {
	key := r.buildKey(guildID)

	var guild types.Guild
	if err := r.cache.Get(ctx, key, &guild); err != nil {
		return nil, err
	}

	return &guild, nil
}

func (r *GuildRepo) Set(ctx context.Context, guildID int, guild *types.Guild) error {
	key := r.buildKey(guildID)
	return r.cache.Set(ctx, key, guild)
}

func (r *GuildRepo) Delete(ctx context.Context, guildID int) error {
	key := r.buildKey(guildID)
	return r.cache.Delete(ctx, key)
}

func (r *GuildRepo) buildKey(guildID int) string {
	return fmt.Sprintf("guild:%d", guildID)
}
