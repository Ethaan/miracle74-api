package repo

import (
	"context"
	"fmt"

	"github.com/ethaan/miracle74-api/internal/types"
	"github.com/ethaan/miracle74-api/pkg/cache"
)

type CharacterRepo struct {
	cache *cache.Client
}

func NewCharacterRepo(cacheClient *cache.Client) *CharacterRepo {
	return &CharacterRepo{
		cache: cacheClient,
	}
}

func (r *CharacterRepo) Get(ctx context.Context, name string) (*types.Character, error) {
	key := r.buildKey(name)

	var character types.Character
	if err := r.cache.Get(ctx, key, &character); err != nil {
		return nil, err
	}

	return &character, nil
}

func (r *CharacterRepo) Set(ctx context.Context, name string, character *types.Character) error {
	key := r.buildKey(name)
	return r.cache.Set(ctx, key, character)
}

func (r *CharacterRepo) Delete(ctx context.Context, name string) error {
	key := r.buildKey(name)
	return r.cache.Delete(ctx, key)
}

func (r *CharacterRepo) buildKey(name string) string {
	return fmt.Sprintf("character:%s", name)
}
