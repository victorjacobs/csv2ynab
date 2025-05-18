package config

import (
	"encoding/json"
	"errors"
	"os"
	"os/user"
	"path/filepath"
)

func Load(configPath string) (Config, error) {
	config := Config{}

	var path string
	if configPath == "" {
		usr, err := user.Current()
		if err != nil {
			return config, err
		}

		dir := usr.HomeDir

		path = filepath.Join(dir, ".ynabrc")
	} else {
		path = configPath
	}

	contents, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return config, nil
	}

	if err != nil {
		return config, err
	}

	err = json.Unmarshal(contents, &config)
	if err != nil {
		return config, err
	}

	if err := validate(config); err != nil {
		return config, err
	}

	return config, nil
}

func validate(c Config) error {
	for _, watchPattern := range c.Watch.Patterns {
		if err := validateWatchPattern(watchPattern); err != nil {
			return err
		}
	}

	return nil
}

func validateWatchPattern(w *WatchPattern) error {
	if w.AccountId == "" {
		return errors.New("account_id not set for pattern")
	}

	if w.Pattern == "" {
		return errors.New("pattern not set for pattern")
	}

	if w.BudgetId == "" {
		return errors.New("budget_id not set for pattern")
	}

	return nil
}
