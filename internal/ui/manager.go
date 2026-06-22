package ui

import (
	"log/slog"

	"aiub-companion/internal/meta"
	"aiub-companion/internal/settings"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

type Manager struct {
	app      *application.App
	window   *application.WebviewWindow
	settings *settings.Service
	state    *settings.WindowState
}

func NewManager(app *application.App, settings *settings.Service) *Manager {
	state := settings.GetSettings().Window
	return &Manager{
		app:      app,
		settings: settings,
		state:    &state,
	}
}

func (m *Manager) ShowMainWindow() {
	if m.window == nil {
		m.createWindow()
		m.loadState()
		m.setupEventHandlers()
	}
	m.window.Show()
	m.window.Focus()
}

func (m *Manager) SetupTray() {
	systray := m.app.SystemTray.New()
	systray.SetIcon(m.app.Config().Icon)
	systray.SetLabel(meta.DisplayName)
	systray.SetTooltip(meta.DisplayName)

	systray.OnDoubleClick(func() {
		m.ShowMainWindow()
	})

	menu := m.app.NewMenu()
	menu.Add("Show").OnClick(func(ctx *application.Context) {
		m.ShowMainWindow()
	})
	menu.Add("Quit").OnClick(func(ctx *application.Context) {
		m.app.Quit()
	})
	systray.SetMenu(menu)
}

func (m *Manager) createWindow() {
	m.window = m.app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            meta.DisplayName,
		Frameless:        true,
		BackgroundColour: application.NewRGBA(0, 0, 0, 255),
		URL:              "/",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
	})
}

func (m *Manager) setupEventHandlers() {
	// State change events
	m.window.OnWindowEvent(events.Common.WindowMinimise, func(event *application.WindowEvent) {
		if err := m.settings.UpdateWindowState(m.state); err != nil {
			slog.Error("Failed to save window state on minimize", "error", err)
		}
	})
	m.window.OnWindowEvent(events.Common.WindowMaximise, func(event *application.WindowEvent) {
		m.state.Maximized = true
	})
	m.window.OnWindowEvent(events.Common.WindowUnMaximise, func(event *application.WindowEvent) {
		m.state.Maximized = false
	})

	// Position and size events
	m.window.OnWindowEvent(events.Common.WindowDidResize, func(event *application.WindowEvent) {
		if !m.state.Maximized {
			m.state.Width, m.state.Height = m.window.Size()
		}
	})
	m.window.OnWindowEvent(events.Common.WindowDidMove, func(event *application.WindowEvent) {
		if !m.state.Maximized {
			m.state.X, m.state.Y = m.window.Position()
		}
	})

	// Save state on close
	m.window.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
		if err := m.settings.UpdateWindowState(m.state); err != nil {
			slog.Error("Failed to save window state on close", "error", err)
		}
		if m.settings.GetSettings().Launch.CloseToTray {
			m.window = nil
		} else {
			event.Cancel()
			m.app.Quit()
		}
	})
}

func (m *Manager) loadState() {
	if m.window == nil {
		return
	}
	m.window.SetSize(m.state.Width, m.state.Height)
	if m.state.X != -1 && m.state.Y != -1 {
		m.window.SetPosition(m.state.X, m.state.Y)
	}
	if m.state.Maximized {
		m.window.Maximise()
	}
}
