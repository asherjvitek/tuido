package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"tuido/util"
)

type Config struct {
	StorageType string `json:"storage_type"`
	RemoteUrl   string `json:"remote_url"`
}

func Load() (Config, error) {
	var config Config
	dir, err := util.GetAppDir()

	if err != nil {
		return config, err
	}

	path := filepath.Join(dir, "config.json")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return config, fmt.Errorf("failed to create directory for database: %w", err)
		}

		config = getDefault()

		conf, err := json.Marshal(config)

		if err != nil {
			return config, fmt.Errorf("failed to marshal default config: %w", err)
		}

		if err := os.WriteFile(path, conf, 0644); err != nil {
			return config, fmt.Errorf("failed to write default config: %w", err)
		}

		return config, nil
	}

	var contents []byte
	if contents, err = os.ReadFile(path); err != nil {
		err := json.Unmarshal(contents, &config)

		if err != nil {
			return config, fmt.Errorf("failed to unmarshal config: %w", err)
		}

		return config, nil
	}

	return config, err
}

func Save(config Config) error {
	dir, err := util.GetAppDir()
	if err != nil {
		return err
	}

	path := filepath.Join(dir, "config.json")

	conf, err := json.Marshal(config)

	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, conf, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func getDefault() Config {
	return Config{
		StorageType: "local",
		RemoteUrl:   "",
	}
}
