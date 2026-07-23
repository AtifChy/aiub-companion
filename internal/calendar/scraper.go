package calendar

import (
	"context"
	"fmt"
	"strings"

	"aiub-companion/internal/fetcher"

	"golang.org/x/net/html"
)

const calendarURL = "https://www.aiub.edu/academic-calendar"

var calendarPaneIDs = map[CalendarType]string{
	CalendarStandard:  "lqd-tab-scedule-item-2",
	CalendarLLBBPharm: "lqd-tab-scedule-item-1",
}

type scraper struct {
	fetcher *fetcher.Fetcher
}

// ScrapeCalendar fetches the academic calendar HTML content for the specified CalendarType.
func (s *scraper) ScrapeCalendar(ctx context.Context, calType CalendarType) (table *html.Node, semester string, err error) {
	paneID, ok := calendarPaneIDs[calType]
	if !ok {
		return nil, "", fmt.Errorf("unsupported calendar type: %s", calType)
	}

	doc, err := s.fetcher.FetchHTML(ctx, calendarURL)
	if err != nil {
		return nil, "", fmt.Errorf("parse HTML: %w", err)
	}

	pane := fetcher.FindNodeByID(doc, paneID)
	if pane == nil {
		return nil, "", fmt.Errorf("calendar pane %q not found", paneID)
	}

	semester = extractSemester(pane)

	table = fetcher.FindNodeByTag(pane, "table")
	if table == nil {
		return nil, "", fmt.Errorf("no table found inside pane %q", paneID)
	}

	return table, semester, nil
}

func extractSemester(n *html.Node) string {
	if h2 := fetcher.FindNodeByTag(n, "h2"); h2 != nil {
		s := fetcher.GetInnerText(h2)
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
