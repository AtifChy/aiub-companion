package updater

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"aiub-companion/internal/config"
	"aiub-companion/internal/event"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/updater"
	"github.com/wailsapp/wails/v3/pkg/updater/providers/github"
)

var ErrUpdaterNotInitialized = errors.New("updater not initialized")

type Service struct {
	config *config.Service

	stopCh   chan struct{}
	reloadCh chan struct{}

	currentVersion string
	githubRepo     string
}

func NewService(cfg *config.Service) *Service {
	return &Service{
		currentVersion: config.Version(),
		githubRepo:     config.Repo,
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
		Repository:   s.githubRepo,
		Prerelease:   false,
		Token:        token,
		AssetMatcher: assetMatcher,
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

	app.Event.On(updater.EventDownloadProgress, func(event *application.CustomEvent) {
		progress, ok := event.Data.(updater.Progress)
		if !ok {
			return
		}
		slog.Info(fmt.Sprintf("Download progress: %d/%d bytes (%.0f KB/s)", progress.Written, progress.Total, progress.Rate/1024))
	})

	app.Event.On(event.EventConfigChanged, func(_ *application.CustomEvent) {
		s.reloadCh <- struct{}{}
	})

	return nil
}

func assetMatcher(req updater.CheckRequest, assets []github.ReleaseAsset) int {
	var filteredAssets []github.ReleaseAsset
	originalIndices := make([]int, 0, len(assets))

	for i, asset := range assets {
		name := strings.ToLower(asset.Name)
		if strings.Contains(name, "installer") {
			continue
		}
		filteredAssets = append(filteredAssets, asset)
		originalIndices = append(originalIndices, i)
	}

	idx := github.DefaultAssetMatcher(req, filteredAssets)
	if idx == -1 {
		return -1
	}
	return originalIndices[idx]
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

	if err = saveState(state{LastCheckedAt: time.Now()}); err != nil {
		slog.Error("Failed to save updater state", "error", err)
	}

	if rel == nil {
		slog.Info("No update available")
		return nil, nil
	}

	slog.Info("Update available", "version", rel.Version)

	return &Release{
		Version: rel.Version,
		Notes:   rel.Notes,
	}, nil
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

	return app.Updater.Restart(ctx)
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
	case "daily":
		return 24 * time.Hour
	case "weekly":
		return 7 * 24 * time.Hour
	case "monthly":
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
