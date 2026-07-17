package config

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"

	"aiub-companion/internal/event"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func init() {
	// Register a custom event whose associated data type is Config.
	application.RegisterEvent[Config](event.EventConfigChanged)
}

type Service struct {
	config atomic.Pointer[Config]
	path   string
	mu     sync.Mutex
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ServiceStartup(ctx context.Context, _ application.ServiceOptions) error {
	path, err := configPath()
	if err != nil {
		return fmt.Errorf("determine config path: %w", err)
	}

	config, err := load(path)
	if err != nil {
		slog.Error("Failed to load config on startup, using defaults", "error", err)
	}

	s.path = path
	s.config.Store(config)

	application.Get().Event.On(event.EventConfigChanged, onConfigChange)

	return nil
}

func onConfigChange(e *application.CustomEvent) {
	cfg, ok := e.Data.(Config)
	if !ok {
		return
	}

	app := application.Get()

	enabled, err := app.Autostart.IsEnabled()
	if err != nil {
		if !errors.Is(err, application.ErrAutostartNotSupported) {
			slog.Error("Failed to check autostart status", "error", err)
		}
	}

	if cfg.Launch.AutoStart != enabled {
		var err error

		if cfg.Launch.AutoStart {
			err = app.Autostart.Enable()
		} else {
			err = app.Autostart.Disable()
		}

		if err != nil {
			slog.Error("Failed to toggle autostart", "error", err)
			return
		}

		slog.Info("Autostart status", "enabled", cfg.Launch.AutoStart)
	}
}

func (s *Service) GetConfig() *Config {
	return s.config.Load()
}

func (s *Service) SaveConfig(config *Config) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	next := *config
	s.config.Store(&next)

	// Emit event to notify listeners of the change
	application.Get().Event.Emit(event.EventConfigChanged, next)

	// Persist to disk
	return save(s.path, &next)
}

func (s *Service) ResetConfig() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	defaults := defaultConfig()
	s.config.Store(defaults)

	// Emit event to notify listeners of the change
	application.Get().Event.Emit(event.EventConfigChanged, *defaults)

	// Persist to disk
	return save(s.path, defaults)
}
