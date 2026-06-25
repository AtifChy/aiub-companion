// Package desktop implements the desktop application service, including window management and system tray integration.
package desktop

import (
	"context"
	"log/slog"

	"aiub-companion/internal/config"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

type Service struct {
	app    *application.App
	window *application.WebviewWindow
	config *config.Service
	state  config.WindowState
}

func NewService(config *config.Service) *Service {
	return &Service{
		config: config,
	}
}

func (s *Service) ServiceStartup(ctx context.Context, _ application.ServiceOptions) error {
	s.app = application.Get()

	s.app.Config().SingleInstance.OnSecondInstanceLaunch = func(data application.SecondInstanceData) {
		slog.Info("Second instance launched", "args", data.Args, "dir", data.WorkingDir, "data", data.AdditionalData)
		s.showMainWindow()
	}

	s.app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(event *application.ApplicationEvent) {
		cfg := s.config.GetConfig()
		s.state = cfg.Window

		s.setupTray()

		if !cfg.Launch.StartMinimized {
			s.showMainWindow()
		}
	})

	return nil
}

func (s *Service) showMainWindow() {
	if s.window == nil {
		s.createMainWindow()
		s.loadState()
		s.setupEventHandlers()
	}
	s.window.Show()
	s.window.Focus()
}

func (s *Service) setupTray() {
	systray := s.app.SystemTray.New()
	systray.SetIcon(s.app.Config().Icon)
	systray.SetLabel(config.DisplayName)
	systray.SetTooltip(config.DisplayName)

	systray.OnDoubleClick(func() {
		s.showMainWindow()
	})

	menu := s.app.NewMenu()
	menu.Add("Show").OnClick(func(ctx *application.Context) {
		s.showMainWindow()
	})
	menu.Add("Quit").OnClick(func(ctx *application.Context) {
		s.app.Quit()
	})
	systray.SetMenu(menu)
}

func (s *Service) createMainWindow() {
	s.window = s.app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:             "main-window",
		Title:            config.DisplayName,
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

func (s *Service) setupEventHandlers() {
	// State change events
	s.window.OnWindowEvent(events.Common.WindowMinimise, func(event *application.WindowEvent) {
		if err := s.config.UpdateWindowState(&s.state); err != nil {
			slog.Error("Failed to save window state on minimize", "error", err)
		}
	})
	s.window.OnWindowEvent(events.Common.WindowMaximise, func(event *application.WindowEvent) {
		s.state.Maximized = true
	})
	s.window.OnWindowEvent(events.Common.WindowUnMaximise, func(event *application.WindowEvent) {
		s.state.Maximized = false
	})

	// Position and size events
	s.window.OnWindowEvent(events.Common.WindowDidResize, func(event *application.WindowEvent) {
		if !s.state.Maximized {
			s.state.Width, s.state.Height = s.window.Size()
		}
	})
	s.window.OnWindowEvent(events.Common.WindowDidMove, func(event *application.WindowEvent) {
		if !s.state.Maximized {
			s.state.X, s.state.Y = s.window.Position()
		}
	})

	// Save state on close
	s.window.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
		if err := s.config.UpdateWindowState(&s.state); err != nil {
			slog.Error("Failed to save window state on close", "error", err)
		}
		if s.config.GetConfig().Launch.CloseToTray {
			s.window = nil
		} else {
			event.Cancel()
			s.window.Close()
			s.app.Quit()
		}
	})
}

func (s *Service) loadState() {
	if s.window == nil || !s.config.GetConfig().Launch.RestoreWindow {
		return
	}
	s.window.SetSize(s.state.Width, s.state.Height)
	if s.state.X != -1 && s.state.Y != -1 {
		s.window.SetPosition(s.state.X, s.state.Y)
	}
	if s.state.Maximized {
		s.window.Maximise()
	}
}
