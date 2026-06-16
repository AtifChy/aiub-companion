package settings

import (
	"log/slog"
	"time"
)

type Service struct {
	logLevel    *slog.LevelVar
	onSyncReset func(d time.Duration)
}

func NewService(logLevel *slog.LevelVar, onSyncReset func(d time.Duration)) *Service {
	return &Service{
		logLevel:    logLevel,
		onSyncReset: onSyncReset,
	}
}

func (s *Service) GetSettings() (*Settings, error) {
	return Load()
}

func (s *Service) SaveSettings(cfg *Settings) error {
	if s.logLevel != nil {
		if level, err := ParseLogLevel(cfg.LogLevel); err == nil {
			s.logLevel.Set(level)
		}
	}

	if s.onSyncReset != nil {
		s.onSyncReset(time.Duration(cfg.Sync.IntervalMinutes) * time.Minute)
	}

	return Save(cfg)
}
