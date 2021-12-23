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

	return config, nil
}
