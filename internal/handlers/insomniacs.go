package handlers

import (
	"context"

	"github.com/ethaan/miracle74-api/internal/api"
)

func (h *Handler) GetInsomniacs(ctx context.Context, params api.GetInsomniacsParams) (api.GetInsomniacsRes, error) {
	includeAll := params.IncludeAll.Value

	insomniacs, err := h.insomniacsService.GetInsomniacs(ctx, includeAll)
	if err != nil {
		return &api.ErrorResponse{
			Error:   "fetch_failed",
			Message: err.Error(),
		}, nil
	}

	var apiInsomniacs []api.Insomniac
	for _, ins := range insomniacs {
		apiInsomniac := api.Insomniac{
			Rank:       ins.Rank,
			Name:       ins.Name,
			Vocation:   ins.Vocation,
			Level:      ins.Level,
			TimeOnline: ins.TimeOnline,
		}
		if ins.Country != "" {
			apiInsomniac.Country.SetTo(ins.Country)
		}
		apiInsomniacs = append(apiInsomniacs, apiInsomniac)
	}

	return &api.InsomniacsResponse{
		Insomniacs: apiInsomniacs,
		Total:      len(apiInsomniacs),
	}, nil
}
