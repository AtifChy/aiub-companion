//go:build production

// Package log provides logging functionality for the application.
package log

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"aiub-companion/internal/config"
	"aiub-companion/internal/log/logger"
)

func SetupLogger() (*logger.Logger, error) {
	path, err := logPath()
	if err != nil {
		return nil, err
	}

	log, err := logger.NewLogger(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	return log, nil
}

func logPath() (string, error) {
	logDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	logDir = path.Join(logDir, config.AppName, "logs")
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create log directory: %w", err)
	}

	return filepath.Join(logDir, config.AppName+".log"), nil
}
