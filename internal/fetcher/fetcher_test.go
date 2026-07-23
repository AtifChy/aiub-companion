package fetcher

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

func parseHTML(t *testing.T, s string) *html.Node {
	t.Helper()
	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		t.Fatalf("failed to parse HTML: %v", err)
	}
	return doc
}

// Fetcher tests

func TestFetcher_FetchHTML_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") != DefaultUserAgent {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("<html><body><h1>Hello</h1></body></html>"))
	}))
	defer srv.Close()

	f := New(WithHTTPClient(srv.Client()))
	doc, err := f.FetchHTML(context.Background(), srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if doc == nil {
		t.Fatal("expected non-nil html node")
	}

	h1 := FindNodeByTag(doc, "h1")
	if h1 == nil || GetInnerText(h1) != "Hello" {
		t.Errorf("expected h1 text 'Hello', got %v", h1)
	}
}

func TestFetcher_CustomUserAgent(t *testing.T) {
	customUA := "CustomUserAgent/1.0"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") != customUA {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("<html><body>ok</body></html>"))
	}))
	defer srv.Close()

	f := New(
		WithHTTPClient(srv.Client()),
		WithUserAgent(customUA),
	)
	_, err := f.FetchHTML(context.Background(), srv.URL)
	if err != nil {
		t.Fatalf("unexpected error with custom UA: %v", err)
	}
}

func TestFetcher_FetchHTML_RetryOn5xx(t *testing.T) {
	var attempts atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	f := New(
		WithHTTPClient(srv.Client()),
		WithMaxAttempts(3),
	)
	f.initialBackoff = 1 * time.Millisecond // Speed up test execution

	_, err := f.FetchHTML(context.Background(), srv.URL)
	if err == nil {
		t.Fatal("expected error after retries failed")
	}

	if attempts.Load() != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts.Load())
	}
}

func TestFetcher_FetchHTML_RetryAndRecover(t *testing.T) {
	var attempts atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if attempts.Add(1) == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("<html><body>success</body></html>"))
	}))
	defer srv.Close()

	f := New(
		WithHTTPClient(srv.Client()),
		WithMaxAttempts(3),
	)
	f.initialBackoff = 1 * time.Millisecond

	doc, err := f.FetchHTML(context.Background(), srv.URL)
	if err != nil {
		t.Fatalf("expected recovery on second attempt, got error: %v", err)
	}
	if doc == nil {
		t.Fatal("expected non-nil document")
	}
	if attempts.Load() != 2 {
		t.Errorf("expected 2 attempts, got %d", attempts.Load())
	}
}

func TestFetcher_FetchHTML_NoRetryOn404(t *testing.T) {
	var attempts atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	f := New(
		WithHTTPClient(srv.Client()),
		WithMaxAttempts(5),
	)
	f.initialBackoff = 1 * time.Millisecond

	_, err := f.FetchHTML(context.Background(), srv.URL)
	if err == nil {
		t.Fatal("expected error on 404")
	}

	if attempts.Load() != 1 {
		t.Errorf("expected exactly 1 attempt for non-retriable 404, got %d", attempts.Load())
	}
}

func TestFetcher_FetchHTML_ContextCanceled(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	f := New(
		WithHTTPClient(srv.Client()),
		WithMaxAttempts(5),
	)
	f.initialBackoff = 100 * time.Millisecond

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	_, err := f.FetchHTML(ctx, srv.URL)
	if err == nil {
		t.Fatal("expected error due to canceled context")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled error, got %v", err)
	}
}

// DOM Utility tests

func TestGetAttr(t *testing.T) {
	doc := parseHTML(t, `<div id="my-id" class="container main">Text</div>`)
	div := FindNodeByTag(doc, "div")

	if val := GetAttr(div, "id"); val != "my-id" {
		t.Errorf("expected 'my-id', got %q", val)
	}
	if val := GetAttr(div, "class"); val != "container main" {
		t.Errorf("expected 'container main', got %q", val)
	}
	if val := GetAttr(div, "nonexistent"); val != "" {
		t.Errorf("expected empty string for non-existent attribute, got %q", val)
	}
	if val := GetAttr(nil, "id"); val != "" {
		t.Errorf("expected empty string for nil node, got %q", val)
	}
}

