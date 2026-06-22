//go:generate go run gen.go

// Package settings holds the application settings.
package settings

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"aiub-companion/internal/meta"
)

type Settings struct {
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
}

type synchronization struct {
	IntervalMinutes int  `json:"interval_minutes"`
	FetchCount      int  `json:"fetch_count"`
	OnStartup       bool `json:"on_startup"`
}

type WindowState struct {
	Enabled   bool `json:"enabled"`
	Width     int  `json:"width"`
	Height    int  `json:"height"`
	X         int  `json:"x"`
	Y         int  `json:"y"`
	Maximized bool `json:"maximized"`
}

func defaultSettings() *Settings {
	return &Settings{
		Notifications: notification{
			Enabled: true,
		},
		Launch: launch{
			StartMinimized: false,
			CloseToTray:    true,
		},
		Sync: synchronization{
			IntervalMinutes: 30,
			FetchCount:      20,
			OnStartup:       false,
		},
		Theme:    "system",
		LogLevel: "INFO",
		Window: WindowState{
			Enabled:   false,
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

func load() (*Settings, error) {
	s := defaultSettings()

	path, err := settingsPath()
	if err != nil {
		return s, fmt.Errorf("failed to get settings path: %w", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return s, nil
		}
		return s, fmt.Errorf("failed to read settings file: %w", err)
	}

	if err := validate(data); err != nil {
		return s, err
	}

	if err := json.Unmarshal(data, s); err != nil {
		if syntaxErr, ok := errors.AsType[*json.SyntaxError](err); ok {
			return s, fmt.Errorf("invalid JSON syntax at offset %d: %w", syntaxErr.Offset, err)
		}
		if typeErr, ok := errors.AsType[*json.UnmarshalTypeError](err); ok {
			return s, fmt.Errorf("invalid JSON type for field %q (expected %s): %w", typeErr.Field, typeErr.Type, err)
		}
		return s, fmt.Errorf("failed to unmarshal settings: %w", err)
	}

	return s, nil
}

func save(s *Settings) error {
	path, err := settingsPath()
	if err != nil {
		return fmt.Errorf("failed to get settings path: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("failed to create settings directory: %w", err)
	}

	data, err := json.MarshalIndent(s, "", "  ")
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

func settingsPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config dir: %w", err)
	}
	return filepath.Join(configDir, meta.Name, "settings.json"), nil
}
