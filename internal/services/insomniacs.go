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

type InsomniacsService struct {
	client *miracle74.Client
	repo   *repo.InsomniacsRepo
}

func NewInsomniacsService(insomniacsRepo *repo.InsomniacsRepo) *InsomniacsService {
	return &InsomniacsService{
		client: miracle74.NewClient(),
		repo:   insomniacsRepo,
	}
}

func (s *InsomniacsService) GetInsomniacs(ctx context.Context, includeAll bool) ([]types.Insomniac, error) {
	insomniacs, err := s.repo.Get(ctx, includeAll)
	if err == nil {
		cacheKey := "insomniacs:page:1"
		if includeAll {
			cacheKey = "insomniacs:all"
		}
		log.Printf("Cache hit for %s", cacheKey)
		return insomniacs, nil
	}

	if !errors.Is(err, cache.ErrCacheMiss) {
		log.Printf("Cache error: %v", err)
	} else {
		cacheKey := "insomniacs:page:1"
		if includeAll {
			cacheKey = "insomniacs:all"
		}
		log.Printf("Cache miss for %s", cacheKey)
	}

	insomniacs, err = s.client.ScrapeInsomniacs(includeAll)
	if err != nil {
		return nil, fmt.Errorf("failed to scrape insomniacs: %w", err)
	}

	if err := s.repo.Set(ctx, insomniacs, includeAll); err != nil {
		log.Printf("Failed to cache insomniacs: %v", err)
	} else {
		log.Printf("Cached %d insomniacs", len(insomniacs))
	}

	return insomniacs, nil
}
