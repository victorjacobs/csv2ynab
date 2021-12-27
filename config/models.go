package config

type Config struct {
	Ynab             Ynab             `json:"ynab"`
	WatchDirectories []WatchDirectory `json:"watch_directories"`
}

type Ynab struct {
	ApiKey    string `json:"api_key"`
	BudgetId  string `json:"budget_id"`
	AccountId string `json:"account_id"`
}

type WatchDirectory struct {
	Path      string `json:"path"`
	BudgetId  string `json:"budget_id"`
	AccountId string `json:"account_id"`
}

func (w *WatchDirectory) Merge(c Ynab) Ynab {
	return Ynab{
		ApiKey:    c.ApiKey,
		BudgetId:  w.BudgetId,
		AccountId: w.AccountId,
	}
}
