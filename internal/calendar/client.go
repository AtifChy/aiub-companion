package calendar

import (
	"context"
	"net/http"
	"time"

	"aiub-companion/internal/fetcher"

	"golang.org/x/net/html"
)

type Client interface {
	ScrapeCalendar(ctx context.Context, calType CalendarType) (table *html.Node, semester string, err error)
}

func NewClient() Client {
	return &scraper{
		fetcher: fetcher.New(fetcher.WithHTTPClient(&http.Client{Timeout: 5 * time.Second})),
	}
}
