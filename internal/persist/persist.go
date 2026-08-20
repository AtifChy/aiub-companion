// Package persist provides functions to load and save application state to a JSON file.
package persist

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"aiub-companion/internal/meta"
)

func Path(name string) (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("get user config dir: %w", err)
	}
	return filepath.Join(configDir, meta.AppName, "state", name), nil
}

func Load[T any](path string) (T, error) {
	var v T

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return v, nil // first run
		}
		return v, fmt.Errorf("read state: %w", err)
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return v, fmt.Errorf("unmarshal state: %w", err)
	}

	return v, nil
}

func Save[T any](path string, v T) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create state dir: %w", err)
	}

	data, err := json.Marshal(v, jsontext.WithIndent("  "))
	if err != nil {
		return fmt.Errorf("marshal state: %w", err)
	}

	tmp, err := os.CreateTemp(dir, filepath.Base(path)+"-*.tmp")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	defer func() { _ = os.Remove(tmp.Name()) }()

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temp file: %w", err)
	}

	// atomic rename to avoid partial writes
	if err := os.Rename(tmp.Name(), path); err != nil {
		return fmt.Errorf("rename temp file: %w", err)
	}

	return nil
}
