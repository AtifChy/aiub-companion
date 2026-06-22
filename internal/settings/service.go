package settings

import (
	"log/slog"
	"sync"
	"sync/atomic"
)

type Service struct {
	writeMu   sync.Mutex
	config    atomic.Pointer[Settings]
	listeners []func(*Settings)
}

func NewService() *Service {
	config, err := load()
	if err != nil {
		slog.Error("Failed to load settings, using defaults", "error", err)
	}
	s := &Service{}
	s.config.Store(config)
	return s
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

	// Trigger side effects
	for _, listener := range s.listeners {
		listener(&next)
	}

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
func (s *Service) OnChange(listener func(*Settings)) {
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	s.listeners = append(s.listeners, listener)
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
