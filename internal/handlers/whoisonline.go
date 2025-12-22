package handlers

import (
	"context"

	"github.com/ethaan/miracle74-api/internal/api"
)

func (h *Handler) GetWhoIsOnline(ctx context.Context, params api.GetWhoIsOnlineParams) (api.GetWhoIsOnlineRes, error) {
	order := string(params.Order.Value)

	onlinePlayers, err := h.whoIsOnlineService.GetWhoIsOnline(ctx, order)
	if err != nil {
		return &api.ErrorResponse{
			Error:   "fetch_failed",
			Message: err.Error(),
		}, nil
	}

	var apiOnlinePlayers []api.OnlinePlayer
	for _, player := range onlinePlayers {
		apiOnlinePlayers = append(apiOnlinePlayers, api.OnlinePlayer{
			Name:     player.Name,
			Level:    player.Level,
			Vocation: player.Vocation,
			Country:  api.NewOptString(player.Country),
		})
	}

	return &api.WhoIsOnlineResponse{
		Players: apiOnlinePlayers,
		Total:   len(apiOnlinePlayers),
	}, nil
}
