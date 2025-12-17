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

type CharacterService struct {
	client *miracle74.Client
	repo   *repo.CharacterRepo
}

func NewCharacterService(characterRepo *repo.CharacterRepo) *CharacterService {
	return &CharacterService{
		client: miracle74.NewClient(),
		repo:   characterRepo,
	}
}

func (s *CharacterService) GetCharacter(ctx context.Context, name string) (*types.Character, error) {
	character, err := s.repo.Get(ctx, name)
	if err == nil {
		log.Printf("Cache hit for character:%s", name)
		return character, nil
	}

	if !errors.Is(err, cache.ErrCacheMiss) {
		log.Printf("Cache error: %v", err)
	} else {
		log.Printf("Cache miss for character:%s", name)
	}

	character, err = s.client.ScrapeCharacter(name)
	if err != nil {
		return nil, fmt.Errorf("failed to scrape character: %w", err)
	}

	if err := s.repo.Set(ctx, name, character); err != nil {
		log.Printf("Failed to cache character: %v", err)
	} else {
		log.Printf("Cached character:%s", name)
	}

	return character, nil
}
