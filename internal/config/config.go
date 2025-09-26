package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Profiles  []string `json:"profiles"`
	NtfyTopic string   `json:"ntfy_topic"`
}

func Load(configPath string) (Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
