package main

import (
	"embed"
	"log/slog"
	"os"
	"time"

	"aiub-companion/internal/calendar"
	"aiub-companion/internal/config"
	"aiub-companion/internal/database"
	"aiub-companion/internal/desktop"
	"aiub-companion/internal/log"
	"aiub-companion/internal/notice"
	"aiub-companion/internal/routine"
	"aiub-companion/internal/worker"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/services/notifications"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var appIcon []byte

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
	notificationService := notifications.New()

	databaseService := database.NewService()
	configService := config.NewService()

	noticeService := notice.NewService(databaseService)
	routineService := routine.NewService(databaseService)
	calendarService := calendar.NewService(noticeService)

	desktopService := desktop.NewService(configService)
	workerService := worker.NewService(noticeService, configService, notificationService)

	loggerService := log.NewService(logger)

	// Create a new Wails application by providing the necessary options.
	app := application.New(application.Options{
		Name:        config.DisplayName,
		Description: config.Description,
		Icon:        appIcon,
		Services: []application.Service{
			// Core
			application.NewService(loggerService),
			application.NewService(databaseService),
			application.NewService(configService),

			// Independent
			application.NewService(notificationService),

			// Domain logic
			application.NewService(noticeService),
			application.NewService(routineService),
			application.NewService(calendarService),

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
			UniqueID: config.ID,
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
