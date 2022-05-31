package config

type Config struct {
	Ynab  Ynab  `json:"ynab"`
	Watch Watch `json:"watch"`
}

type Ynab struct {
	ApiKey    string `json:"api_key"`
	BudgetId  string `json:"budget_id"`
	AccountId string `json:"account_id"`
}

type Watch struct {
	Directory string         `json:"directory"`
	Patterns  []WatchPattern `json:"patterns"`
}

type WatchPattern struct {
	Pattern   string `json:"pattern"`
	BudgetId  string `json:"budget_id"`
	AccountId string `json:"account_id"`
}

func (w *WatchPattern) Merge(c Ynab) Ynab {
	return Ynab{
		ApiKey:    c.ApiKey,
		BudgetId:  w.BudgetId,
		AccountId: w.AccountId,
	}
}
