package crawler

import (
	"context"
	"time"
)

// CrawlerResult представляет результат краулинга страницы
type CrawlerResult struct {
	URL             string
	Title           string
	MetaDescription string
	WordCount       int
	Links           []string
	Error           error
}

// Crawler определяет интерфейс для краулера
type Crawler interface {
	// CrawlPage краулит одну страницу и возвращает результат
	CrawlPage(ctx context.Context, targetURL string) (*CrawlerResult, error)

	// IsValidURL проверяет, является ли URL валидным для краулинга
	IsValidURL(targetURL string) bool
}

// Config содержит конфигурацию краулера
type Config struct {
	UserAgent    string
	RequestDelay time.Duration
	MaxRetries   int
	Timeout      time.Duration
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() *Config {
	return &Config{
		UserAgent:    "RankVision Bot/1.0",
		RequestDelay: time.Second,
		MaxRetries:   3,
		Timeout:      30 * time.Second,
	}
}
