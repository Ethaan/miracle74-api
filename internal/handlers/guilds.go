package handlers

import (
	"context"

	"github.com/ethaan/miracle74-api/internal/api"
)

func (h *Handler) GetGuild(ctx context.Context, params api.GetGuildParams) (api.GetGuildRes, error) {
	guild, err := h.guildService.GetGuild(ctx, params.GuildId)
	if err != nil {
		return &api.GetGuildInternalServerError{
			Error:   "fetch_failed",
			Message: err.Error(),
		}, nil
	}

	var members []api.GuildMember
	for _, m := range guild.Members {
		members = append(members, api.GuildMember{
			Rank:     m.Rank,
			Name:     m.Name,
			Vocation: m.Vocation,
			Level:    m.Level,
			Status:   m.Status,
		})
	}

	response := &api.GuildResponse{
		GuildID: guild.GuildID,
		Members: members,
		Total:   len(members),
	}

	return response, nil
}
