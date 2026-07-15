package autostart

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"aiub-companion/internal/config"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

const shortcutName = config.AppName + ".lnk"

func Set(enable bool) error {
	dir, err := startupPath()
	if err != nil {
		return fmt.Errorf("getting startup path: %w", err)
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
		return fmt.Errorf("creating startup directory: %w", err)
	}

	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("getting executable path: %w", err)
	}

	return createShortcut(exePath, shortcutPath)
}

func IsEnabled() (bool, error) {
	dir, err := startupPath()
	if err != nil {
		return false, fmt.Errorf("getting startup path: %w", err)
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

func createShortcut(target, shortcutPath string) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	err := ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY)
	if err != nil {
		return fmt.Errorf("CoInitializeEx: %w", err)
	}
	defer ole.CoUninitialize()

	shell, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return fmt.Errorf("creating WScript.Shell object: %w", err)
	}
	defer shell.Release()

	dispatch, err := shell.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return err
	}
	defer dispatch.Release()

	shortcutRaw, err := oleutil.CallMethod(dispatch, "CreateShortcut", shortcutPath)
	if err != nil {
		return fmt.Errorf("CreateShortcut call: %w", err)
	}
	shortcut := shortcutRaw.ToIDispatch()
	defer shortcut.Release()

	if _, err := oleutil.PutProperty(shortcut, "TargetPath", target); err != nil {
		return fmt.Errorf("setting TargetPath: %w", err)
	}
	if _, err := oleutil.PutProperty(shortcut, "WorkingDirectory", filepath.Dir(target)); err != nil {
		return fmt.Errorf("setting WorkingDirectory: %w", err)
	}
	if _, err := oleutil.PutProperty(shortcut, "IconLocation", target); err != nil {
		return fmt.Errorf("setting IconLocation: %w", err)
	}
	if _, err := oleutil.CallMethod(shortcut, "Save"); err != nil {
		return fmt.Errorf("saving shortcut: %w", err)
	}

	return nil
}

func startupPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "Microsoft", "Windows", "Start Menu", "Programs", "Startup"), nil
}
