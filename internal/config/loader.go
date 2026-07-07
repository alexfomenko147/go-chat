package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func Load(path string) (*Config, error) {
	cfg := Default()

	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return cfg, nil
		}
		path = filepath.Join(home, ".config", "go-chat", "config.yaml")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("read config: %w", err)
	}

	ext := filepath.Ext(path)
	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("parse yaml: %w", err)
		}
	case ".json":
		if err := json.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("parse json: %w", err)
		}
	case ".toml":
		return nil, fmt.Errorf("toml support requires burntsushi/toml: implement with go-toml")
	default:
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("parse config (default yaml): %w", err)
		}
	}

	return cfg, nil
}

func Save(path string, cfg *Config) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}
