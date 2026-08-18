//go:build production

package log

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"time"

	"aiub-companion/internal/meta"

	"github.com/lmittmann/tint"
)

const maxLogFiles = 5

// newHandler creates a new slog.Handler that writes logs to both the console and a log file.
func newHandler(level *slog.LevelVar) (slog.Handler, io.Closer, error) {
	path, err := logPath()
	if err != nil {
		return nil, nil, err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, nil, err
	}

	// Best effort to prune old logs, ignoring any errors.
	pruneOldLogs()

	consoleHandler := tint.NewTextHandler(os.Stdout, &tint.Options{
		Level:     level,
		AddSource: false,
	})
	fileHandler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level:     level,
		AddSource: false,
	})

	return &multiHandler{handlers: []slog.Handler{consoleHandler, fileHandler}}, file, nil
}

type multiHandler struct {
	handlers []slog.Handler
}

func (h *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h *multiHandler) Handle(ctx context.Context, record slog.Record) error {
	var errs []error
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, record.Level) {
			if err := handler.Handle(ctx, record.Clone()); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errors.Join(errs...)
}

func (h *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithAttrs(attrs)
	}
	return &multiHandler{handlers: handlers}
}

func (h *multiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithGroup(name)
	}
	return &multiHandler{handlers: handlers}
}

func pruneOldLogs() {
	dir, err := logDir()
	if err != nil {
		return
	}

	pattern := filepath.Join(dir, meta.AppName+"-*.log")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) <= maxLogFiles {
		return
	}

	// Sort lexicographically, timestamps in filenames ensure that the most recent logs are last.
	slices.Sort(matches)

	// Remove old log files except the most recent maxLogFiles
	for _, old := range matches[:len(matches)-maxLogFiles] {
		_ = os.Remove(old)
	}
}

// logPath returns the full path to the log file.
func logPath() (string, error) {
	dir, err := logDir()
	if err != nil {
		return "", err
	}
	timestamp := time.Now().Format("2006-01-02T15-04-05")
	filename := fmt.Sprintf("%s-%s.log", meta.AppName, timestamp)
	return filepath.Join(dir, filename), nil
}
