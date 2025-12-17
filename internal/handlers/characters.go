package handlers

import (
	"context"

	"github.com/ethaan/miracle74-api/internal/api"
)

func (h *Handler) GetCharacter(ctx context.Context, params api.GetCharacterParams) (api.GetCharacterRes, error) {
	character, err := h.characterService.GetCharacter(ctx, params.Name)
	if err != nil {
		return &api.GetCharacterInternalServerError{
			Error:   "fetch_failed",
			Message: err.Error(),
		}, nil
	}

	var deaths []api.Death
	for _, d := range character.Deaths {
		deaths = append(deaths, api.Death{
			Date:     d.Date,
			Level:    d.Level,
			KilledBy: d.KilledBy,
		})
	}

	response := &api.CharacterResponse{
		Name:      character.Name,
		Sex:       character.Sex,
		IsPremium: character.IsPremium,
	}

	if character.Vocation != "" {
		response.Vocation.SetTo(character.Vocation)
	}
	if character.Level > 0 {
		response.Level.SetTo(character.Level)
	}
	if character.Residence != "" {
		response.Residence.SetTo(character.Residence)
	}
	if character.Guild != "" {
		response.Guild.SetTo(character.Guild)
	}
	if character.GuildRank != "" {
		response.GuildRank.SetTo(character.GuildRank)
	}
	if character.GuildURL != "" {
		response.GuildURL.SetTo(character.GuildURL)
	}
	if character.LastLogin != nil {
		response.LastLogin.SetTo(*character.LastLogin)
	}
	if character.Country != "" {
		response.Country.SetTo(character.Country)
	}

	response.Deaths = deaths

	return response, nil
}
