package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateSchema(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
	}{
		{
			name: "Valid config",
			json: `{
				"appearance": { "theme": "default", "color": "system" },
				"updates": { "interval": "weekly" },
				"logging": { "level": "DEBUG" },
				"sync": {
					"interval": 15,
					"fetch_count": 10,
					"on_startup": true
				},
				"launch": {
					"auto_start": false,
					"start_minimized": true,
					"close_to_tray": false,
					"keep_alive": false,
					"restore_window": false,
					"sidebar_open": true
				},
				"notifications": {
					"enabled": true
				}
			}`,
			wantErr: false,
		},
		{
			name: "Invalid color value",
			json: `{
				"appearance": { "theme": "default", "color": "invalid-color-value" },
				"updates": { "interval": "weekly" },
				"logging": { "level": "DEBUG" },
				"sync": { "interval": 15, "fetch_count": 10, "on_startup": true },
				"launch": { "start_minimized": true, "close_to_tray": false, "keep_alive": false, "restore_window": false, "sidebar_open": true },
				"notifications": { "enabled": true }
			}`,
			wantErr: true,
		},
		{
			name: "Invalid log level",
			json: `{
				"appearance": { "theme": "default", "color": "system" },
				"updates": { "interval": "weekly" },
				"logging": { "level": "TRACE" },
				"sync": { "interval": 15, "fetch_count": 10, "on_startup": true },
				"launch": { "start_minimized": true, "close_to_tray": false, "keep_alive": false, "restore_window": false, "sidebar_open": true },
				"notifications": { "enabled": true }
			}`,
			wantErr: true,
		},
		{
			name: "Malformed JSON syntax",
			json: `{
				"appearance": { "theme": "default", "color": "system" },
				"updates": { "interval": "weekly" },
				"logging": { "level": "TRACE" },
				"sync": { "interval": 15, "fetch_count": 10, "on_startup": true },
				"launch": { "start_minimized": true, "close_to_tray": false, "keep_alive": false, "restore_window": false, "sidebar_open": true },
				"notifications": { "enabled": true }
			}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate([]byte(tt.json))
			if (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	t.Run("Non-existent file returns default config", func(t *testing.T) {
		tempDir := t.TempDir()
		nonExistentPath := filepath.Join(tempDir, "nonexistent.json")

		cfg, err := load(nonExistentPath)
		if err != nil {
			t.Fatalf("expected no error for non-existent file, got: %v", err)
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
		if cfg.Logging.Level != "DEBUG" {
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

	t.Run("Invalid JSON syntax returns syntactic error", func(t *testing.T) {
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "config.json")
		invalidJSON := `{
			"appearance": { "color": "dark",
		}`
		if err := os.WriteFile(configPath, []byte(invalidJSON), 0o644); err != nil {
			t.Fatalf("failed to write config file: %v", err)
		}

		_, err := load(configPath)
		if err == nil {
			t.Fatal("expected syntax error, got nil")
		}
		if !strings.Contains(err.Error(), "invalid JSON syntax") && !strings.Contains(err.Error(), "syntax") {
			t.Errorf("expected syntax error message, got %v", err)
		}
	})

	t.Run("Schema validation failure returns error", func(t *testing.T) {
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "config.json")
		invalidEnumJSON := `{
			"appearance": { "color": "invalid_color", "theme": "default" },
			"updates": { "interval": "weekly" },
			"logging": { "level": "DEBUG" },
			"sync": { "interval": 15, "fetch_count": 10, "on_startup": true },
			"launch": {
				"auto_start": false,
				"start_minimized": true,
				"close_to_tray": false,
				"keep_alive": false,
				"restore_window": false,
				"sidebar_open": true
			},
			"notifications": { "enabled": true }
		}`
		if err := os.WriteFile(configPath, []byte(invalidEnumJSON), 0o644); err != nil {
			t.Fatalf("failed to write config file: %v", err)
		}

		_, err := load(configPath)
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
		if !strings.Contains(err.Error(), "invalid config") {
			t.Errorf("expected schema validation error, got: %v", err)
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
			Level: "INFO",
		},
		Sync: sync_{
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
