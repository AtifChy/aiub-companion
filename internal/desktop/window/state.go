package window

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"aiub-companion/internal/config"
)

// WindowState holds the last known state of a window.
type WindowState struct {
	Width     int  `json:"width"`
	Height    int  `json:"height"`
	X         int  `json:"x"`
	Y         int  `json:"y"`
	Maximized bool `json:"maximized"`
}

type AppState struct {
	Window map[string]WindowState `json:"window"`
}

func defaultState() WindowState {
	return WindowState{
		Width:     1024,
		Height:    768,
		X:         -1,
		Y:         -1,
		Maximized: false,
	}
}

func loadState(name string) WindowState {
	path, err := statePath()
	if err != nil {
		return defaultState()
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return defaultState()
	}
	var appState AppState
	if err := json.Unmarshal(data, &appState); err != nil {
		return defaultState()
	}
	if state, ok := appState.Window[name]; ok {
		return state
	}
	return defaultState()
}

func saveState(name string, state WindowState) error {
	path, err := statePath()
	if err != nil {
		return fmt.Errorf("failed to determine state path: %w", err)
	}

	var appState AppState
	data, err := os.ReadFile(path)
	if err == nil {
		_ = json.Unmarshal(data, &appState)
	}
	if appState.Window == nil {
		appState.Window = make(map[string]WindowState)
	}

	appState.Window[name] = state

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("failed to create state directory: %w", err)
	}

	data, err = json.MarshalIndent(appState, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}

func statePath() (string, error) {
	path, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(path, config.AppName, "state.json"), nil
}
