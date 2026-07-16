//go:build production

// Package log provides logging functionality for the application.
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
	path, err := logPath()
	if err != nil {
		return nil, err
	}

	log, err := logger.NewLogger(path)
	if err != nil {
		return nil, fmt.Errorf("create logger: %w", err)
	}

	return log, nil
}

func logPath() (string, error) {
	logDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	logDir = path.Join(logDir, meta.AppName, "logs")
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return "", fmt.Errorf("create log directory: %w", err)
	}

	return filepath.Join(logDir, meta.AppName+".log"), nil
}
