//go:build !production

package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func getDefaultLevel() slog.Level {
	return slog.LevelInfo
}

func getHandler(_ *os.File, level *slog.LevelVar) slog.Handler {
	return tint.NewHandler(os.Stdout, &tint.Options{
		Level:     level,
		AddSource: true,
	})
}
