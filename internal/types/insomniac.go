package types

type Insomniac struct {
	Rank       int    `json:"rank"`
	Name       string `json:"name"`
	Country    string `json:"country,omitempty"`
	Vocation   string `json:"vocation"`
	Level      int    `json:"level"`
	TimeOnline string `json:"time_online"`
}
