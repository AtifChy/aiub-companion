//go:build !production

// Package log provides logging functionality for the AIUB Companion application.
package log

import (
	"aiub-companion/internal/log/logger"
)

func SetupLogger() (*logger.Logger, error) {
	return logger.NewLogger("")
}
