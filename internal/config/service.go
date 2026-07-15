package config

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"

	"aiub-companion/internal/autostart"
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

	if enabled, err := autostart.IsEnabled(); err != nil {
		slog.Error("Failed to check autostart status", "error", err)
	} else if config.Launch.AutoStart != enabled {
		if err := autostart.Set(config.Launch.AutoStart); err != nil {
			slog.Error("Failed to set autostart", "error", err)
		}
	}

	s.path = path
	s.config.Store(config)

	return nil
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
