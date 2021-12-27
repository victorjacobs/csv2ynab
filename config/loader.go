package config

import (
	"encoding/json"
	"errors"
	"os"
	"os/user"
	"path/filepath"
)

func Load() (Config, error) {
	config := Config{}

	usr, err := user.Current()
	if err != nil {
		return config, err
	}

	dir := usr.HomeDir

	path := filepath.Join(dir, ".ynabrc")

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
	for _, watchDirectory := range c.WatchDirectories {
		if err := validateWatchDirectory(watchDirectory); err != nil {
			return err
		}
	}

	return nil
}

func validateWatchDirectory(d WatchDirectory) error {
	if d.AccountId == "" {
		return errors.New("account_id not set for watch_directory")
	}

	if d.Path == "" {
		return errors.New("path not set for watch_directory")
	}

	if d.BudgetId == "" {
		return errors.New("budget_id not set for watch_directory")
	}

	return nil
}
