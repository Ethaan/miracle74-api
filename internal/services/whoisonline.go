package services

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/ethaan/miracle74-api/internal/repo"
	"github.com/ethaan/miracle74-api/internal/types"
	"github.com/ethaan/miracle74-api/pkg/cache"
	"github.com/ethaan/miracle74-api/pkg/miracle74"
)

type WhoIsOnlineService struct {
	client *miracle74.Client
	repo   *repo.WhoIsOnlineRepo
}

func NewWhoIsOnlineService(whoIsOnlineRepo *repo.WhoIsOnlineRepo) *WhoIsOnlineService {
	return &WhoIsOnlineService{
		client: miracle74.NewClient(),
		repo:   whoIsOnlineRepo,
	}
}

func (s *WhoIsOnlineService) GetWhoIsOnline(ctx context.Context, order string) ([]types.OnlinePlayer, error) {
	onlinePlayers, err := s.repo.Get(ctx, order)
	if err == nil {
		cacheKey := s.repo.BuildKey(order)
		log.Printf("Cache hit for %s", cacheKey)
		return onlinePlayers, nil
	}

	if !errors.Is(err, cache.ErrCacheMiss) {
		log.Printf("Cache error: %v", err)
	} else {
		cacheKey := s.repo.BuildKey(order)
		log.Printf("Cache miss for %s", cacheKey)
	}

	onlinePlayers, err = s.client.ScrapeWhoIsOnline(order)
	if err != nil {
		return nil, fmt.Errorf("failed to scrape who is online: %w", err)
	}

	if err := s.repo.Set(ctx, onlinePlayers, order); err != nil {
		log.Printf("Failed to cache who is online: %v", err)
	} else {
		log.Printf("Cached %d online players", len(onlinePlayers))
	}

	return onlinePlayers, nil
}
