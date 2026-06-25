package main

import (
	"embed"
	"log/slog"
	"os"
	"time"

	"aiub-companion/internal/database"
	"aiub-companion/internal/desktop"
	"aiub-companion/internal/log"
	"aiub-companion/internal/meta"
	"aiub-companion/internal/notice"
	"aiub-companion/internal/routine"
	"aiub-companion/internal/settings"
	"aiub-companion/internal/worker"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/services/notifications"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var appIcon []byte

// build time variables
var (
	version   = "dev"
	commit    string
	buildTime string
)

func main() {
	// Setup structured logging with slog and our custom log package.
	logger, err := log.SetupLogger()
	if err != nil {
		slog.Error("Failed to setup logger", "error", err)
		os.Exit(1)
	}
	slog.SetDefault(logger.Logger)
	defer func() {
		_ = logger.Close()
	}()

	// Initialize Services
	metaService := meta.NewService()
	notificationService := notifications.New()

	databaseService := database.NewService()
	settingsService := settings.NewService()

	noticeService := notice.NewService(databaseService)
	routineService := routine.NewService(databaseService)

	desktopService := desktop.NewService(settingsService)
	workerService := worker.NewService(noticeService, settingsService, notificationService)

	loggerService := log.NewService(logger)

	// Create a new Wails application by providing the necessary options.
	app := application.New(application.Options{
		Name:        meta.DisplayName,
		Description: meta.Description,
		Icon:        appIcon,
		Services: []application.Service{
			// Core
			application.NewService(loggerService),
			application.NewService(databaseService),
			application.NewService(settingsService),

			// Independent
			application.NewService(metaService),
			application.NewService(notificationService),

			// Domain logic
			application.NewService(noticeService),
			application.NewService(routineService),

			// Consumers
			application.NewService(desktopService),
			application.NewService(workerService),
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

	// Run the application. This blocks until the application has been exited.
	err = app.Run()
	if err != nil {
		slog.Error("Error running application", "error", err)
		os.Exit(1)
	}
}
