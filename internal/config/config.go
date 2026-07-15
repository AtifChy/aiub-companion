//go:generate go run gen.go

// Package config holds the application configuration.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"aiub-companion/internal/meta"
)

type Config struct {
	Appearance    appearance   `json:"appearance"`
	Updates       updates      `json:"updates"`
	Logging       logging      `json:"logging"`
	Sync          sync_        `json:"sync"`
	Launch        launch       `json:"launch"`
	Notifications notification `json:"notifications"`
}

type appearance struct {
	Theme string `json:"theme" jsonschema:"enum=light,enum=dark,enum=system"`
}

type notification struct {
	Enabled bool `json:"enabled"`
}

type launch struct {
	AutoStart      bool `json:"auto_start"`
	StartMinimized bool `json:"start_minimized"`
	CloseToTray    bool `json:"close_to_tray"`
	KeepAlive      bool `json:"keep_alive"`
	RestoreWindow  bool `json:"restore_window"`
	SidebarOpen    bool `json:"sidebar_open"`
}

type sync_ struct {
	Interval   int  `json:"interval"`
	FetchCount int  `json:"fetch_count"`
	OnStartup  bool `json:"on_startup"`
}

type updates struct {
	Interval string `json:"interval" jsonschema:"enum=never,enum=daily,enum=weekly,enum=monthly"`
}

type logging struct {
	Level string `json:"level" jsonschema:"enum=DEBUG,enum=INFO,enum=WARN,enum=ERROR"`
}

func defaultConfig() *Config {
	return &Config{
		Appearance: appearance{
			Theme: "system",
		},
		Notifications: notification{
			Enabled: true,
		},
		Launch: launch{
			AutoStart:      false,
			StartMinimized: false,
			CloseToTray:    true,
			KeepAlive:      false,
			RestoreWindow:  false,
			SidebarOpen:    true,
		},
		Sync: sync_{
			Interval:   30,
			FetchCount: 20,
			OnStartup:  false,
		},
		Updates: updates{
			Interval: "weekly",
		},
		Logging: logging{
			Level: "INFO",
		},
	}
}

func ParseLogLevel(s string) (slog.Level, error) {
	var level slog.Level
	return level, level.UnmarshalText([]byte(s))
}

func load(path string) (*Config, error) {
	cfg := defaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("read config file: %w", err)
	}

	if err := validate(data); err != nil {
		return cfg, err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		if syntaxErr, ok := errors.AsType[*json.SyntaxError](err); ok {
			return cfg, fmt.Errorf("invalid JSON syntax at offset %d: %w", syntaxErr.Offset, err)
		}
		if typeErr, ok := errors.AsType[*json.UnmarshalTypeError](err); ok {
			return cfg, fmt.Errorf("invalid JSON type for field %q (expected %s): %w", typeErr.Field, typeErr.Type, err)
		}
		return cfg, fmt.Errorf("unmarshal config: %w", err)
	}

	return cfg, nil
}

func save(path string, cfg *Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	tmp, err := os.CreateTemp(filepath.Dir(path), "config-*.json")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	defer func() { _ = os.Remove(tmp.Name()) }()

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("write temp file: %w", err)
	}

	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temp file: %w", err)
	}

	// atomic rename to avoid partial writes
	if err := os.Rename(tmp.Name(), path); err != nil {
		return fmt.Errorf("rename temp file: %w", err)
	}

	return nil
}

func configPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("get user config dir: %w", err)
	}
	return filepath.Join(configDir, meta.AppName, "config.json"), nil
}
