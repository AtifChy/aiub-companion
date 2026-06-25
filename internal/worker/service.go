// Package worker implements background tasks such as syncing notices and sending notifications.
package worker

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"aiub-companion/internal/config"
	"aiub-companion/internal/meta"
	"aiub-companion/internal/notice"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/services/notifications"
)

const EventNoticesSynced = "notices:synced"

func init() {
	// Register a custom event whose associated data type is string.
	application.RegisterEvent[int](EventNoticesSynced)
}

type Service struct {
	notice       *notice.Service
	settings     *config.Service
	notification *notifications.NotificationService

	intervalCh chan time.Duration
	cancel     context.CancelFunc
}

func NewService(notice *notice.Service, settings *config.Service, notification *notifications.NotificationService) *Service {
	return &Service{
		notice:       notice,
		settings:     settings,
		notification: notification,
		intervalCh:   make(chan time.Duration, 1),
	}
}

func (s *Service) ServiceStartup(ctx context.Context, _ application.ServiceOptions) error {
	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel

	application.Get().Event.On(config.EventSettingsChanged, func(event *application.CustomEvent) {
		if settings, ok := event.Data.(config.Settings); ok {
			s.intervalCh <- time.Duration(settings.Sync.IntervalMinutes) * time.Minute
		}
	})

	go s.run(ctx)
	return nil
}

func (s *Service) run(ctx context.Context) {
	cfg := s.settings.GetSettings()
	interval := time.Duration(cfg.Sync.IntervalMinutes) * time.Minute
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	doSync := func() {
		count, err := s.notice.SyncNotices(ctx, cfg.Sync.FetchCount)
		if err != nil {
			slog.Error("Failed to sync notices", "error", err)
			return
		}
		if count == 0 {
			slog.Info("No new notices found")
			return
		}

		title := fmt.Sprintf("%d new notice(s)", count)
		err = s.notification.SendNotification(notifications.NotificationOptions{
			ID:    fmt.Sprintf("sync-%d", time.Now().Local().UnixMilli()),
			Title: meta.DisplayName,
			Body:  title + " available",
		})
		if err != nil {
			slog.Error("Failed to send notification", "error", err)
		}

		app := application.Get()
		app.Event.Emit(EventNoticesSynced, count)
		slog.Info("Synced notices", "count", count)
	}

	if cfg.Sync.OnStartup {
		slog.Info("Syncing notices on startup")
		time.AfterFunc(5*time.Second, doSync)
	}

	for {
		select {
		case <-ticker.C:
			doSync()
		case d := <-s.intervalCh:
			ticker.Reset(d)
		case <-ctx.Done():
			slog.Info("Worker service shutting down")
			return
		}
	}
}
