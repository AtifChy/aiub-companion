package notice

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"aiub-companion/internal/database"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type Service struct {
	db     *database.Service
	repo   Repository
	client Client
}

func NewService(db *database.Service) *Service {
	return &Service{db: db}
}

func (s *Service) ServiceStartup(ctx context.Context, _ application.ServiceOptions) error {
	s.repo = NewRepository(s.db.DB())
	s.client = NewClient()
	return nil
}

func (s *Service) SyncNotices(ctx context.Context, count int) (int64, error) {
	notices, err := s.client.ScrapeNotices(ctx, count)
	if err != nil {
		return 0, fmt.Errorf("scrape notices: %w", err)
	}

	var newCount int64

	err = s.repo.WithTx(ctx, func(txRepo Repository) error {
		latestOrder, err := txRepo.GetLatestNoticeSourceOrder(ctx)
		if err != nil {
			return fmt.Errorf("get latest notice info: %w", err)
		}
		base := latestOrder + 1

		for i := range slices.Backward(notices) {
			orderOffset := int64(len(notices) - 1 - i)
			n, err := txRepo.InsertNotice(ctx, notices[i], base+orderOffset)
			if err != nil {
				return fmt.Errorf("insert notice %s: %w", notices[i].ID, err)
			}

			newCount += n

			if n == 0 {
				err = txRepo.UpdateNotice(ctx, notices[i])
				if err != nil {
					return fmt.Errorf("update notice %s: %w", notices[i].ID, err)
				}
			}
		}
		return nil
	})

	return newCount, err
}

func (s *Service) GetNotices(ctx context.Context, filter Filter) ([]Notice, error) {
	rawQuery := strings.TrimSpace(filter.Search)
	filter.Search = ""

	notices, err := s.repo.GetNotices(ctx, filter)
	if err != nil {
		return nil, err
	}

	if rawQuery != "" {
		notices = fuzzyFilter(notices, rawQuery)
	}

	return notices, nil
}

func (s *Service) GetNoticeDetails(ctx context.Context, id string) (*Notice, error) {
	notice, err := s.repo.GetNoticeByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get notice: %w", err)
	}

	if !notice.IsCached {
		if err := s.populateNoticeCache(ctx, notice); err != nil {
			return nil, fmt.Errorf("populate notice cache: %w", err)
		}
	} else {
		notice.Attachments, err = s.repo.GetNoticeAttachments(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("get notice attachments: %w", err)
		}
	}

	return notice, nil
}

func (s *Service) populateNoticeCache(ctx context.Context, notice *Notice) error {
	details, err := s.client.ScrapeNoticeDetails(ctx, notice.ID)
	if err != nil {
		return fmt.Errorf("scrape notice details: %w", err)
	}

	err = s.repo.WithTx(ctx, func(txRepo Repository) error {
		err := txRepo.UpdateNoticeDetails(ctx, notice.ID, details.FullTitle, details.Content)
		if err != nil {
			return fmt.Errorf("update notice details: %w", err)
		}
		for i := range details.Attachments {
			err := txRepo.UpsertNoticeAttachment(ctx, notice.ID, details.Attachments[i])
			if err != nil {
				return fmt.Errorf("upsert notice attachment: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	// Update the notice object in memory
	notice.Title = details.FullTitle
	notice.Content = details.Content
	notice.Attachments = details.Attachments

	return nil
}

func (s *Service) ToggleNoticePinned(ctx context.Context, id string, pinned bool) error {
	return s.repo.SetPinState(ctx, id, pinned)
}

func (s *Service) ToggleNoticeRead(ctx context.Context, id string, read bool) error {
	return s.repo.SetReadState(ctx, id, read)
}

func (s *Service) ClearNotices(ctx context.Context) error {
	return s.repo.ClearNotices(ctx)
}
