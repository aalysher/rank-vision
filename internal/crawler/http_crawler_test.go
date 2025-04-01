package crawler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPCrawler_CrawlPage(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Test Page</title>
			<meta name="description" content="Test Description">
		</head>
		<body>
			<h1>Hello World</h1>
			<p>This is a test page with some text.</p>
			<a href="https://example.com">Example Link</a>
		</body>
		</html>
		`
		w.Write([]byte(html))
	}))
	defer server.Close()

	// Создаем краулер
	crawler := NewHTTPCrawler(DefaultConfig())

	// Тестируем краулинг страницы
	result, err := crawler.CrawlPage(context.Background(), server.URL)
	if err != nil {
		t.Fatalf("CrawlPage failed: %v", err)
	}

	// Проверяем результаты
	if result.URL != server.URL {
		t.Errorf("Expected URL %s, got %s", server.URL, result.URL)
	}
	if result.Title != "Test Page" {
		t.Errorf("Expected title 'Test Page', got '%s'", result.Title)
	}
	if result.MetaDescription != "Test Description" {
		t.Errorf("Expected description 'Test Description', got '%s'", result.MetaDescription)
	}
	if result.WordCount < 10 {
		t.Errorf("Expected word count > 10, got %d", result.WordCount)
	}
	if len(result.Links) != 1 {
		t.Errorf("Expected 1 link, got %d", len(result.Links))
	}
	if result.Links[0] != "https://example.com" {
		t.Errorf("Expected link 'https://example.com', got '%s'", result.Links[0])
	}
}

func TestHTTPCrawler_IsValidURL(t *testing.T) {
	crawler := NewHTTPCrawler(DefaultConfig())

	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{"valid http", "http://example.com", true},
		{"valid https", "https://example.com", true},
		{"valid with path", "https://example.com/path", true},
		{"invalid scheme", "ftp://example.com", true},
		{"invalid format", "not a url", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := crawler.IsValidURL(tt.url)
			if result != tt.expected {
				t.Errorf("IsValidURL(%s) = %v; want %v", tt.url, result, tt.expected)
			}
		})
	}
}
