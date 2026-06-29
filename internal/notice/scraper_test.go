package notice

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/net/html"
)

func okHandler(html string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(html))
	}
}

func newTestScraper(t *testing.T, handler http.Handler) (*scraper, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(handler)
	s := &scraper{
		client:  srv.Client(),
		baseURL: srv.URL,
	}
	return s, srv
}

func parseHTML(t *testing.T, s string) *html.Node {
	t.Helper()
	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		t.Fatalf("failed to parse HTML: %v", err)
	}
	return doc
}

// fetchHTML tests

func TestFetchHTML_RetryOn5xxFailure(t *testing.T) {
	var attempts atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		w.WriteHeader(http.StatusServiceUnavailable) // 503 Service Unavailable (retriable)
	}))
	defer srv.Close()

	s := &scraper{
		client:  srv.Client(),
		baseURL: srv.URL,
	}

	_, err := s.fetchHTML(context.Background(), srv.URL)
	if err == nil {
		t.Fatal("expected error from repeated 503 responses, got nil")
	}

	expectedAttempts := int32(5) // maxAttempts configured in fetchHTML
	actualAttempts := attempts.Load()
	if actualAttempts != expectedAttempts {
		t.Errorf("expected %d fetch attempts, got %d", expectedAttempts, actualAttempts)
	}
}

func TestFetchHTML_RetryAndRecover(t *testing.T) {
	var attempts atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currentAttempt := attempts.Add(1)
		if currentAttempt == 1 {
			w.WriteHeader(http.StatusInternalServerError) // 500 failure first
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("<html><body><p>success</p></body></html>"))
	}))
	defer srv.Close()

	s := &scraper{
		client:  srv.Client(),
		baseURL: srv.URL,
	}

	doc, err := s.fetchHTML(context.Background(), srv.URL)
	if err != nil {
		t.Fatalf("expected recovery and success on second attempt, got error: %v", err)
	}
	if doc == nil {
		t.Fatal("expected parsed document, got nil")
	}

	actualAttempts := attempts.Load()
	if actualAttempts != 2 {
		t.Errorf("expected recovery on attempt 2, but took %d attempts", actualAttempts)
	}
}

func TestFetchHTML_NoRetryOn404(t *testing.T) {
	var attempts atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		w.WriteHeader(http.StatusNotFound) // 404 Not Found
	}))
	defer srv.Close()

	s := &scraper{
		client:  srv.Client(),
		baseURL: srv.URL,
	}

	_, err := s.fetchHTML(context.Background(), srv.URL)
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}

	actualAttempts := attempts.Load()
	if actualAttempts != 1 {
		t.Errorf("expected exactly 1 attempt for non-retriable 404, got %d", actualAttempts)
	}
}

func TestFetchHTML_ContextCancellation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable) // 503 to trigger retry logic
	}))
	defer srv.Close()

	s := &scraper{
		client:  srv.Client(),
		baseURL: srv.URL,
	}

	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(1*time.Second, cancel) // Cancel after 1 second

	_, err := s.fetchHTML(ctx, srv.URL)
	if err == nil {
		t.Fatal("expected error due to context cancellation, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected error chain to contain `context.Canceled`, got %v", err)
	}
}

