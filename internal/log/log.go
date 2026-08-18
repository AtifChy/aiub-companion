// Package log provides logging functionality for the AIUB Companion application.
package log

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"aiub-companion/internal/meta"
)

type Logger struct {
	*slog.Logger
	Level  *slog.LevelVar
	closer io.Closer
}

func NewLogger() (*Logger, error) {
	level := new(slog.LevelVar)
	level.Set(slog.LevelInfo)

	handler, closer, err := newHandler(level)
	if err != nil {
		return nil, fmt.Errorf("create log handler: %w", err)
	}

	return &Logger{
		Logger: slog.New(handler),
		Level:  level,
		closer: closer,
	}, nil
}

// SetLevel sets the logging level at runtime.
func (l *Logger) SetLevel(level slog.Level) {
	l.Level.Set(level)
}

// Close releases underlying resources used by the logger, if any.
func (l *Logger) Close() error {
	if l.closer != nil {
		return l.closer.Close()
	}
	return nil
}

// logDir returns the directory for log files, creating it if it doesn't exist.
func logDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(configDir, meta.AppName, "logs")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}
