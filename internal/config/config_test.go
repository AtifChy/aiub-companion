package config

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	t.Run("Non-existent file returns default config", func(t *testing.T) {
		tempDir := t.TempDir()
		nonExistentPath := filepath.Join(tempDir, "does-not-exist.json")

		cfg, err := load(nonExistentPath)
		if err != nil {
			t.Fatalf("expected no error for non-existent file, got %v", err)
		}
		expected := defaultConfig()
		if *cfg != *expected {
			t.Errorf("expected default config %+v, got %+v", expected, cfg)
		}
	})

	t.Run("Valid config file loads successfully", func(t *testing.T) {
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "config.json")
		validJSON := `{
			"appearance": { "color": "dark", "theme": "dracula" },
			"updates": { "interval": "weekly" },
			"logging": { "level": "DEBUG" },
			"sync": { "interval": 45, "fetch_count": 50, "on_startup": false },
			"launch": {
				"auto_start": true,
				"start_minimized": true,
				"close_to_tray": false,
				"keep_alive": true,
				"restore_window": false,
				"sidebar_open": true
			},
			"notifications": { "enabled": false }
		}`
		if err := os.WriteFile(configPath, []byte(validJSON), 0o644); err != nil {
			t.Fatalf("failed to write config file: %v", err)
		}

		cfg, err := load(configPath)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if cfg.Appearance.Color != "dark" || cfg.Appearance.Theme != "dracula" {
			t.Errorf("unexpected appearance: %+v", cfg.Appearance)
		}
		if cfg.Updates.Interval != "weekly" {
			t.Errorf("unexpected updates interval: %s", cfg.Updates.Interval)
		}
		if cfg.Logging.Level != slog.LevelDebug {
			t.Errorf("unexpected logging level: %s", cfg.Logging.Level)
		}
		if cfg.Sync.Interval != 45 || cfg.Sync.FetchCount != 50 || cfg.Sync.OnStartup != false {
			t.Errorf("unexpected sync settings: %+v", cfg.Sync)
		}
		if !cfg.Launch.AutoStart || !cfg.Launch.StartMinimized || cfg.Launch.CloseToTray || !cfg.Launch.KeepAlive || cfg.Launch.RestoreWindow || !cfg.Launch.SidebarOpen {
			t.Errorf("unexpected launch settings: %+v", cfg.Launch)
		}
		if cfg.Notifications.Enabled {
			t.Errorf("expected notifications to be disabled")
		}
	})

	t.Run("Schema violation preserves valid settings and loads without fatal failure", func(t *testing.T) {
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "config.json")
		schemaViolationJSON := `{
			"appearance": { "color": "invalid_color", "theme": "dracula" },
			"updates": { "interval": "weekly" },
			"logging": { "level": "DEBUG" },
			"sync": { "interval": 15, "fetch_count": 10, "on_startup": true },
			"launch": {
				"auto_start": true,
				"start_minimized": true,
				"close_to_tray": false,
				"keep_alive": true,
				"restore_window": false,
				"sidebar_open": true
			},
			"notifications": { "enabled": false }
		}`
		if err := os.WriteFile(configPath, []byte(schemaViolationJSON), 0o644); err != nil {
			t.Fatalf("failed to write config file: %v", err)
		}

		cfg, err := load(configPath)
		if err != nil {
			t.Fatalf("expected load to succeed with warning, got error: %v", err)
		}

		// Valid fields preserved
		if cfg.Appearance.Theme != "dracula" {
			t.Errorf("expected theme 'dracula' to be preserved, got %s", cfg.Appearance.Theme)
		}
		if !cfg.Launch.AutoStart || !cfg.Launch.StartMinimized || cfg.Launch.CloseToTray || !cfg.Launch.KeepAlive || cfg.Launch.RestoreWindow || !cfg.Launch.SidebarOpen {
			t.Errorf("expected launch settings to be preserved, got %+v", cfg.Launch)
		}
		if cfg.Notifications.Enabled {
			t.Errorf("expected notifications enabled=false to be preserved")
		}
	})
}

func TestSaveAndRoundtrip(t *testing.T) {
	tempDir := t.TempDir()
	nestedPath := filepath.Join(tempDir, "nested", "dir", "config.json")

	cfg := &Config{
		Appearance: appearance{
			Color: "dark",
			Theme: "nord",
		},
		Updates: updates{
			Interval: "monthly",
		},
		Logging: logging{
			Level: slog.LevelInfo,
		},
		Sync: syn_{
			Interval:   60,
			FetchCount: 100,
			OnStartup:  false,
		},
		Launch: launch{
			AutoStart:      true,
			StartMinimized: true,
			CloseToTray:    false,
			KeepAlive:      true,
			RestoreWindow:  false,
			SidebarOpen:    true,
		},
		Notifications: notification{
			Enabled: false,
		},
	}

	if err := save(nestedPath, cfg); err != nil {
		t.Fatalf("save() failed: %v", err)
	}

	// Verify file was written and is formatted with indentation
	raw, err := os.ReadFile(nestedPath)
	if err != nil {
		t.Fatalf("failed to read saved file: %v", err)
	}
	if !strings.Contains(string(raw), "  \"appearance\":") {
		t.Errorf("expected indented JSON output, got:\n%s", string(raw))
	}

	// Load and verify equality
	loaded, err := load(nestedPath)
	if err != nil {
		t.Fatalf("load() failed on saved config: %v", err)
	}

	if *loaded != *cfg {
		t.Errorf("roundtrip mismatch:\nexpected: %+v\ngot:      %+v", *cfg, *loaded)
	}
}

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"DEBUG", false},
		{"INFO", false},
		{"WARN", false},
		{"ERROR", false},
		{"TRACE", true},
		{"invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			_, err := ParseLogLevel(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLogLevel(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}
