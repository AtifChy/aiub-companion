package autostart

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func setupTestAppData(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("APPDATA", dir)
	return dir
}

func readShortcutTarget(t *testing.T, path string) string {
	t.Helper()

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if err := ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY); err != nil {
		t.Fatalf("CoInitializeEx: %v", err)
	}
	defer ole.CoUninitialize()

	shell, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		t.Fatalf("CreateObject: %v", err)
	}
	defer shell.Release()

	wshell, err := shell.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		t.Fatalf("QueryInterface: %v", err)
	}
	defer wshell.Release()

	shortcut, err := oleutil.CallMethod(wshell, "CreateShortcut", path)
	if err != nil {
		t.Fatalf("CreateShortcut: %v", err)
	}
	shortcutDispatch := shortcut.ToIDispatch()
	defer shortcutDispatch.Release()

	target, err := oleutil.GetProperty(shortcutDispatch, "TargetPath")
	if err != nil {
		t.Fatalf("GetProperty TargetPath: %v", err)
	}
	return target.ToString()
}

func TestIsEnabled_InitiallyDisabled(t *testing.T) {
	setupTestAppData(t)

	s := NewService()
	enabled, err := s.IsEnabled()
	if err != nil {
		t.Fatalf("IsEnabled returned error: %v", err)
	}
	if enabled {
		t.Fatalf("expected IsEnabled to return false, got true")
	}
}

func TestSet_Enable_CreatesShortcut(t *testing.T) {
	appData := setupTestAppData(t)

	s := NewService()
	if err := s.Set(true); err != nil {
		t.Fatalf("Set(true) returned error: %v", err)
	}

	path := filepath.Join(appData, "Microsoft", "Windows", "Start Menu", "Programs", "Startup", shortcutName)
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected shortcut to exist at %s, but got error: %v", path, err)
	}
}

func TestSet_Enable_TargetMatchesExecutable(t *testing.T) {
	appData := setupTestAppData(t)

	s := NewService()
	if err := s.Set(true); err != nil {
		t.Fatalf("Set(true) returned error: %v", err)
	}

	exePath, err := os.Executable()
	if err != nil {
		t.Fatalf("os.Executable returned error: %v", err)
	}

	path := filepath.Join(appData, "Microsoft", "Windows", "Start Menu", "Programs", "Startup", shortcutName)
	target := readShortcutTarget(t, path)

	if target != exePath {
		t.Fatalf("expected shortcut target to be %s, got %s", exePath, target)
	}
}

func TestSet_Enable_Idempotent(t *testing.T) {
	setupTestAppData(t)

	s := NewService()
	if err := s.Set(true); err != nil {
		t.Fatalf("Set(true) returned error: %v", err)
	}
	if err := s.Set(true); err != nil {
		t.Fatalf("Set(true) called again returned error: %v", err)
	}

	enabled, err := s.IsEnabled()
	if err != nil {
		t.Fatalf("IsEnabled returned error: %v", err)
	}
	if !enabled {
		t.Fatalf("expected IsEnabled to return true after enabling, got false")
	}
}

func TestSet_Disable_RemovesShortcut(t *testing.T) {
	setupTestAppData(t)

	s := NewService()
	if err := s.Set(true); err != nil {
		t.Fatalf("Set(true) returned error: %v", err)
	}
	if err := s.Set(false); err != nil {
		t.Fatalf("Set(false) returned error: %v", err)
	}

	enabled, err := s.IsEnabled()
	if err != nil {
		t.Fatalf("IsEnabled returned error: %v", err)
	}
	if enabled {
		t.Fatalf("expected IsEnabled to return false after disabling, got true")
	}
}

func TestSet_Disable_NoOpWhenNotEnabled(t *testing.T) {
	setupTestAppData(t)

	s := NewService()
	if err := s.Set(false); err != nil {
		t.Fatalf("Set(false) returned error: %v", err)
	}

	enabled, err := s.IsEnabled()
	if err != nil {
		t.Fatalf("IsEnabled returned error: %v", err)
	}
	if enabled {
		t.Fatalf("expected IsEnabled to return false after disabling when not enabled, got true")
	}
}
