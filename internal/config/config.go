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
)

type Config struct {
	// Appearance
	Theme string `json:"theme" jsonschema:"enum=light,enum=dark,enum=system"`

	// Advanced
	LogLevel string `json:"log_level" jsonschema:"enum=DEBUG,enum=INFO,enum=WARN,enum=ERROR"`

	// Sync
	Sync synchronization `json:"sync"`

	// Launch
	Launch launch `json:"launch"`

	// Notifications
	Notifications notification `json:"notifications"`
}

type notification struct {
	Enabled bool `json:"enabled"`
}

type launch struct {
	StartMinimized bool `json:"start_minimized"`
	CloseToTray    bool `json:"close_to_tray"`
	KeepAlive      bool `json:"keep_alive"`
	RestoreWindow  bool `json:"restore_window"`
	SidebarOpen    bool `json:"sidebar_open"`
}

type synchronization struct {
	IntervalMinutes int  `json:"interval_minutes"`
	FetchCount      int  `json:"fetch_count"`
	OnStartup       bool `json:"on_startup"`
}

func defaultConfig() *Config {
	return &Config{
		Notifications: notification{
			Enabled: true,
		},
		Launch: launch{
			StartMinimized: false,
			CloseToTray:    true,
			KeepAlive:      false,
			RestoreWindow:  false,
			SidebarOpen:    true,
		},
		Sync: synchronization{
			IntervalMinutes: 30,
			FetchCount:      20,
			OnStartup:       false,
		},
		Theme:    "system",
		LogLevel: "INFO",
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
		return cfg, fmt.Errorf("failed to read config file: %w", err)
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
		return cfg, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}

func save(path string, cfg *Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	tmp, err := os.CreateTemp(filepath.Dir(path), "config-*.json")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer func() {
		_ = os.Remove(tmp.Name())
	}()

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	if err := tmp.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	// atomic rename to avoid partial writes
	if err := os.Rename(tmp.Name(), path); err != nil {
		return fmt.Errorf("failed to remove temp file: %w", err)
	}

	return nil
}

func configPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config dir: %w", err)
	}
	return filepath.Join(configDir, AppName, "config.json"), nil
}
