package notice

import (
	"context"
	"net/http"
	"time"
)

type Client interface {
	ScrapeNotices(ctx context.Context, count int) ([]Notice, error)
	ScrapeNoticeDetails(ctx context.Context, id string) (NoticeDetails, error)
}

func NewClient() Client {
	return &scraper{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: defaultBaseURL,
	}
}
