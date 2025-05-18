package config

import (
	"errors"
	"fmt"
	"path/filepath"
)

type Config struct {
	YNAB  *YNAB  `json:"ynab"`
	Watch *Watch `json:"watch"`
}

func (c *Config) YNABConfigForFile(path string) (*YNAB, error) {
	baseFileName := filepath.Base(path)

	var matchedPattern *WatchPattern
	for _, p := range c.Watch.Patterns {
		match, err := filepath.Match(p.Pattern, baseFileName)
		if err != nil {
			return nil, fmt.Errorf("Matching failed: %w", err)
		}

		if match {
			matchedPattern = p
			break
		}
	}

	if matchedPattern == nil {
		return nil, errors.New("didn't find matching config")
	}

	return matchedPattern.Merge(c.YNAB), nil
}

type YNAB struct {
	ApiKey    string `json:"api_key"`
	BudgetId  string `json:"budget_id"`
	AccountId string `json:"account_id"`
}

type Watch struct {
	Directory string          `json:"directory"`
	Patterns  []*WatchPattern `json:"patterns"`
}

type WatchPattern struct {
	Pattern   string `json:"pattern"`
	BudgetId  string `json:"budget_id"`
	AccountId string `json:"account_id"`
}

func (w *WatchPattern) Merge(c *YNAB) *YNAB {
	return &YNAB{
		ApiKey:    c.ApiKey,
		BudgetId:  w.BudgetId,
		AccountId: w.AccountId,
	}
}
