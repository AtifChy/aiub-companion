package notice

import (
	"context"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"aiub-companion/internal/fetcher"

	"github.com/microcosm-cc/bluemonday"
	"github.com/wailsapp/wails/v3/pkg/application"
	"golang.org/x/net/html"
)

const defaultBaseURL = "https://www.aiub.edu"

type scraper struct {
	fetcher *fetcher.Fetcher
	baseURL string
}

func (s *scraper) ScrapeNotices(ctx context.Context, count int) ([]Notice, error) {
	url := s.baseURL + "/category/notices?pageNo=1&pageSize=" + strconv.Itoa(count)

	doc, err := s.fetcher.FetchHTML(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("fetch notices: %w", err)
	}

	cards := fetcher.FindNodesByClass(doc, "div", "notification")

	notices := make([]Notice, 0, len(cards))

	for _, card := range cards {
		fields := extractCardFields(card)
		if fields.slug == "" {
			continue
		}

		notices = append(notices, Notice{
			ID:         fields.slug,
			Title:      fields.title,
			Summary:    fields.summary,
			PostedDate: parseDate(fields.slug, fields.date),
			Category:   string(matchCategory(fields.title)),
			IsUrgent:   isUrgent(fields.title),
		})
	}

	return notices, nil
}

func extractCardFields(card *html.Node) struct{ slug, title, summary, date string } {
	var f struct{ slug, title, summary, date string }
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			cls := fetcher.GetAttr(n, "class")
			switch {
			case f.slug == "" && n.Data == "a":
				href := fetcher.GetAttr(n, "href")
				if href != "" && href != "#" && !strings.HasPrefix(href, "javascript:") {
					f.slug = strings.TrimPrefix(href, "/")
				}
			case n.Data == "h2" && strings.Contains(cls, "title"):
				f.title = fetcher.GetInnerText(n)
			case n.Data == "p" && strings.Contains(cls, "desc"):
				f.summary = fetcher.GetInnerText(n)
			case n.Data == "div" && strings.Contains(cls, "date-custom"):
				f.date = fetcher.GetInnerText(n)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(card)
	return f
}

func parseDate(slug, raw string) string {
	text := strings.Join(strings.Fields(raw), " ")
	date, err := time.Parse("2 Jan 2006", text)
	if err != nil {
		application.Get().Logger.Warn("Failed to parse date, defaulting to current date", "slug", slug, "raw", raw, "error", err)
		return time.Now().Format(time.DateOnly)
	}
	return date.Format(time.DateOnly)
}

func isUrgent(title string) bool {
	words := strings.Fields(strings.TrimSpace(title))

	filtered := make([]string, 0, len(words))
	for i := range words {
		if len(words[i]) > 1 {
			filtered = append(filtered, words[i])
		}
	}

	if len(filtered) < 3 {
		return false
	}

	upperCount := 0
	for i := range filtered {
		if filtered[i] == strings.ToUpper(filtered[i]) {
			upperCount++
		}
	}

	return float64(upperCount)/float64(len(filtered)) >= 0.5
}

func (s *scraper) ScrapeNoticeDetails(ctx context.Context, slug string) (NoticeDetails, error) {
	url := s.baseURL + "/" + slug

	doc, err := s.fetcher.FetchHTML(ctx, url)
	if err != nil {
		return NoticeDetails{}, fmt.Errorf("fetch notice details: %w", err)
	}

	contents := fetcher.FindNodesByClass(doc, "div", "question-column")
	if len(contents) < 2 {
		return NoticeDetails{}, errors.New("invalid page structure")
	}

	title := fetcher.GetInnerText(contents[0])
	attachments := s.extractAttachments(contents[1])
	s.removeAttachments(contents[1])
	body := s.extractBodyHTML(contents[1])

	details := NoticeDetails{
		FullTitle:   title,
		Content:     body,
		Attachments: attachments,
	}

	return details, nil
}

func (s *scraper) extractBodyHTML(node *html.Node) string {
	var sb strings.Builder
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if fetcher.GetInnerText(c) != "" {
			_ = html.Render(&sb, c)
		}
	}
	return bluemonday.UGCPolicy().Sanitize(sb.String())
}

func (s *scraper) extractAttachments(node *html.Node) []Attachment {
	var attachments []Attachment

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if s.isAttachment(n) {
			href := fetcher.GetAttr(n, "href")
			clean := strings.ReplaceAll(href, "\\", "/")
			if !strings.HasPrefix(clean, "http") {
				clean = s.baseURL + clean
			}

			base := path.Base(clean)
			id := strings.TrimSuffix(base, filepath.Ext(base))

			label := strings.TrimSpace(fetcher.GetInnerText(n))
			if label == "" {
				label = "Document"
			}

			attachments = append(attachments, Attachment{
				ID:    id,
				Label: label,
				URL:   clean,
			})
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(node)
	return attachments
}

func (s *scraper) removeAttachments(node *html.Node) {
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		var next *html.Node
		for c := n.FirstChild; c != nil; c = next {
			next = c.NextSibling
			if s.isAttachment(c) {
				n.RemoveChild(c)
			} else {
				walk(c)
			}
		}
	}
	walk(node)
}

func (s *scraper) isAttachment(n *html.Node) bool {
	if n == nil || n.Type != html.ElementNode || n.Data != "a" {
		return false
	}
	href := fetcher.GetAttr(n, "href")
	clean := strings.ReplaceAll(href, "\\", "/")
	lower := strings.ToLower(clean)
	return strings.Contains(lower, "uploads")
}
