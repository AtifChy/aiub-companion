package updater

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"aiub-companion/internal/meta"
)

type state struct {
	LastCheckedAt time.Time `json:"last_checked_at"`
}

func loadState() (state, error) {
	path, err := statePath()
	if err != nil {
		return state{}, fmt.Errorf("determine state path: %w", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return state{}, nil // first run
		}
		return state{}, fmt.Errorf("read updater state: %w", err)
	}

	var s state
	if err := json.Unmarshal(data, &s); err != nil {
		return state{}, fmt.Errorf("unmarshal updater state: %w", err)
	}
	return s, nil
}

func saveState(s state) error {
	path, err := statePath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create state dir: %w", err)
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal updater state: %w", err)
	}

	return os.WriteFile(path, data, 0o644)
}

func statePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("get user config dir: %w", err)
	}
	return filepath.Join(configDir, meta.AppName, "state", "updater.json"), nil
}
