package autostart

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupTextConfigDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	return dir
}

func TestIsEnabled_InitiallyDisabled(t *testing.T) {
	setupTextConfigDir(t)

	enabled, err := IsEnabled()
	if err != nil {
		t.Fatalf("IsEnabled: %v", err)
	}
	if enabled {
		t.Fatalf("expected IsEnabled to be false, got true")
	}
}

func TestSet_Enable_CreatesDesktopFile(t *testing.T) {
	configDir := setupTextConfigDir(t)

	if err := Set(true); err != nil {
		t.Fatalf("Set(true) returned error: %v", err)
	}

	path := filepath.Join(configDir, "autostart", shortcutName)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("expected desktop file to exist at %s, but got error: %v", path, err)
	}

	content := string(data)
	wantSubstrings := []string{
		"[Desktop Entry]",
		"Type=Application",
		"Name=",
		"Exec=",
	}
	for _, substr := range wantSubstrings {
		if !strings.Contains(content, substr) {
			t.Fatalf("expected desktop file to contain %q, but it did not", substr)
		}
	}
}

func TestSet_Enable_ExecPathMatchesExecutable(t *testing.T) {
	setupTextConfigDir(t)

	if err := Set(true); err != nil {
		t.Fatalf("Set(true) returned error: %v", err)
	}

	enabled, err := IsEnabled()
	if err != nil {
		t.Fatalf("IsEnabled returned error: %v", err)
	}
	if !enabled {
		t.Fatalf("expected IsEnabled to be true after Set(true), got false")
	}

	exePath, err := os.Executable()
	if err != nil {
		t.Fatalf("os.Executable returned error: %v", err)
	}

	dir, _ := autostartPath()
	data, err := os.ReadFile(filepath.Join(dir, shortcutName))
	if err != nil {
		t.Fatalf("reading desktop file returned error: %v", err)
	}
	if !strings.Contains(string(data), "Exec="+exePath) {
		t.Fatalf("expected desktop file to contain Exec=%s, but it did not", exePath)
	}
}

func TestSet_enable_Idempotent(t *testing.T) {
	setupTextConfigDir(t)

	if err := Set(true); err != nil {
		t.Fatalf("Set(true) returned error: %v", err)
	}
	if err := Set(true); err != nil {
		t.Fatalf("Set(true) called again returned error: %v", err)
	}

	enabled, err := IsEnabled()
	if err != nil {
		t.Fatalf("IsEnabled returned error: %v", err)
	}
	if !enabled {
		t.Fatalf("expected IsEnabled to be true after enabling, got false")
	}
}

func TestSet_Disable_RemovesDesktopFile(t *testing.T) {
	setupTextConfigDir(t)

	if err := Set(true); err != nil {
		t.Fatalf("Set(true) returned error: %v", err)
	}
	if err := Set(false); err != nil {
		t.Fatalf("Set(false) returned error: %v", err)
	}

	enabled, err := IsEnabled()
	if err != nil {
		t.Fatalf("IsEnabled returned error: %v", err)
	}
	if enabled {
		t.Fatalf("expected IsEnabled to be false after disabling, got true")
	}
}

func TestSet_Disable_NoOpWhenNotEnabled(t *testing.T) {
	setupTextConfigDir(t)

	if err := Set(false); err != nil {
		t.Fatalf("Set(false) returned error: %v", err)
	}
}

func TestSet_Enable_CreatesAutostartDirIfNotExists(t *testing.T) {
	configDir := setupTextConfigDir(t)

	_, err := os.Stat(filepath.Join(configDir, "autostart"))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("unexpected error checking autostart dir: %v", err)
	}

	if err := Set(true); err != nil {
		t.Fatalf("Set(true) returned error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(configDir, "autostart")); err != nil {
		t.Fatalf("expected autostart dir to exist after Set(true), but got error: %v", err)
	}
}
