// Package fetcher provides a Fetcher type that can fetch HTML documents from URLs with retry logic and customizable options.
package fetcher

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/html"
)

const DefaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148.0.0.0 Safari/537.36"

type Option func(*Fetcher)

type Fetcher struct {
	client         *http.Client
	userAgent      string
	maxAttempts    int
	initialBackoff time.Duration
}

func New(opts ...Option) *Fetcher {
	fetcher := &Fetcher{
		client:         &http.Client{Timeout: 15 * time.Second},
		userAgent:      DefaultUserAgent,
		maxAttempts:    5,
		initialBackoff: 1 * time.Second,
	}
	for _, opt := range opts {
		opt(fetcher)
	}
	return fetcher
}

func WithHTTPClient(client *http.Client) Option {
	return func(f *Fetcher) { f.client = client }
}

func WithUserAgent(userAgent string) Option {
	return func(f *Fetcher) { f.userAgent = userAgent }
}

func WithMaxAttempts(attempts int) Option {
	return func(f *Fetcher) { f.maxAttempts = attempts }
}

func (f *Fetcher) FetchHTML(ctx context.Context, url string) (*html.Node, error) {
	backoff := f.initialBackoff

	for attempt := 1; attempt <= f.maxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		doc, retriable, err := f.fetchWithRetriable(ctx, url)
		if err == nil {
			return doc, nil
		}

		if attempt == f.maxAttempts || !retriable {
			if attempt == 1 {
				return nil, err
			}
			return nil, fmt.Errorf("after %d attempts: %w", attempt, err)
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(backoff):
			backoff *= 2
		}
	}

	panic("FetchHTML: unreachable code reached")
}

func (f *Fetcher) fetchWithRetriable(ctx context.Context, url string) (*html.Node, bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, false, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("User-Agent", f.userAgent)

	resp, err := f.client.Do(req)
	if err != nil {
		if netErr, ok := errors.AsType[net.Error](err); ok && netErr.Timeout() {
			return nil, true, fmt.Errorf("%s: timeout: %w", url, err)
		}
		if dnsErr, ok := errors.AsType[*net.DNSError](err); ok {
			return nil, true, fmt.Errorf("%s: DNS error: %w", url, dnsErr)
		}
		if opErr, ok := errors.AsType[*net.OpError](err); ok {
			return nil, true, fmt.Errorf("%s: operation error: %w", url, opErr)
		}
		return nil, false, fmt.Errorf("%s: %w", url, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		is5xx := resp.StatusCode >= 500 && resp.StatusCode < 600
		return nil, is5xx, fmt.Errorf("%s: unexpected status %d", url, resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, false, fmt.Errorf("%s: parse HTML: %w", url, err)
	}

	return doc, false, nil
}
