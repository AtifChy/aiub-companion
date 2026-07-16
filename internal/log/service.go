package log

import (
	"context"
	"fmt"

	"aiub-companion/internal/config"
	"aiub-companion/internal/event"
	"aiub-companion/internal/log/logger"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type Service struct {
	logger *logger.Logger
}

func NewService(logger *logger.Logger) *Service {
	return &Service{logger: logger}
}

func (s *Service) ServiceStartup(ctx context.Context, _ application.ServerOptions) error {
	application.Get().Event.On(event.EventConfigChanged, s.onConfigChanged)
	return nil
}

func (s *Service) onConfigChanged(ev *application.CustomEvent) {
	cfg, ok := ev.Data.(config.Config)
	if !ok {
		return
	}

	level, err := config.ParseLogLevel(cfg.Logging.Level)
	if err != nil {
		s.logger.L().Warn("Failed to parse log level", "error", err)
		return
	}

	s.logger.SetLevel(level)
	s.logger.L().Info("Log level changed", "level", level.String())
}

func (s *Service) Debug(message string) {
	s.logger.L().Debug("[frontend] " + message)
}

func (s *Service) Info(message string) {
	s.logger.L().Info("[frontend] " + message)
}

func (s *Service) Warn(message string) {
	s.logger.L().Warn("[frontend] " + message)
}

func (s *Service) Error(message string) {
	s.logger.L().Error("[frontend] " + message)
}

func (s *Service) OpenLogFile() error {
	path, err := logPath()
	if err != nil {
		return err
	}

	if err := application.Get().Browser.OpenFile(path); err != nil {
		return fmt.Errorf("open log file: %w", err)
	}

	return nil
}
