package notice

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/wailsapp/wails/v3/pkg/application"
	"golang.org/x/net/html"
)

const defaultBaseURL = "https://www.aiub.edu"

type scraper struct {
	client  *http.Client
	baseURL string
}

func (s *scraper) ScrapeNotices(ctx context.Context, count int) ([]Notice, error) {
	url := s.baseURL + "/category/notices?pageNo=1&pageSize=" + strconv.Itoa(count)

	doc, err := s.fetchHTML(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch notices: %w", err)
	}

	cards := findNodesByClass(doc, "div", "notification")

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
			cls := getAttr(n, "class")
			switch {
			case f.slug == "" && n.Data == "a":
				href := getAttr(n, "href")
				if href != "" && href != "#" && !strings.HasPrefix(href, "javascript:") {
					f.slug = strings.TrimPrefix(href, "/")
				}
			case n.Data == "h2" && strings.Contains(cls, "title"):
				f.title = getInnerText(n)
			case n.Data == "p" && strings.Contains(cls, "desc"):
				f.summary = getInnerText(n)
			case n.Data == "div" && strings.Contains(cls, "date-custom"):
				f.date = getInnerText(n)
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
		return time.Now().Format("2006-01-02")
	}
	return date.Format("2006-01-02")
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

	doc, err := s.fetchHTML(ctx, url)
	if err != nil {
		return NoticeDetails{}, fmt.Errorf("failed to fetch notice details: %w", err)
	}

	contents := findNodesByClass(doc, "div", "question-column")
	if len(contents) < 2 {
		return NoticeDetails{}, fmt.Errorf("invalid notice page structure: expected title and body columns")
	}

	title := getInnerText(contents[0])
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
		if getInnerText(c) != "" {
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
			href := getAttr(n, "href")
			clean := strings.ReplaceAll(href, "\\", "/")
			if !strings.HasPrefix(clean, "http") {
				clean = s.baseURL + clean
			}

			base := path.Base(clean)
			id := strings.TrimSuffix(base, filepath.Ext(base))

			label := strings.TrimSpace(getInnerText(n))
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
	href := getAttr(n, "href")
	clean := strings.ReplaceAll(href, "\\", "/")
	lower := strings.ToLower(clean)
	return strings.Contains(lower, "uploads")
}

func (s *scraper) fetchHTML(ctx context.Context, url string) (*html.Node, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148.0.0.0 Safari/537.36")

	resp, err := s.client.Do(req)
	if err != nil {
		if netErr, ok := errors.AsType[net.Error](err); ok && netErr.Timeout() {
			return nil, fmt.Errorf("timeout fetching %s: %w", url, err)
		}
		if dnsErr, ok := errors.AsType[*net.DNSError](err); ok {
			if dnsErr.IsNotFound {
				return nil, fmt.Errorf("no internet or host found: %w", dnsErr)
			}
			return nil, fmt.Errorf("DNS error fetching %s: %w", url, dnsErr)
		}
		if opErr, ok := errors.AsType[*net.OpError](err); ok {
			if sysErr, ok := errors.AsType[syscall.Errno](err); ok {
				switch sysErr {
				case syscall.ECONNREFUSED:
					return nil, fmt.Errorf("connection refused: %w", opErr)
				case syscall.ECONNRESET:
					return nil, fmt.Errorf("connection reset: %w", opErr)
				case syscall.ENETUNREACH:
					return nil, fmt.Errorf("network unreachable: %w", opErr)
				}
			}
			return nil, fmt.Errorf("operation error fetching %s: %w", url, opErr)
		}
		return nil, fmt.Errorf("failed to fetch %s: %w", url, err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			application.Get().Logger.Warn("Failed to close response body", "url", url, "error", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d for %s", resp.StatusCode, url)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}
	return doc, nil
}

func findNodesByClass(n *html.Node, tagName, className string) []*html.Node {
	var result []*html.Node
	if n.Type == html.ElementNode && (tagName == "" || n.Data == tagName) {
		for _, a := range n.Attr {
			if a.Key != "class" {
				continue
			}
			classes := strings.FieldsSeq(a.Val)
			for c := range classes {
				if c == className {
					result = append(result, n)
					break
				}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result = append(result, findNodesByClass(c, tagName, className)...)
	}
	return result
}

func getAttr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

func getInnerText(n *html.Node) string {
	var sb strings.Builder
	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.TextNode {
			sb.WriteString(node.Data)
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return strings.TrimSpace(sb.String())
}
