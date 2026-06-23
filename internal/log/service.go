package log

import "log/slog"

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Debug(message string) {
	slog.Debug("[frontend] " + message)
}

func (s *Service) Info(message string) {
	slog.Info("[frontend] " + message)
}

func (s *Service) Warn(message string) {
	slog.Warn("[frontend] " + message)
}

func (s *Service) Error(message string) {
	slog.Error("[frontend] " + message)
}
