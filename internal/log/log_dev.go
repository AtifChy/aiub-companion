//go:build !production

// Package log provides logging functionality for the AIUB Companion application.
package log

import (
	"errors"

	"aiub-companion/internal/log/logger"
)

var ErrNotImplemented = errors.New("not implemented in development mode")

func SetupLogger() (*logger.Logger, error) {
	return logger.NewLogger("")
}

func logPath() (string, error) {
	return "", ErrNotImplemented
}
