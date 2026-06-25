package log

import (
	"context"

	"aiub-companion/internal/config"
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
	application.Get().Event.On(config.EventConfigChanged, func(event *application.CustomEvent) {
		if cfg, ok := event.Data.(config.Config); ok {
			if level, err := config.ParseLogLevel(cfg.LogLevel); err == nil {
				s.logger.SetLevel(level)
			} else {
				s.logger.L().Warn("invalid log level in config", "error", err)
			}
		}
	})

	return nil
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
