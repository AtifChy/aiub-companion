// Package logger provides a simple logging interface.
package logger

import (
	"fmt"
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger
	Level *slog.LevelVar
	file  *os.File
}

func NewLogger(filepath string) (*Logger, error) {
	var file *os.File
	if filepath != "" {
		var err error
		file, err = os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			return nil, fmt.Errorf("open log file: %w", err)
		}
	}

	level := new(slog.LevelVar)
	level.Set(getDefaultLevel())

	handler := getHandler(file, level)

	return &Logger{
		Logger: slog.New(handler),
		Level:  level,
		file:   file,
	}, nil
}

func (l *Logger) L() *slog.Logger {
	return l.Logger
}

func (l *Logger) SetLevel(level slog.Level) {
	l.Level.Set(level)
}

func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}
