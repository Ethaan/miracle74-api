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

type GuildService struct {
	client *miracle74.Client
	repo   *repo.GuildRepo
}

func NewGuildService(guildRepo *repo.GuildRepo) *GuildService {
	return &GuildService{
		client: miracle74.NewClient(),
		repo:   guildRepo,
	}
}

func (s *GuildService) GetGuild(ctx context.Context, guildID int) (*types.Guild, error) {
	guild, err := s.repo.Get(ctx, guildID)
	if err == nil {
		log.Printf("Cache hit for guild:%d", guildID)
		return guild, nil
	}

	if !errors.Is(err, cache.ErrCacheMiss) {
		log.Printf("Cache error: %v", err)
	} else {
		log.Printf("Cache miss for guild:%d", guildID)
	}

	guild, err = s.client.ScrapeGuild(guildID)
	if err != nil {
		return nil, fmt.Errorf("failed to scrape guild: %w", err)
	}

	if err := s.repo.Set(ctx, guildID, guild); err != nil {
		log.Printf("Failed to cache guild: %v", err)
	} else {
		log.Printf("Cached guild:%d", guildID)
	}

	return guild, nil
}
