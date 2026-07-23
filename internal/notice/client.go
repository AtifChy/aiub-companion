package notice

import (
	"context"
	"net/http"
	"time"

	"aiub-companion/internal/fetcher"
)

type Client interface {
	ScrapeNotices(ctx context.Context, count int) ([]Notice, error)
	ScrapeNoticeDetails(ctx context.Context, id string) (NoticeDetails, error)
}

func NewClient() Client {
	fetcher := fetcher.New(
		fetcher.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}),
	)
	return &scraper{
		fetcher: fetcher,
		baseURL: defaultBaseURL,
	}
}
