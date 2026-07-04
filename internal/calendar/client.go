package calendar

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/net/html"
)

type Client interface {
	ScrapeCalendar(ctx context.Context, calType CalendarType) (table *html.Node, semester string, err error)
}

func NewClient() Client {
	return &scraper{
		client: &http.Client{Timeout: 30 * time.Second},
	}
}
