package config

import (
	"encoding/json"
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
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(contents, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
