package types

type Guild struct {
	GuildID int           `json:"guild_id"`
	Members []GuildMember `json:"members"`
}

type GuildMember struct {
	Rank     string `json:"rank"`
	Name     string `json:"name"`
	Vocation string `json:"vocation"`
	Level    int    `json:"level"`
	Status   string `json:"status"`
}
