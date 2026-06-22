//go:build production

// Package log provides logging functionality for the AIUB Companion metalication.
package log

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"aiub-companion/internal/log/logger"
	"aiub-companion/internal/meta"
)

func SetupLogger() (*logger.Logger, error) {
	logDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	logDir = path.Join(logDir, meta.Name, "logs")
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	filename := filepath.Join(logDir, meta.Name+".log")

	log, err := logger.NewLogger(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	return log, nil
}
