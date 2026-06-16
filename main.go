package main

import (
	"embed"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/services/notifications"

	appinfo "aiub-companion/internal/app"
	"aiub-companion/internal/database"
	"aiub-companion/internal/log"
	"aiub-companion/internal/notice"
	"aiub-companion/internal/routine"
	"aiub-companion/internal/settings"
)

//go:embed all:frontend/dist
var assets embed.FS

// build time variables
var version = "dev"

const EventNoticesSynced = "notices:synced"

var mainWindow *application.WebviewWindow

func init() {
	// Register a custom event whose associated data type is string.
	application.RegisterEvent[int](EventNoticesSynced)
}

func main() {
	// Setup structured logging with slog and our custom log package.
	logger, err := log.SetupLogger()
	if err != nil {
		slog.Error("Failed to setup logger", "error", err)
		os.Exit(1)
	}
	slog.SetDefault(logger.Logger)
	defer func() {
		if err := logger.Close(); err != nil {
			slog.Warn("Failed to close logger", "error", err)
		}
	}()

	// Set application version from build variable
	appinfo.SetVersion(version)

	cfg, err := settings.Load()
	if err != nil {
		slog.Warn("Failed to load settings, using defaults", "error", err)
	}

	if level, err := settings.ParseLogLevel(cfg.LogLevel); err != nil {
		slog.Warn("Invalid log level in settings", "error", err)
	} else {
		logger.SetLevel(level)
	}

	dbInstance, err := database.InitDB()
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := dbInstance.Close(); err != nil {
			slog.Warn("Failed to close database", "error", err)
		}
	}()

	// Initialize Repositories
	noticeRepo := notice.NewRepository(dbInstance.SQLDB)
	routineRepo := routine.NewRepository(dbInstance.SQLDB)

	// Initialize Client
	noticeClient := notice.NewClient()

	// Initialize Services
	noticeService := notice.NewService(noticeRepo, noticeClient)
	routineService := routine.NewService(routineRepo)
	appService := appinfo.NewService()

	syncInterval := make(chan time.Duration, 1)
	settingsService := settings.NewService(logger.Level, func(d time.Duration) {
		syncInterval <- d
	})

	// Wails native services
	notifier := notifications.New()

	// Create a new Wails application by providing the necessary options.
	var app *application.App
	app = application.New(application.Options{
		Name:        appinfo.ID,
		Description: appinfo.Description,
		Services: []application.Service{
			application.NewService(noticeService),
			application.NewService(routineService),
			application.NewService(settingsService),
			application.NewService(appService),
			application.NewService(notifier),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Windows: application.WindowsOptions{
			DisableQuitOnLastWindowClosed: true,
		},
		Linux: application.LinuxOptions{
			DisableQuitOnLastWindowClosed: true,
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
		SingleInstance: &application.SingleInstanceOptions{
			UniqueID: appinfo.ID,
			OnSecondInstanceLaunch: func(data application.SecondInstanceData) {
				slog.Info("Second instance launched", "args", data.Args, "dir", data.WorkingDir, "data", data.AdditionalData)
				showMainWindow(app)
			},
			AdditionalData: map[string]string{
				"launchTime": time.Now().Format(time.RFC1123),
			},
		},
		Logger: logger.Logger,
	})

	// Setup the system tray icon and menu.
	setupSystemTray(app)

	// Show the main window on startup
	if !cfg.Launch.StartMinimized {
		showMainWindow(app)
	}

	// Sync notices on startup and every 30 minutes.
	go startBackgroundSync(app, noticeService, notifier, cfg, syncInterval)

	// Run the application. This blocks until the application has been exited.
	err = app.Run()
	if err != nil {
		slog.Error("Error running application", "error", err)
		os.Exit(1)
	}
}

func showMainWindow(app *application.App) {
	if mainWindow == nil {
		// Create if doesn't exist
		mainWindow = app.Window.NewWithOptions(application.WebviewWindowOptions{
			Title:  appinfo.DisplayName,
			Width:  1024,
			Height: 768,
			Mac: application.MacWindow{
				InvisibleTitleBarHeight: 50,
				Backdrop:                application.MacBackdropTranslucent,
				TitleBar:                application.MacTitleBarHiddenInset,
			},
			BackgroundColour: application.NewRGBA(0, 0, 0, 255),
			URL:              "/",
		})
		mainWindow.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
			mainWindow = nil
		})
	}

	// Show and focus the main window.
	mainWindow.Show()
	mainWindow.Focus()
}

func setupSystemTray(app *application.App) {
	systray := app.SystemTray.New()
	systray.SetIcon(app.Config().Icon)
	systray.SetTooltip(appinfo.DisplayName)
	systray.SetLabel(appinfo.DisplayName)
	systray.OnDoubleClick(func() {
		showMainWindow(app)
	})

	// Create the system tray menu.
	menu := app.NewMenu()
	menu.Add("Show").OnClick(func(ctx *application.Context) {
		showMainWindow(app)
	})
	menu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})
	systray.SetMenu(menu)
}

func startBackgroundSync(
	app *application.App,
	noticeService *notice.Service,
	notifier *notifications.NotificationService,
	cfg *settings.Settings,
	syncInterval chan time.Duration,
) {
	interval := time.Duration(cfg.Sync.IntervalMinutes) * time.Minute
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	doSync := func() {
		newCount, err := noticeService.SyncNotices(app.Context(), cfg.Sync.NoticeFetchCount)
		if err != nil {
			slog.Error("Error syncing notices", "error", err)
			return
		}
		if newCount == 0 {
			slog.Warn("No new notices to sync")
			return
		}

		title := fmt.Sprintf("%d new notice(s)", newCount)
		err = notifier.SendNotification(notifications.NotificationOptions{
			ID:    fmt.Sprintf("sync-%d", time.Now().Local().UnixMilli()),
			Title: appinfo.DisplayName,
			Body:  title + " available",
		})
		if err != nil {
			slog.Error("Error sending notification", "error", err)
		}

		app.Event.Emit(EventNoticesSynced, newCount)
		slog.Info("Synced notices", "count", newCount)
	}

	// Run the sync after 5 seconds on startup
	if cfg.Sync.OnStartup {
		slog.Info("Syncing notices on startup")
		time.AfterFunc(5*time.Second, doSync)
	}

	for {
		select {
		case <-ticker.C:
			doSync()
		case d := <-syncInterval:
			ticker.Reset(d)
		}
	}
}
