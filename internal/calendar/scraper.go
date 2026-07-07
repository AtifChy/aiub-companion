package calendar

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

const calendarURL = "https://www.aiub.edu/academic-calendar"

var calendarPaneIDs = map[CalendarType]string{
	CalendarStandard:  "lqd-tab-scedule-item-2",
	CalendarLLBBPharm: "lqd-tab-scedule-item-1",
}

type scraper struct {
	client *http.Client
}

// ScrapeCalendar fetches the academic calendar HTML content for the specified CalendarType.
func (s *scraper) ScrapeCalendar(ctx context.Context, calType CalendarType) (table *html.Node, semester string, err error) {
	paneID, ok := calendarPaneIDs[calType]
	if !ok {
		return nil, "", fmt.Errorf("unsupported calendar type: %s", calType)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, calendarURL, nil)
	if err != nil {
		return nil, "", fmt.Errorf("new request: %w", err)
	}
	req.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148.0.0.0 Safari/537.36",
	)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("fetching calendar: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("parse HTML: %w", err)
	}

	pane := findNodeByID(doc, paneID)
	if pane == nil {
		return nil, "", fmt.Errorf("calendar pane %q not found", paneID)
	}

	semester = extractSemester(pane)

	table = findNode(pane, "table")
	if table == nil {
		return nil, "", fmt.Errorf("no table found inside pane %q", paneID)
	}

	return table, semester, nil
}

func findNodeByID(n *html.Node, id string) *html.Node {
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == "id" && attr.Val == id {
				return n
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := findNodeByID(c, id); result != nil {
			return result
		}
	}
	return nil
}

func extractSemester(n *html.Node) string {
	if h2 := findNode(n, "h2"); h2 != nil {
		s := extractText(h2)
		if strings.TrimSpace(strings.ReplaceAll(s, "\u00a0", " ")) != "" {
			return strings.ToLower(s)
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := extractSemester(c); result != "" {
			return result
		}
	}
	return ""
}
