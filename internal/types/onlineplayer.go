package types

type OnlinePlayer struct {
	Name     string `json:"name"`
	Level    int    `json:"level"`
	Vocation string `json:"vocation"`
	Country  string `json:"country"`
}
