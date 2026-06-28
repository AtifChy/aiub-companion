// Package window provides a thin wrapper around Wails WebviewWindow to manage window lifecycle and state persistence.
package window

import (
	"log/slog"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// WindowOptions holds the configuration for a managed window.
type WindowOptions struct {
	OnClose func()

	application.WebviewWindowOptions

	HideOnClose   bool
	RestoreWindow bool
}

// Window is a thin lifecycle wrapper around a Wails WebviewWindow.
// It is safe to call Show and Hide concurrently, but the underlying WebviewWindow is not thread-safe.
type Window struct {
	app    *application.App
	handle *application.WebviewWindow
	timer  *time.Timer
	opts   WindowOptions
	state  WindowState
}

// NewWindow creates a Window descriptor. No OS window is created until Show is called.
func NewWindow(app *application.App, opts WindowOptions) *Window {
	return &Window{
		app:  app,
		opts: opts,
	}
}

// SetHideOnClose sets the HideOnClose option for the window.
// If true, closing the window will hide it instead of destroying it.
func (w *Window) SetHideOnClose(hide bool) {
	w.opts.HideOnClose = hide
}

// HideOnClose returns the current value of the HideOnClose option.
func (w *Window) HideOnClose() bool {
	return w.opts.HideOnClose
}

// SetRestoreWindow sets the RestoreWindow option for the window.
// If true, the window will restore its size and position from the last session.
func (w *Window) SetRestoreWindow(restore bool) {
	w.opts.RestoreWindow = restore
}

// Show brings the window to foreground, creating it if necessary.
func (w *Window) Show() {
	if w.handle == nil {
		w.handle = w.app.Window.NewWithOptions(w.opts.WebviewWindowOptions)

		w.handle.RegisterHook(events.Common.WindowClosing, w.onClose)

		w.restoreState()
		w.setupStateTracking()
	}

	w.handle.Show()
	w.handle.Focus()
}

// Hide hides the window without destroying it. No-op if not yet created.
func (w *Window) Hide() {
	if w.handle != nil {
		w.handle.Hide()
	}
}

// Close destroys the window if it exists. No-op if not yet created.
func (w *Window) Close() {
	if w.handle != nil {
		w.handle.Close()
		w.handle = nil
	}
}

// onClose is called from the WindowClosing event. Must NOT be called when mu is held.
func (w *Window) onClose(e *application.WindowEvent) {
	if w.opts.OnClose != nil {
		w.opts.OnClose()
	}

	if w.opts.HideOnClose {
		e.Cancel()
		if w.handle != nil {
			w.handle.Hide()
		}
	} else {
		w.handle = nil
	}
}

func (w *Window) restoreState() {
	if !w.opts.RestoreWindow || w.opts.Name == "" {
		return
	}

	w.state = loadState(w.opts.Name)

	w.handle.SetSize(w.state.Width, w.state.Height)

	if w.state.X != -1 && w.state.Y != -1 {
		w.handle.SetPosition(w.state.X, w.state.Y)
	}

	if w.state.Maximized {
		w.handle.Maximise()
	}
}

func (w *Window) setupStateTracking() {
	if !w.opts.RestoreWindow || w.opts.Name == "" {
		return
	}

	w.handle.OnWindowEvent(events.Common.WindowMinimise, func(event *application.WindowEvent) {
		w.flushSaveState()
	})
	w.handle.OnWindowEvent(events.Common.WindowMaximise, func(event *application.WindowEvent) {
		w.state.Maximized = true
		w.flushSaveState()
	})
	w.handle.OnWindowEvent(events.Common.WindowUnMaximise, func(event *application.WindowEvent) {
		w.state.Maximized = false
		w.flushSaveState()
	})

	w.handle.OnWindowEvent(events.Common.WindowDidResize, func(event *application.WindowEvent) {
		if !w.state.Maximized {
			w.state.Width, w.state.Height = w.handle.Size()
			w.scheduleSaveState()
		}
	})
	w.handle.OnWindowEvent(events.Common.WindowDidMove, func(event *application.WindowEvent) {
		if !w.state.Maximized {
			w.state.X, w.state.Y = w.handle.Position()
			w.scheduleSaveState()
		}
	})
}

func (w *Window) scheduleSaveState() {
	if w.timer != nil {
		w.timer.Stop()
	}

	w.timer = time.AfterFunc(500*time.Millisecond, func() {
		if err := saveState(w.opts.Name, w.state); err != nil {
			slog.Error("Failed to save window state", "error", err)
		}
	})
}

func (w *Window) flushSaveState() {
	if w.timer != nil {
		w.timer.Stop()
	}

	if err := saveState(w.opts.Name, w.state); err != nil {
		slog.Error("Failed to save window state", "error", err)
	}
}
