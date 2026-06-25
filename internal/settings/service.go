package settings

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"

	"github.com/wailsapp/wails/v3/pkg/application"
)

const EventSettingsChanged = "settings:changed"

func init() {
	// Register a custom event whose associated data type is *Settings.
	application.RegisterEvent[Settings](EventSettingsChanged)
}

type Service struct {
	writeMu sync.Mutex
	config  atomic.Pointer[Settings]
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ServiceStartup(ctx context.Context, _ application.ServiceOptions) error {
	config, err := load()
	if err != nil {
		slog.Error("Failed to load settings on startup, using defaults", "error", err)
	}
	s.config.Store(config)
	return nil
}

func (s *Service) GetSettings() *Settings {
	return s.config.Load()
}

func (s *Service) SaveSettings(config *Settings) error {
	s.writeMu.Lock()
	defer s.writeMu.Unlock()

	current := s.config.Load()
	next := *config
	next.Window = current.Window

	s.config.Store(&next)

	// Emit event to notify listeners of the change
	application.Get().Event.Emit(EventSettingsChanged, next)

	// Persist to disk
	return save(&next)
}

func (s *Service) ResetSettings() error {
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	defaults := defaultSettings()
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
