// Package updater provides functionality for checking for updates and downloading the latest release from GitHub.
package updater

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"aiub-companion/internal/config"
	"aiub-companion/internal/event"
	"aiub-companion/internal/meta"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/updater"
	"github.com/wailsapp/wails/v3/pkg/updater/providers/github"
	"golang.org/x/mod/semver"
)

var ErrUpdaterNotInitialized = errors.New("updater not initialized")

type Service struct {
	config *config.Service

	stopCh   chan struct{}
	reloadCh chan struct{}

	pendingRelease *Release
	currentVersion string
	githubRepo     string

	mu sync.Mutex
}

func NewService(cfg *config.Service) *Service {
	return &Service{
		currentVersion: meta.Version(),
		githubRepo:     meta.Repo,
		config:         cfg,
		stopCh:         make(chan struct{}),
		reloadCh:       make(chan struct{}, 1),
	}
}

func (s *Service) ServiceStartup(ctx context.Context, _ application.ServiceOptions) error {
	go s.scheduleLoop(ctx)
	return nil
}

func (s *Service) ServiceShutdown(ctx context.Context) error {
	close(s.stopCh)
	return nil
}

//wails:ignore
func (s *Service) Init(app *application.App) error {
	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		slog.Info("Using GitHub token for updater")
	}

	gh, err := github.New(github.Config{
		Repository:    s.githubRepo,
		Token:         token,
		Prerelease:    false,
		ChecksumAsset: "SHA256SUMS",
		HTTPClient:    &http.Client{Timeout: 5 * time.Minute},
	})
	if err != nil {
		return fmt.Errorf("github provider: %w", err)
	}

	err = app.Updater.Init(updater.Config{
		CurrentVersion: s.currentVersion,
		Providers:      []updater.Provider{gh},
		Window:         updater.WindowNone,
		CheckInterval:  0,
	})
	if err != nil {
		return fmt.Errorf("updater init: %w", err)
	}

	state, err := loadState()
	if err != nil {
		slog.Error("Failed to load updater state", "error", err)
	}

	if state.PendingRelease != nil && semver.Compare(state.PendingRelease.Version, s.currentVersion) > 0 {
		s.pendingRelease = state.PendingRelease
	}

	app.Event.On(event.EventConfigChanged, func(_ *application.CustomEvent) {
		s.reloadCh <- struct{}{}
	})

	return nil
}

type Release struct {
	Version string `json:"version"`
	Notes   string `json:"notes"`
}

func (s *Service) CheckForUpdates(ctx context.Context) (*Release, error) {
	app := application.Get()
	if app == nil || app.Updater == nil {
		return nil, ErrUpdaterNotInitialized
	}

	rel, err := app.Updater.Check(ctx)
	if err != nil {
		return nil, fmt.Errorf("check for updates: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	var pending *Release
	if rel != nil {
		pending = &Release{
			Version: rel.Version,
			Notes:   rel.Notes,
		}
	}
	s.pendingRelease = pending

	err = saveState(state{
		LastCheckedAt:  time.Now(),
		PendingRelease: pending,
	})
	if err != nil {
		slog.Error("Failed to save updater state", "error", err)
	}

	if rel == nil {
		slog.Info("No update available")
		return nil, nil
	}

	slog.Info("Update available", "version", rel.Version)

	return pending, nil
}

func (s *Service) DownloadUpdate(ctx context.Context) error {
	app := application.Get()

	if app == nil || app.Updater == nil {
		return ErrUpdaterNotInitialized
	}

	err := app.Updater.DownloadAndInstall(ctx)
	if err != nil {
		return fmt.Errorf("download and install update: %w", err)
	}

	return nil
}

func (s *Service) InstallUpdate(ctx context.Context) error {
	app := application.Get()

	if app == nil || app.Updater == nil {
		return ErrUpdaterNotInitialized
	}

	s.mu.Lock()
	s.pendingRelease = nil
	s.mu.Unlock()

	return app.Updater.Restart(ctx)
}

func (s *Service) GetPendingUpdate(ctx context.Context) (*Release, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.pendingRelease, nil
}

func (s *Service) scheduleLoop(ctx context.Context) {
	for {
		interval := s.checkInterval()
		if interval == 0 {
			select {
			case <-s.stopCh:
				return
			case <-s.reloadCh:
				continue
			}
		}

		state, err := loadState()
		if err != nil {
			slog.Error("Failed to load updater state", "error", err)
		}

		elapsed := time.Since(state.LastCheckedAt)
		wait := max(interval-elapsed, 0)

		timer := time.NewTimer(wait)
		select {
		case <-s.stopCh:
			timer.Stop()
			return
		case <-s.reloadCh:
			timer.Stop()
			continue
		case <-timer.C:
			s.runBackgroundCheck(ctx)
		}
	}
}

func (s *Service) checkInterval() time.Duration {
	cfg := s.config.GetConfig()
	switch cfg.Updates.Interval {
	case config.UpdateIntervalDaily:
		return 24 * time.Hour
	case config.UpdateIntervalWeekly:
		return 7 * 24 * time.Hour
	case config.UpdateIntervalMonthly:
		return 30 * 24 * time.Hour
	default:
		return 0
	}
}

func (s *Service) runBackgroundCheck(ctx context.Context) {
	slog.Info("Running scheduled update check")

	if err := saveState(state{LastCheckedAt: time.Now()}); err != nil {
		slog.Error("Failed to save updater state", "error", err)
	}

	rel, err := s.CheckForUpdates(ctx)
	if err != nil {
		slog.Error("Scheduled update check failed", "error", err)
		return
	}
	if rel == nil {
		slog.Info("No update found during scheduled check")
		return
	}

	slog.Info("Update found during scheduled check", "version", rel.Version)
}

func init() {
	application.RegisterEvent[*updater.Release](updater.EventUpdateAvailable)
	application.RegisterEvent[updater.Progress](updater.EventDownloadProgress)
}
