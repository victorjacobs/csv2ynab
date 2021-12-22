package config

type Config struct {
	Ynab Ynab `json:"ynab"`
}

type Ynab struct {
	ApiKey string `json:"api_key"`
}
