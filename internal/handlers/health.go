package handlers

import (
	"context"
	"time"

	"github.com/ethaan/miracle74-api/internal/api"
	"github.com/ethaan/miracle74-api/internal/services"
)

type Handler struct {
	characterService   *services.CharacterService
	powerGamersService *services.PowerGamersService
	insomniacsService  *services.InsomniacsService
}

func NewHandler(characterService *services.CharacterService, powerGamersService *services.PowerGamersService, insomniacsService *services.InsomniacsService) *Handler {
	return &Handler{
		characterService:   characterService,
		powerGamersService: powerGamersService,
		insomniacsService:  insomniacsService,
	}
}

func (h *Handler) GetHealth(ctx context.Context) (*api.HealthResponse, error) {
	return &api.HealthResponse{
		Status:    api.HealthResponseStatusHealthy,
		Timestamp: time.Now().UTC(),
		Version:   "0.1.0",
	}, nil
}
