package types

import "time"

type Character struct {
	Name      string     `json:"name"`
	Sex       string     `json:"sex"`
	Vocation  string     `json:"vocation,omitempty"`
	Level     int        `json:"level,omitempty"`
	Residence string     `json:"residence,omitempty"`
	Guild     string     `json:"guild,omitempty"`
	GuildRank string     `json:"guild_rank,omitempty"`
	GuildURL  string     `json:"guild_url,omitempty"`
	LastLogin *time.Time `json:"last_login,omitempty"`
	IsPremium bool       `json:"is_premium"`
	Country   string     `json:"country,omitempty"`
	Deaths    []Death    `json:"deaths,omitempty"`
}

type Death struct {
	Date     string `json:"date"`
	Level    int    `json:"level"`
	KilledBy string `json:"killed_by"`
}

