package crawler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSiteCrawler_CrawlSite(t *testing.T) {
	// Создаем тестовый сервер с несколькими страницами
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			html := `
			<!DOCTYPE html>
			<html>
			<head>
				<title>Home Page</title>
				<meta name="description" content="Home page description">
			</head>
			<body>
				<h1>Welcome</h1>
				<a href="/about">About</a>
				<a href="/contact">Contact</a>
			</body>
			</html>
			`
			w.Write([]byte(html))
		case "/about":
			html := `
			<!DOCTYPE html>
			<html>
			<head>
				<title>About Page</title>
				<meta name="description" content="About page description">
			</head>
			<body>
				<h1>About Us</h1>
				<a href="/">Home</a>
				<a href="/contact">Contact</a>
			</body>
			</html>
			`
			w.Write([]byte(html))
		case "/contact":
			html := `
			<!DOCTYPE html>
			<html>
			<head>
				<title>Contact Page</title>
				<meta name="description" content="Contact page description">
			</head>
			<body>
				<h1>Contact Us</h1>
				<a href="/">Home</a>
				<a href="/about">About</a>
			</body>
			</html>
			`
			w.Write([]byte(html))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	// Создаем краулер
	crawler, err := NewSiteCrawler(server.URL, DefaultConfig(), 10, 3)
	if err != nil {
		t.Fatalf("Failed to create crawler: %v", err)
	}

	// Запускаем краулинг
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := crawler.CrawlSite(ctx)
	if err != nil {
		t.Fatalf("CrawlSite failed: %v", err)
	}

	// Проверяем результаты
	if result.BaseURL != server.URL {
		t.Errorf("Expected base URL %s, got %s", server.URL, result.BaseURL)
	}
	if result.TotalPages != 3 {
		t.Errorf("Expected 3 pages, got %d", result.TotalPages)
	}
	if result.Statistics.TotalWordCount < 10 {
		t.Errorf("Expected word count > 10, got %d", result.Statistics.TotalWordCount)
	}
	if result.Statistics.AverageWordCount < 3 {
		t.Errorf("Expected average word count > 3, got %d", result.Statistics.AverageWordCount)
	}

	// Проверяем наличие всех страниц
	expectedPages := []string{"/", "/about", "/contact"}
	for _, page := range expectedPages {
		fullURL := server.URL + page
		if _, exists := result.Pages[fullURL]; !exists {
			t.Errorf("Page %s not found in results", fullURL)
		}
	}
}

func TestSiteCrawler_IsValidInternalURL(t *testing.T) {
	crawler, err := NewSiteCrawler("https://example.com", DefaultConfig(), 10, 3)
	if err != nil {
		t.Fatalf("Failed to create crawler: %v", err)
	}

	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{"same domain", "https://example.com/page", true},
		{"same domain with path", "https://example.com/about/contact", true},
		{"different domain", "https://other.com/page", false},
		{"file", "https://example.com/file.pdf", false},
		{"image", "https://example.com/image.jpg", false},
		{"invalid url", "not a url", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := crawler.isValidInternalURL(tt.url)
			if result != tt.expected {
				t.Errorf("isValidInternalURL(%s) = %v; want %v", tt.url, result, tt.expected)
			}
		})
	}
}
