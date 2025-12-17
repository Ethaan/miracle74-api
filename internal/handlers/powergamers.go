package handlers

import (
	"context"

	"github.com/ethaan/miracle74-api/internal/api"
)

func (h *Handler) GetPowerGamers(ctx context.Context, params api.GetPowerGamersParams) (api.GetPowerGamersRes, error) {
	includeAll := params.IncludeAll.Value

	powerGamers, err := h.powerGamersService.GetPowerGamers(ctx, includeAll)
	if err != nil {
		return &api.ErrorResponse{
			Error:   "fetch_failed",
			Message: err.Error(),
		}, nil
	}

	var apiPowerGamers []api.PowerGamer
	for _, pg := range powerGamers {
		apiPowerGamers = append(apiPowerGamers, api.PowerGamer{
			Rank:     pg.Rank,
			Name:     pg.Name,
			Vocation: pg.Vocation,
			Level:    pg.Level,
			Today:    pg.Today,
		})
	}

	return &api.PowerGamersResponse{
		PowerGamers: apiPowerGamers,
		Total:       len(apiPowerGamers),
	}, nil
}