func TestGetAttrInt(t *testing.T) {
	doc := parseHTML(t, `
		<table>
			<tr>
				<td colspan="3" rowspan="invalid">Cell</td>
			</tr>
		</table>
	`)
	td := FindNodeByTag(doc, "td")

	if val := GetAttrInt(td, "colspan", 1); val != 3 {
		t.Errorf("expected 3, got %d", val)
	}
	if val := GetAttrInt(td, "rowspan", 1); val != 1 {
		t.Errorf("expected default 1 for invalid int, got %d", val)
	}
	if val := GetAttrInt(td, "missing", 5); val != 5 {
		t.Errorf("expected default 5 for missing attr, got %d", val)
	}
	if val := GetAttrInt(nil, "colspan", 1); val != 1 {
		t.Errorf("expected default 1 for nil node, got %d", val)
	}
}

func TestGetInnerText(t *testing.T) {
	doc := parseHTML(t, `<div>  Hello <span>World</span> ! </div>`)
	div := FindNodeByTag(doc, "div")

	if text := GetInnerText(div); text != "Hello World !" {
		t.Errorf("expected 'Hello World !', got %q", text)
	}
	if text := GetInnerText(nil); text != "" {
		t.Errorf("expected empty string for nil node, got %q", text)
	}
}

func TestGetInnerHTML(t *testing.T) {
	doc := parseHTML(t, `<div class="box"><p>Paragraph <strong>bold</strong></p></div>`)
	div := FindNodeByTag(doc, "div")

	expected := `<p>Paragraph <strong>bold</strong></p>`
	if htmlStr := GetInnerHTML(div); htmlStr != expected {
		t.Errorf("expected %q, got %q", expected, htmlStr)
	}
	if htmlStr := GetInnerHTML(nil); htmlStr != "" {
		t.Errorf("expected empty string for nil node, got %q", htmlStr)
	}
}

func TestStripTags(t *testing.T) {
	input := `<p>Hello <a href="#">World</a>!</p>`
	expected := `Hello World!`
	if result := StripTags(input); result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestFindNodeByID(t *testing.T) {
	doc := parseHTML(t, `<div><p id="target">Found</p></div>`)

	node := FindNodeByID(doc, "target")
	if node == nil {
		t.Fatal("expected to find node with ID 'target'")
	}
	if GetInnerText(node) != "Found" {
		t.Errorf("expected text 'Found', got %q", GetInnerText(node))
	}

	if missing := FindNodeByID(doc, "not-found"); missing != nil {
		t.Errorf("expected nil for non-existent ID, got %v", missing)
	}
	if nilNode := FindNodeByID(nil, "id"); nilNode != nil {
		t.Errorf("expected nil for nil root node, got %v", nilNode)
	}
}

func TestFindNodesByClass(t *testing.T) {
	doc := parseHTML(t, `
		<div>
			<span class="item primary">1</span>
			<span class="item">2</span>
			<p class="item">3</p>
		</div>
	`)

	spans := FindNodesByClass(doc, "span", "item")
	if len(spans) != 2 {
		t.Errorf("expected 2 span elements with class 'item', got %d", len(spans))
	}

	all := FindNodesByClass(doc, "", "item")
	if len(all) != 3 {
		t.Errorf("expected 3 total elements with class 'item', got %d", len(all))
	}

	none := FindNodesByClass(doc, "div", "item")
	if len(none) != 0 {
		t.Errorf("expected 0 div elements with class 'item', got %d", len(none))
	}

	if nilNodes := FindNodesByClass(nil, "div", "item"); len(nilNodes) != 0 {
		t.Errorf("expected empty slice for nil node, got %d", len(nilNodes))
	}
}
