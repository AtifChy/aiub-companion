package config

import (
	"context"
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
	config  atomic.Pointer[Config]
	writeMu sync.Mutex
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ServiceStartup(ctx context.Context, _ application.ServiceOptions) error {
	config, err := load()
	if err != nil {
		slog.Error("Failed to load config on startup, using defaults", "error", err)
	}
	s.config.Store(config)
	return nil
}

func (s *Service) GetConfig() *Config {
	return s.config.Load()
}

func (s *Service) SaveConfig(config *Config) error {
	s.writeMu.Lock()
	defer s.writeMu.Unlock()

	current := s.config.Load()
	next := *config
	next.Window = current.Window

	s.config.Store(&next)

	// Emit event to notify listeners of the change
	application.Get().Event.Emit(event.EventConfigChanged, next)

	// Persist to disk
	return save(&next)
}

func (s *Service) ResetConfig() error {
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	defaults := defaultConfig()
	s.config.Store(defaults)
	return save(defaults)
}

//wails:ignore
func (s *Service) UpdateWindowState(state *WindowState) error {
	s.writeMu.Lock()
	defer s.writeMu.Unlock()

	current := s.config.Load()
	next := *current
	if state != nil {
		next.Window = *state
	}

	s.config.Store(&next)
	return save(&next)
}
