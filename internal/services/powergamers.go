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

type PowerGamersService struct {
	client *miracle74.Client
	repo   *repo.PowerGamersRepo
}

func NewPowerGamersService(powerGamersRepo *repo.PowerGamersRepo) *PowerGamersService {
	return &PowerGamersService{
		client: miracle74.NewClient(),
		repo:   powerGamersRepo,
	}
}

func (s *PowerGamersService) GetPowerGamers(ctx context.Context, includeAll bool, list string, vocation string) ([]types.PowerGamer, error) {
	powerGamers, err := s.repo.Get(ctx, includeAll, list, vocation)
	if err == nil {
		cacheKey := s.repo.BuildKey(includeAll, list, vocation)
		log.Printf("Cache hit for %s", cacheKey)
		return powerGamers, nil
	}

	if !errors.Is(err, cache.ErrCacheMiss) {
		log.Printf("Cache error: %v", err)
	} else {
		cacheKey := s.repo.BuildKey(includeAll, list, vocation)
		log.Printf("Cache miss for %s", cacheKey)
	}

	powerGamers, err = s.client.ScrapePowerGamers(includeAll, list, vocation)
	if err != nil {
		return nil, fmt.Errorf("failed to scrape power gamers: %w", err)
	}

	if err := s.repo.Set(ctx, powerGamers, includeAll, list, vocation); err != nil {
		log.Printf("Failed to cache power gamers: %v", err)
	} else {
		log.Printf("Cached %d power gamers", len(powerGamers))
	}

	return powerGamers, nil
}
