// Package worker implements background tasks such as syncing notices and sending notifications.
package worker

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"aiub-companion/internal/calendar"
	"aiub-companion/internal/config"
	"aiub-companion/internal/event"
	"aiub-companion/internal/notice"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/services/notifications"
)

func init() {
	// Register a custom event whose associated data type is string.
	application.RegisterEvent[bool](event.EventNoticeSyncing)
	application.RegisterEvent[int](event.EventNoticeSynced)
	application.RegisterEvent[string](event.EventNoticeOpen)
}

type Service struct {
	notice       *notice.Service
	calendar     *calendar.Service
	config       *config.Service
	notification *notifications.NotificationService

	intervalCh chan time.Duration
	cancel     context.CancelFunc
}

func NewService(notice *notice.Service, calendar *calendar.Service, config *config.Service, notification *notifications.NotificationService) *Service {
	return &Service{
		notice:       notice,
		calendar:     calendar,
		config:       config,
		notification: notification,
		intervalCh:   make(chan time.Duration, 1),
	}
}

func (s *Service) ServiceStartup(ctx context.Context, _ application.ServiceOptions) error {
	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel

	app := application.Get()

	app.Event.On(event.EventConfigChanged, func(ev *application.CustomEvent) {
		if cfg, ok := ev.Data.(config.Config); ok {
			s.intervalCh <- time.Duration(cfg.Sync.Interval) * time.Minute
		}
	})

	s.notification.OnNotificationResponse(s.handleNotificationResponse)

	go s.run(ctx)

	return nil
}

func (s *Service) ServiceShutdown() error {
	if s.cancel != nil {
		s.cancel()
	}
	return nil
}

func (s *Service) run(ctx context.Context) {
	cfg := s.config.GetConfig()
	interval := time.Duration(cfg.Sync.Interval) * time.Minute
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	doSync := func() {
		go s.calendar.SyncAll(ctx)
		s.syncNotices(ctx)
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

func (s *Service) syncNotices(ctx context.Context) {
	app := application.Get()
	cfg := s.config.GetConfig()

	app.Event.Emit(event.EventNoticeSyncing, true)
	defer func() {
		time.Sleep(1 * time.Second)
		app.Event.Emit(event.EventNoticeSyncing, false)
	}()

	newNotices, err := s.notice.SyncNotices(ctx, int(cfg.Sync.FetchCount))
	if err != nil {
		slog.Error("Failed to sync notices", "error", err)
		return
	}

	count := len(newNotices)
	if count == 0 {
		slog.Info("No new notices found")
		return
	}

	payload := notice.BuildNotificationPayload(newNotices)
	err = s.notification.SendNotification(notifications.NotificationOptions{
		ID:    payload.ID,
		Title: payload.Title,
		Body:  payload.Body,
	})
	if err != nil {
		slog.Error("Failed to send notification", "error", err)
	}

	app.Event.Emit(event.EventNoticeSynced, int(count))
	slog.Info("Synced notices", "count", count)
}

func (s *Service) handleNotificationResponse(result notifications.NotificationResult) {
	if result.Error != nil {
		slog.Error("Failed to send notification", "error", result.Error)
		return
	}

	app := application.Get()

	app.Event.Emit(event.EventMainWindowShow)

	id := result.Response.ID
	if !strings.HasPrefix(id, "sync-") {
		s.notice.SetPendingNotice(id)
		app.Event.Emit(event.EventNoticeOpen, id)
		slog.Info("Opening notice from notification", "id", id)
	}
}
