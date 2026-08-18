//go:build !production

package log

import (
	"io"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

// newHandler creates a new slog.Handler that writes logs to the console with colorized output.
func newHandler(level *slog.LevelVar) (slog.Handler, io.Closer, error) {
	return tint.NewTextHandler(os.Stdout, &tint.Options{
		Level:     level,
		AddSource: true,
	}), nil, nil
}
