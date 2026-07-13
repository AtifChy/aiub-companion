package updater

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/updater"
	"github.com/wailsapp/wails/v3/pkg/updater/providers/github"
)

type Service struct {
	currentVersion string
	githubRepo     string
}

func NewService(currentVersion, githubRepo string) *Service {
	return &Service{
		currentVersion: currentVersion,
		githubRepo:     githubRepo,
	}
}

//wails:ignore
func (s *Service) Init(app *application.App) error {
	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		slog.Info("Using GitHub token for updater")
	}

	gh, err := github.New(github.Config{
		Repository: s.githubRepo,
		Prerelease: false,
		Token:      token,
		AssetMatcher: func(req updater.CheckRequest, assets []github.ReleaseAsset) int {
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
		},
	})
	if err != nil {
		return fmt.Errorf("github provider: %w", err)
	}

	err = app.Updater.Init(updater.Config{
		CurrentVersion: s.currentVersion,
		Providers:      []updater.Provider{gh},
		Window:         updater.WindowNone,
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

	return nil
}

type Release struct {
	Version string `json:"version"`
	Notes   string `json:"notes"`
}

func (s *Service) CheckForUpdates(ctx context.Context) (*Release, error) {
	app := application.Get()

	if app == nil || app.Updater == nil {
		return nil, fmt.Errorf("updater not initialized")
	}

	rel, err := app.Updater.Check(ctx)
	if err != nil {
		return nil, fmt.Errorf("check for updates: %w", err)
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
		return fmt.Errorf("updater not initialized")
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
		return fmt.Errorf("updater not initialized")
	}

	return app.Updater.Restart(ctx)
}

func init() {
	application.RegisterEvent[*updater.Release](updater.EventUpdateAvailable)
	application.RegisterEvent[updater.Progress](updater.EventDownloadProgress)
}
