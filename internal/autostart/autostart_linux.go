package autostart

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"aiub-companion/internal/meta"
)

const shortcutName = meta.AppName + ".desktop"

// Set set the application to start on login in linux
func Set(enable bool) error {
	dir, err := autostartPath()
	if err != nil {
		return fmt.Errorf("getting autostart path: %w", err)
	}
	shortcutPath := filepath.Join(dir, shortcutName)

	if !enable {
		err := os.Remove(shortcutPath)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("removing shortcut: %w", err)
		}
		return nil
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating autostart directory: %w", err)
	}

	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("getting executable path: %w", err)
	}

	content := fmt.Sprintf(`[Desktop Entry]
Type=Application
Name=%s
Exec=%s
Terminal=false
X-GNOME-Autostart-enabled=true
`, meta.DisplayName, exePath)

	return os.WriteFile(shortcutPath, []byte(content), 0o644)
}

func IsEnabled() (bool, error) {
	dir, err := autostartPath()
	if err != nil {
		return false, fmt.Errorf("getting autostart path: %w", err)
	}
	shortcutPath := filepath.Join(dir, shortcutName)

	_, err = os.Stat(shortcutPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func autostartPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "autostart"), nil
}
