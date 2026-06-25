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

	// Window state
	Window WindowState `json:"window"`
}

type notification struct {
	Enabled bool `json:"enabled"`
}

type launch struct {
	StartMinimized bool `json:"start_minimized"`
	CloseToTray    bool `json:"close_to_tray"`
	RestoreWindow  bool `json:"restore_window"`
	SidebarOpen    bool `json:"sidebar_open"`
}

type synchronization struct {
	IntervalMinutes int  `json:"interval_minutes"`
	FetchCount      int  `json:"fetch_count"`
	OnStartup       bool `json:"on_startup"`
}

type WindowState struct {
	Width     int  `json:"width"`
	Height    int  `json:"height"`
	X         int  `json:"x"`
	Y         int  `json:"y"`
	Maximized bool `json:"maximized"`
}

func defaultConfig() *Config {
	return &Config{
		Notifications: notification{
			Enabled: true,
		},
		Launch: launch{
			StartMinimized: false,
			CloseToTray:    true,
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
		Window: WindowState{
			Width:     1024,
			Height:    768,
			X:         -1,
			Y:         -1,
			Maximized: false,
		},
	}
}

func ParseLogLevel(s string) (slog.Level, error) {
	var level slog.Level
	return level, level.UnmarshalText([]byte(s))
}

func load() (*Config, error) {
	cfg := defaultConfig()

	path, err := configPath()
	if err != nil {
		return cfg, fmt.Errorf("failed to get settings path: %w", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("failed to read settings file: %w", err)
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
		return cfg, fmt.Errorf("failed to unmarshal settings: %w", err)
	}

	return cfg, nil
}

func save(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return fmt.Errorf("failed to get settings path: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("failed to create settings directory: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	tmp, err := os.CreateTemp(filepath.Dir(path), "settings-*.json")
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
