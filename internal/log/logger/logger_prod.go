//go:build production

package logger

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func getDefaultLevel() slog.Level {
	return slog.LevelInfo
}

func getHandler(file *os.File, level *slog.LevelVar) slog.Handler {
	console := tint.NewHandler(os.Stdout, &tint.Options{
		Level:     level,
		AddSource: false,
	})
	if file == nil {
		return console
	}
	fileHandler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level:     level,
		AddSource: false,
	})
	return &multiHandler{handlers: []slog.Handler{console, fileHandler}}
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
