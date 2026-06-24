package main

import (
	"embed"
	"fmt"
	"log/slog"
	"os"
	"time"

	"aiub-companion/internal/database"
	"aiub-companion/internal/log"
	"aiub-companion/internal/meta"
	"aiub-companion/internal/notice"
	"aiub-companion/internal/routine"
	"aiub-companion/internal/settings"
	"aiub-companion/internal/ui"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/services/notifications"
)

//go:embed all:frontend/dist
var assets embed.FS

// build time variables
var version = "dev"

const EventNoticesSynced = "notices:synced"

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
	meta.SetVersion(version)

	// Initialize settings service
	settingsService := settings.NewService()

	syncInterval := make(chan time.Duration, 1)
	applyLogLevel := func(levelStr string) {
		if level, err := settings.ParseLogLevel(levelStr); err != nil {
			slog.Warn("Invalid log level in settings", "error", err)
		} else {
			logger.SetLevel(level)
		}
	}

	settingsService.OnChange(func(new *settings.Settings) {
		syncInterval <- time.Duration(new.Sync.IntervalMinutes) * time.Minute
		applyLogLevel(new.LogLevel)
	})

	// Apply initial settings
	cfg := settingsService.GetSettings()
	applyLogLevel(cfg.LogLevel)

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
	loggerService := log.NewService()
	metaService := meta.NewService()

	// Wails native services
	notificationService := notifications.New()

	// Create a new Wails application by providing the necessary options.
	app := application.New(application.Options{
		Name:        meta.DisplayName,
		Description: meta.Description,
		Services: []application.Service{
			application.NewService(noticeService),
			application.NewService(routineService),
			application.NewService(settingsService),
			application.NewService(loggerService),
			application.NewService(metaService),
			application.NewService(notificationService),
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
			UniqueID: meta.ID,
			AdditionalData: map[string]string{
				"launchTime": time.Now().Format(time.RFC1123),
			},
		},
		Logger: logger.Logger,
	})

	// Setup UI Manager
	manager := ui.NewManager(app, settingsService)
	manager.SetupTray()

	app.Config().SingleInstance.OnSecondInstanceLaunch = func(data application.SecondInstanceData) {
		slog.Info("Second instance launched", "args", data.Args, "dir", data.WorkingDir, "data", data.AdditionalData)
		manager.ShowMainWindow()
	}

	// Show the main window on startup
	if !cfg.Launch.StartMinimized {
		manager.ShowMainWindow()
	}

	// Sync notices on startup and every 30 minutes.
	go startBackgroundSync(app, noticeService, notificationService, cfg, syncInterval)

	// Run the application. This blocks until the application has been exited.
	err = app.Run()
	if err != nil {
		slog.Error("Error running application", "error", err)
		os.Exit(1)
	}
}

func startBackgroundSync(
	app *application.App,
	noticeService *notice.Service,
	notifier *notifications.NotificationService,
	config *settings.Settings,
	syncInterval chan time.Duration,
) {
	interval := time.Duration(config.Sync.IntervalMinutes) * time.Minute
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	doSync := func() {
		newCount, err := noticeService.SyncNotices(app.Context(), config.Sync.FetchCount)
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
			Title: meta.DisplayName,
			Body:  title + " available",
		})
		if err != nil {
			slog.Error("Error sending notification", "error", err)
		}

		app.Event.Emit(EventNoticesSynced, newCount)
		slog.Info("Synced notices", "count", newCount)
	}

	// Run the sync after 5 seconds on startup
	if config.Sync.OnStartup {
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
