package calendar

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

type mockRoundTripper struct {
	roundTrip func(req *http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTrip(req)
}

func TestScraper_ScrapeCalendar_Success(t *testing.T) {
	mockHTML := `
	<html>
		<body>
			<div id="lqd-tab-scedule-item-2">
				<h2>&nbsp;</h2>
				<h2>Fall 2025-2026</h2>
				<table>
					<tr><td>Standard Content</td></tr>
				</table>
			</div>
			<div id="lqd-tab-scedule-item-1">
				<h2>&nbsp;</h2>
				<h2>Fall 2025-2026 [Revised]</h2>
				<table>
					<tr><td>LLB Content</td></tr>
				</table>
			</div>
		</body>
	</html>
	`

	httpClient := &http.Client{
		Transport: &mockRoundTripper{
			roundTrip: func(req *http.Request) (*http.Response, error) {
				if req.URL.String() != calendarURL {
					return nil, errors.New("unexpected URL: " + req.URL.String())
				}
				if req.Header.Get("User-Agent") == "" {
					return nil, errors.New("User-Agent header missing")
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(mockHTML)),
				}, nil
			},
		},
	}

	s := &scraper{client: httpClient}

	t.Run("Standard Calendar", func(t *testing.T) {
		table, semester, err := s.ScrapeCalendar(context.Background(), CalendarStandard)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if semester != "fall 2025-2026" {
			t.Errorf("unexpected semester: %q", semester)
		}
		if table == nil {
			t.Fatal("expected table to be non-nil")
		}
	})

	t.Run("LLB Calendar", func(t *testing.T) {
		table, semester, err := s.ScrapeCalendar(context.Background(), CalendarLLBBPharm)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if semester != "fall 2025-2026 [revised]" {
			t.Errorf("unexpected semester: %q", semester)
		}
		if table == nil {
			t.Fatal("expected table to be non-nil")
		}
	})
}

func TestScraper_ScrapeCalendar_Errors(t *testing.T) {
	t.Run("Unsupported Calendar Type", func(t *testing.T) {
		s := &scraper{client: http.DefaultClient}
		_, _, err := s.ScrapeCalendar(context.Background(), CalendarType("invalid"))
		if err == nil {
			t.Error("expected error for unsupported calendar type")
		}
	})

	t.Run("HTTP Error Status", func(t *testing.T) {
		httpClient := &http.Client{
			Transport: &mockRoundTripper{
				roundTrip: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusInternalServerError,
						Body:       io.NopCloser(bytes.NewReader(nil)),
					}, nil
				},
			},
		}
		s := &scraper{client: httpClient}
		_, _, err := s.ScrapeCalendar(context.Background(), CalendarStandard)
		if err == nil {
			t.Error("expected error for non-200 status code")
		}
	})

	t.Run("Missing Pane", func(t *testing.T) {
		mockHTML := "<html><body><div>No Calendar here</div></body></html>"
		httpClient := &http.Client{
			Transport: &mockRoundTripper{
				roundTrip: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(mockHTML)),
					}, nil
				},
			},
		}
		s := &scraper{client: httpClient}
		_, _, err := s.ScrapeCalendar(context.Background(), CalendarStandard)
		if err == nil {
			t.Error("expected error when calendar pane is missing")
		}
	})

	t.Run("Missing Table inside Pane", func(t *testing.T) {
		mockHTML := `<html><body><div id="lqd-tab-scedule-item-2"><h2>No Table</h2></div></body></html>`
		httpClient := &http.Client{
			Transport: &mockRoundTripper{
				roundTrip: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(mockHTML)),
					}, nil
				},
			},
		}
		s := &scraper{client: httpClient}
		_, _, err := s.ScrapeCalendar(context.Background(), CalendarStandard)
		if err == nil {
			t.Error("expected error when table is missing inside calendar pane")
		}
	})
}
