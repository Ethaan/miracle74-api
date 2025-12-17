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

func (s *PowerGamersService) GetPowerGamers(ctx context.Context, includeAll bool) ([]types.PowerGamer, error) {
	powerGamers, err := s.repo.Get(ctx, includeAll)
	if err == nil {
		cacheKey := "powergamers:today:page:1"
		if includeAll {
			cacheKey = "powergamers:today:all"
		}
		log.Printf("Cache hit for %s", cacheKey)
		return powerGamers, nil
	}

	if !errors.Is(err, cache.ErrCacheMiss) {
		log.Printf("Cache error: %v", err)
	} else {
		cacheKey := "powergamers:today:page:1"
		if includeAll {
			cacheKey = "powergamers:today:all"
		}
		log.Printf("Cache miss for %s", cacheKey)
	}

	powerGamers, err = s.client.ScrapePowerGamers(includeAll)
	if err != nil {
		return nil, fmt.Errorf("failed to scrape power gamers: %w", err)
	}

	if err := s.repo.Set(ctx, powerGamers, includeAll); err != nil {
		log.Printf("Failed to cache power gamers: %v", err)
	} else {
		log.Printf("Cached %d power gamers", len(powerGamers))
	}

	return powerGamers, nil
}
