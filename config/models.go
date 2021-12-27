package config

type Config struct {
	Ynab             Ynab             `json:"ynab"`
	WatchDirectories []WatchDirectory `json:"watch_directories"`
}

type Ynab struct {
	ApiKey string `json:"api_key"`
}

type WatchDirectory struct {
	Path string `json:"path"`
}
