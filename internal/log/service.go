package log

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"

	"aiub-companion/internal/config"
	"aiub-companion/internal/event"
	"aiub-companion/internal/meta"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type Service struct {
	logger *Logger
	config *config.Service
}

func NewService(logger *Logger, config *config.Service) *Service {
	return &Service{
		logger: logger,
		config: config,
	}
}

func (s *Service) ServiceStartup(ctx context.Context, _ application.ServiceOptions) error {
	application.Get().Event.On(event.EventConfigChanged, s.onConfigChanged)

	if cfg := s.config.GetConfig(); cfg != nil {
		s.logger.SetLevel(cfg.Logging.Level)
	} else {
		slog.Warn("No config found during log service startup")
	}

	return nil
}

func (s *Service) onConfigChanged(ev *application.CustomEvent) {
	cfg, ok := ev.Data.(config.Config)
	if !ok {
		return
	}

	level := cfg.Logging.Level
	if level == s.logger.Level() {
		return
	}
	s.logger.SetLevel(level)

	s.logger.Info("Log level changed", "level", level.String())
}

func (s *Service) Debug(message string) {
	s.logger.Debug("[frontend] " + message)
}

func (s *Service) Info(message string) {
	s.logger.Info("[frontend] " + message)
}

func (s *Service) Warn(message string) {
	s.logger.Warn("[frontend] " + message)
}

func (s *Service) Error(message string) {
	s.logger.Error("[frontend] " + message)
}

func (s *Service) OpenLogFile() error {
	dir, err := logDir()
	if err != nil {
		return err
	}

	pattern := filepath.Join(dir, meta.AppName+"-*.log")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("glob log files: %w", err)
	}

	if len(matches) == 0 {
		return errors.New("no log files found")
	}

	latestLogFile := matches[len(matches)-1]

	if err := application.Get().Browser.OpenFile(latestLogFile); err != nil {
		return fmt.Errorf("open log file: %w", err)
	}

	return nil
}
