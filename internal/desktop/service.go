// Package desktop provides the desktop service for the application, managing windows, system tray, and related events.
package desktop

import (
	"context"
	"log/slog"
	"time"

	"aiub-companion/internal/config"
	"aiub-companion/internal/desktop/window"
	"aiub-companion/internal/event"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

type Service struct {
	app    *application.App
	config *config.Service

	main  *window.Window
	about *window.Window
}

func NewService(config *config.Service) *Service {
	return &Service{
		config: config,
	}
}

func (s *Service) ServiceStartup(ctx context.Context, _ application.ServiceOptions) error {
	app := application.Get()

	app.Config().SingleInstance.OnSecondInstanceLaunch = func(data application.SecondInstanceData) {
		slog.Info("Second instance launched", "args", data.Args, "dir", data.WorkingDir, "data", data.AdditionalData)
		s.main.Show()
	}

	app.Event.On(event.EventConfigChanged, func(e *application.CustomEvent) {
		if cfg, ok := e.Data.(config.Config); ok {
			s.main.SetHideOnClose(cfg.Launch.CloseToTray && cfg.Launch.KeepAlive)
			s.main.SetRestoreWindow(cfg.Launch.RestoreWindow)
		}
	})

	app.Event.OnApplicationEvent(events.Common.ApplicationStarted, s.onApplicationStarted)

	s.app = app

	return nil
}

func (s *Service) ShowAboutWindow() {
	s.about.Show()
}

func (s *Service) onApplicationStarted(_ *application.ApplicationEvent) {
	cfg := s.config.GetConfig()
	s.main = window.NewWindow(s.app, window.WindowOptions{
		HideOnClose:   cfg.Launch.CloseToTray && cfg.Launch.KeepAlive,
		RestoreWindow: cfg.Launch.RestoreWindow,
		WebviewWindowOptions: application.WebviewWindowOptions{
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
		},
		OnClose: s.handleClose,
	})

	s.about = window.NewWindow(s.app, window.WindowOptions{
		HideOnClose: true,
		WebviewWindowOptions: application.WebviewWindowOptions{
			Name:                "about-window",
			Title:               "About " + config.DisplayName,
			Width:               400,
			Height:              350,
			Frameless:           true,
			MinimiseButtonState: application.ButtonDisabled,
			MaximiseButtonState: application.ButtonDisabled,
			DisableResize:       true,
			URL:                 "/#/about",
		},
	})

	s.setupTray()

	if !cfg.Launch.StartMinimized {
		s.main.Show()
	}
}

func (s *Service) handleClose() {
	if s.main.HideOnClose() {
		s.about.Hide()
	} else {
		s.about.Close()
	}
	if !s.config.GetConfig().Launch.CloseToTray {
		time.AfterFunc(500*time.Millisecond, func() {
			s.app.Quit()
		})
	}
}

func (s *Service) setupTray() {
	systray := s.app.SystemTray.New()
	systray.SetIcon(s.app.Config().Icon)
	systray.SetLabel(config.DisplayName)
	systray.SetTooltip(config.DisplayName)

	systray.OnDoubleClick(s.main.Show)

	menu := s.app.NewMenu()
	menu.Add("Show").OnClick(func(ctx *application.Context) {
		s.main.Show()
	})
	menu.Add("Quit").OnClick(func(ctx *application.Context) {
		s.app.Quit()
	})

	systray.SetMenu(menu)
}
