package notice

import (
	"context"

	"aiub-companion/internal/fetcher"
)

type Client interface {
	ScrapeNotices(ctx context.Context, count int) ([]Notice, error)
	ScrapeNoticeDetails(ctx context.Context, id string) (NoticeDetails, error)
}

func NewClient() Client {
	fetcher := fetcher.New()
	return &scraper{
		fetcher: fetcher,
		baseURL: defaultBaseURL,
	}
}