func TestFetchHTML_OK(t *testing.T) {
	s, srv := newTestScraper(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<html><body><p>hello</p></body></html>`))
	}))
	defer srv.Close()

	doc, err := s.fetchHTML(context.Background(), srv.URL)
	if err != nil {
		t.Fatalf("failed to fetch HTML: %v", err)
	}
	if doc == nil {
		t.Fatal("expected non-nil document")
	}
}

// findNodesByClass tests

func TestFindNodesByClass_MatchesCorrectNodes(t *testing.T) {
	doc := parseHTML(t, `<html><body>
		<div class="notification active">First</div>
		<div class="notification">Second</div>
		<div class="other">Third</div>
	</body></html>`)

	nodes := findNodesByClass(doc, "div", "notification")
	if len(nodes) != 2 {
		t.Errorf("expected 2 nodes, got %d", len(nodes))
	}
}

func TestFindNodesByClass_NoMatch(t *testing.T) {
	doc := parseHTML(t, `<html><body><div class="other">Text</div></body></html>`)

	nodes := findNodesByClass(doc, "div", "notification")
	if len(nodes) != 0 {
		t.Errorf("expected 0 nodes, got %d", len(nodes))
	}
}

func TestFindNodesByClass_WrongTag(t *testing.T) {
	doc := parseHTML(t, `<html><body><p class="notification">Text</p></body></html>`)

	nodes := findNodesByClass(doc, "div", "notification")
	if len(nodes) != 0 {
		t.Errorf("expected 0 nodes for wrong tag, got %d", len(nodes))
	}
}

func TestFindNodesByClass_PartialClassNameNoMatch(t *testing.T) {
	// "notif" should NOT match "notification"
	doc := parseHTML(t, `<html><body><div class="notification">Text</div></body></html>`)

	nodes := findNodesByClass(doc, "div", "notif")
	if len(nodes) != 0 {
		t.Errorf("partial class name should not match, got %d nodes", len(nodes))
	}
}

// ScrapeNotices tests

const noticesListHTML = `<html><body>
	<div class="notification">
		<a href="/notices/exam-schedule-2025">link</a>
		<h2 class="title">Final Exam Schedule 2025</h2>
		<p class="desc">Check your exam schedule here.</p>
		<div class="date-custom">12 May 2025</div>
	</div>
	<div class="notification">
		<a href="/notices/eid-holiday">link</a>
		<h2 class="title">Eid Holiday Notice</h2>
		<p class="desc">University will be closed for Eid.</p>
		<div class="date-custom">1 Apr 2025</div>
	</div>
</body></html>`

func TestScrapeNotices_ReturnsCorrectNotices(t *testing.T) {
	s, srv := newTestScraper(t, okHandler(noticesListHTML))
	defer srv.Close()

	notices, err := s.ScrapeNotices(context.Background(), 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(notices) != 2 {
		t.Fatalf("expected 2 notices, got %d", len(notices))
	}

	if notices[0].ID != "notices/exam-schedule-2025" {
		t.Errorf("expected slug 'notices/exam-schedule-2025', got %q", notices[0].ID)
	}
	if notices[0].Title != "Final Exam Schedule 2025" {
		t.Errorf("expected title 'Final Exam Schedule 2025', got %q", notices[0].Title)
	}
	if notices[0].Category != "exam" {
		t.Errorf("expected category 'exam', got %q", notices[0].Category)
	}
	if notices[0].PostedDate != "2025-05-12" {
		t.Errorf("expected posted_date '2025-05-12', got %q", notices[0].PostedDate)
	}
	if notices[0].Summary != "Check your exam schedule here." {
		t.Errorf("expected summary, got %q", notices[0].Summary)
	}
}

func TestScrapeNotices_SkipsCardWithNoSlug(t *testing.T) {
	s, srv := newTestScraper(t, okHandler(`<html><body>
		<div class="notification">
			<h2 class="title">Orphan Notice</h2>
		</div>
		<div class="notification">
			<a href="/notices/valid">link</a>
			<h2 class="title">Valid Notice</h2>
			<div class="date-custom">1 Jan 2025</div>
		</div>
	</body></html>`))
	defer srv.Close()

	notices, err := s.ScrapeNotices(context.Background(), 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(notices) != 1 {
		t.Fatalf("expected 1 notice, got %d", len(notices))
	}
}

func TestScrapeNotices_EmptyPage(t *testing.T) {
	s, srv := newTestScraper(t, okHandler(`<html><body></body></html>`))
	defer srv.Close()

	notices, err := s.ScrapeNotices(context.Background(), 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(notices) != 0 {
		t.Errorf("expected 0 notices, got %d", len(notices))
	}
}

func TestScrapeNotices_FetchError(t *testing.T) {
	s, srv := newTestScraper(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	_, err := s.ScrapeNotices(context.Background(), 20)
	if err == nil {
		t.Fatal("expected error for 503, got nil")
	}
}

// ScrapeNoticeDetails tests

const noticePageHTML = `<html><body>
	<div class="question-column">Notice Title</div>
	<div class="question-column">
		<p>This is the first paragraph of the notice with enough text.</p>
		<p>Login to the portal for more information.</p>
		<p>Second paragraph with sufficient content here for testing.</p>
		<ul>
			<li>Important list item with sufficient length here.</li>
		</ul>
		<a href="/files/uploads/notice-2025.pdf">Download Notice PDF</a>
		<a href="/files/uploads/results-2025.pdf">Download Results PDF</a>
		<a href="https://external.com/page">External Link</a>
	</div>
</body></html>`

func TestScrapeNoticeDetails_Content(t *testing.T) {
	s, srv := newTestScraper(t, okHandler(noticePageHTML))
	defer srv.Close()

	details, err := s.ScrapeNoticeDetails(context.Background(), "notices/test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if details.FullTitle != "Notice Title" {
		t.Errorf("expected title 'Notice Title', got: %s", details.FullTitle)
	}
	if details.Content == "" {
		t.Fatal("expected non-empty content")
	}
	if !strings.Contains(details.Content, "first paragraph") {
		t.Errorf("expected content to include first paragraph, got: %s", details.Content)
	}
	if strings.Contains(details.Content, "notice-2025.pdf") || strings.Contains(details.Content, "Download Notice PDF") {
		t.Errorf("content body should not contain attachment links/labels, got: %s", details.Content)
	}
}

func TestScrapeNoticeDetails_Attachments(t *testing.T) {
	s, srv := newTestScraper(t, okHandler(noticePageHTML))
	defer srv.Close()

	details, err := s.ScrapeNoticeDetails(context.Background(), "notices/test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(details.Attachments) != 2 {
		t.Errorf("expected 2 attachments, got %d", len(details.Attachments))
	}
	for _, att := range details.Attachments {
		if att.Label == "" {
			t.Errorf("attachment label should not be empty")
		}
		if att.URL == "" {
			t.Errorf("attachment url should not be empty")
		}
		if !strings.HasSuffix(strings.ToLower(att.URL), ".pdf") {
			t.Errorf("attachment URL should end in .pdf, got: %s", att.URL)
		}
	}
}

func TestScrapeNoticeDetails_EmptyLabelFallback(t *testing.T) {
	page := `<html><body>
	<div class="question-column">Title</div>
	<div class="question-column">
		<a href="/files/uploads/doc.pdf"></a>
	</div>
</body></html>`
	s, srv := newTestScraper(t, okHandler(page))
	defer srv.Close()

	details, err := s.ScrapeNoticeDetails(context.Background(), "notices/test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(details.Attachments) != 1 {
		t.Errorf("expected 1 attachment, got %d", len(details.Attachments))
	}
	if details.Attachments[0].Label != "Document" {
		t.Errorf("attachment label should be 'Document', got: %s", details.Attachments[0].Label)
	}
}

func TestScrapeNoticeDetails_NoAttachments(t *testing.T) {
	page := `<html><body>
	<div class="question-column">Title</div>
	<div class="question-column">
		<p>Just a plain notice with no attachments.</p>
	</div>
</body></html>`
	s, srv := newTestScraper(t, okHandler(page))
	defer srv.Close()

	details, err := s.ScrapeNoticeDetails(context.Background(), "notices/plain")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(details.Attachments) != 0 {
		t.Errorf("expected 0 attachment, got %d", len(details.Attachments))
	}
}

func TestScrapeNoticeDetails_FetchError(t *testing.T) {
	s, srv := newTestScraper(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	_, err := s.ScrapeNoticeDetails(context.Background(), "notices/bad")
	if err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}
}

// getInnerText tests
func TestGetInnerText_NestedElements(t *testing.T) {
	doc := parseHTML(t, `<html><body><p id="target">  Hello <span>World</span>  </p></body></html>`)

	var pNode *html.Node
	var find func(*html.Node)
	find = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "p" {
			pNode = n
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			find(c)
		}
	}
	find(doc)

	if pNode == nil {
		t.Fatal("could not find <p> node")
	}
	text := getInnerText(pNode)
	if text != "Hello World" {
		t.Errorf("expected 'Hello World', got %q", text)
	}
}

func TestGetInnerText_EmptyNode(t *testing.T) {
	doc := parseHTML(t, `<html><body><p></p></body></html>`)

	var pNode *html.Node
	var find func(*html.Node)
	find = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "p" {
			pNode = n
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			find(c)
		}
	}
	find(doc)

	text := getInnerText(pNode)
	if text != "" {
		t.Errorf("expected empty string, got %q", text)
	}
}

// determineCategory tests

func TestDetermineCategory(t *testing.T) {
	tests := []struct {
		title    string
		expected string
	}{
		{"Final Exam Schedule Released", "exam"},
		{"Seat Plan for Mid-Term", "exam"},
		{"Viva Schedule Announced", "general"},
		{"Probation List Updated", "general"},
		{"Programming Contest Registration", "registration"},
		{"Job Fair at AIUB Campus", "general"},
		{"Sports Day Announced", "general"},
		{"Holiday Notice for Eid", "holiday"},
		{"University Closed Due to Cyclone", "general"},
		{"Weather Advisory Issued", "general"},
		{"Course Registration Guidelines", "registration"},
		{"General Announcement", "general"},
		{"", "general"},
		{"FINAL EXAM UPPERCASE", "general"},
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			got := string(matchCategory(tc.title))
			if got != tc.expected {
				t.Errorf("matchCategory(%q) = %q, want %q", tc.title, got, tc.expected)
			}
		})
	}
}
