package window

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"aiub-companion/internal/config"
)

// windowState holds the last known state of a window.
type windowState struct {
	Width     int  `json:"width"`
	Height    int  `json:"height"`
	X         int  `json:"x"`
	Y         int  `json:"y"`
	Maximized bool `json:"maximized"`
}

type AppState struct {
	Window map[string]windowState `json:"window"`
}

func defaultState() windowState {
	return windowState{
		Width:     1024,
		Height:    768,
		X:         -1,
		Y:         -1,
		Maximized: false,
	}
}

func loadState(name string) (windowState, error) {
	path, err := statePath()
	if err != nil {
		return windowState{}, fmt.Errorf("determine state path: %w", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return defaultState(), nil // first run
		}
		return windowState{}, fmt.Errorf("read window state: %w", err)
	}

	var appState AppState
	if err := json.Unmarshal(data, &appState); err != nil {
		return windowState{}, fmt.Errorf("unmarshal state: %w", err)
	}

	if state, ok := appState.Window[name]; ok {
		return state, nil
	}

	return windowState{}, nil
}

func saveState(name string, state windowState) error {
	path, err := statePath()
	if err != nil {
		return fmt.Errorf("determine state path: %w", err)
	}

	var appState AppState
	data, err := os.ReadFile(path)
	if err == nil {
		_ = json.Unmarshal(data, &appState)
	}
	if appState.Window == nil {
		appState.Window = make(map[string]windowState)
	}

	appState.Window[name] = state

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create state directory: %w", err)
	}

	data, err = json.MarshalIndent(appState, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal state: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}

func statePath() (string, error) {
	path, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(path, config.AppName, "state", "window.json"), nil
}
