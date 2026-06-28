package window

import (
	"testing"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func TestWindowOptionsAndGettersSetters(t *testing.T) {
	app := &application.App{}
	opts := WindowOptions{
		HideOnClose:   true,
		RestoreWindow: true,
	}
	opts.Name = "test-window"

	w := NewWindow(app, opts)

	if w.HideOnClose() != true {
		t.Errorf("expected HideOnClose to be true, got false")
	}

	w.SetHideOnClose(false)
	if w.HideOnClose() != false {
		t.Errorf("expected HideOnClose to be false, got true")
	}

	w.SetRestoreWindow(false)
}

func TestStateLoadingAndSaving(t *testing.T) {
	tempDir := t.TempDir()

	// Sandboxing environment variables for platform-independent config folder testing
	t.Setenv("APPDATA", tempDir)
	t.Setenv("XDG_CONFIG_HOME", tempDir)
	t.Setenv("HOME", tempDir)

	name := "main-window"
	state := WindowState{
		Width:     800,
		Height:    600,
		X:         100,
		Y:         150,
		Maximized: true,
	}

	// 1. Loading non-existent state should return defaultState
	loaded := loadState(name)
	if loaded.Width != 1024 || loaded.Height != 768 {
		t.Errorf("expected default size 1024x768, got %dx%d", loaded.Width, loaded.Height)
	}

	// 2. Save active state
	err := saveState(name, state)
	if err != nil {
		t.Fatalf("failed to save window state: %v", err)
	}

	// 3. Load saved state and assert values
	loaded = loadState(name)
	if loaded.Width != state.Width || loaded.Height != state.Height {
		t.Errorf("expected width %d and height %d, got %dx%d", state.Width, state.Height, loaded.Width, loaded.Height)
	}
	if loaded.X != state.X || loaded.Y != state.Y {
		t.Errorf("expected X %d and Y %d, got X %d and Y %d", state.X, state.Y, loaded.X, loaded.Y)
	}
	if loaded.Maximized != state.Maximized {
		t.Errorf("expected Maximized to be %t, got %t", state.Maximized, loaded.Maximized)
	}
}
