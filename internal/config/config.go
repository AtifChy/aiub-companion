//go:generate go run schema.go

// Package config holds the application configuration.
package config

import (
	"encoding"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"aiub-companion/internal/meta"
)

type Config struct {
	Appearance    appearance   `json:"appearance"`
	Updates       updates      `json:"updates"`
	Logging       logging      `json:"logging"`
	Sync          syn_         `json:"sync"`
	Launch        launch       `json:"launch"`
	Notifications notification `json:"notifications"`
}

func defaultConfig() *Config {
	return &Config{
		Appearance: appearance{
			Color: ColorSystem,
			Theme: "default",
		},
		Notifications: notification{
			Enabled: true,
		},
		Launch: launch{
			AutoStart:      false,
			StartMinimized: false,
			CloseToTray:    true,
			KeepAlive:      false,
			RestoreWindow:  true,
			SidebarOpen:    false,
		},
		Sync: syn_{
			Interval:   30,
			FetchCount: 20,
			OnStartup:  true,
		},
		Updates: updates{
			Interval: UpdateIntervalDaily,
		},
		Logging: logging{
			Level: slog.LevelWarn,
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

	loadFields(data, cfg)

	return cfg, nil
}

var (
	textUnmarshalerType = reflect.TypeFor[encoding.TextUnmarshaler]()
	jsonUnmarshalerType = reflect.TypeFor[json.Unmarshaler]()
)

func loadFields(raw jsontext.Value, dst any) {
	if len(raw) == 0 {
		return
	}

	var fields map[string]jsontext.Value
	if err := json.Unmarshal(raw, &fields); err != nil {
		slog.Warn("Invalid section", "type", reflect.TypeOf(dst).Elem().Name(), "error", err)
		return
	}

	rv := reflect.ValueOf(dst).Elem()
	rt := rv.Type()

	for sf := range rt.Fields() {
		tag, _, _ := strings.Cut(sf.Tag.Get("json"), ",")
		if tag == "" || tag == "-" {
			continue
		}
		rawField, ok := fields[tag]
		if !ok {
			continue
		}

		fv := rv.FieldByIndex(sf.Index)
		fvPtr := fv.Addr()
		hasOwnDecoder := fvPtr.Type().Implements(jsonUnmarshalerType) || fvPtr.Type().Implements(textUnmarshalerType)

		if fv.Kind() == reflect.Struct && !hasOwnDecoder {
			loadFields(rawField, fvPtr.Interface())
			continue
		}

		if err := json.Unmarshal(rawField, fvPtr.Interface()); err != nil {
			slog.Warn("Invalid field", "type", rt.Name(), "field", tag, "error", err)
		}
	}
}

func save(path string, cfg *Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	data, err := json.Marshal(cfg, jsontext.WithIndent("  "))
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
