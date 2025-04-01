package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"rank-vision/internal/crawler"
)

const defaultUserAgent = "RankVision Bot/1.0 (+https://rank-vision.com; bot@rank-vision.com) Compatible/Go-http-client/1.1"

func main() {
	// Парсим аргументы командной строки
	url := flag.String("url", "", "URL для краулинга")
	timeout := flag.Duration("timeout", 30*time.Second, "Таймаут запроса")
	userAgent := flag.String("user-agent", defaultUserAgent, "User-Agent для запросов")
	maxPages := flag.Int("max-pages", 100, "Максимальное количество страниц для краулинга")
	maxDepth := flag.Int("max-depth", 3, "Максимальная глубина краулинга")
	singlePage := flag.Bool("single", false, "Краулить только одну страницу")
	flag.Parse()

	if *url == "" {
		log.Fatal("Необходимо указать URL с помощью флага -url")
	}

	// Создаем конфигурацию краулера
	config := &crawler.Config{
		UserAgent:    *userAgent,
		RequestDelay: time.Second,
		MaxRetries:   3,
		Timeout:      *timeout,
	}

	if *singlePage {
		// Краулим одну страницу
		c := crawler.NewHTTPCrawler(config)
		if !c.IsValidURL(*url) {
			log.Fatal("Невалидный URL")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *timeout)
		defer cancel()

		fmt.Printf("Начинаем краулинг %s...\n", *url)
		fmt.Printf("Используется User-Agent: %s\n", *userAgent)
		result, err := c.CrawlPage(ctx, *url)
		if err != nil {
			log.Fatalf("Ошибка при краулинге: %v", err)
		}

		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(result); err != nil {
			log.Fatalf("Ошибка при выводе результата: %v", err)
		}
	} else {
		// Краулим весь сайт
		sc, err := crawler.NewSiteCrawler(*url, config, *maxPages, *maxDepth)
		if err != nil {
			log.Fatalf("Ошибка при создании краулера: %v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), *timeout)
		defer cancel()

		fmt.Printf("Начинаем краулинг сайта %s...\n", *url)
		fmt.Printf("Используется User-Agent: %s\n", *userAgent)
		fmt.Printf("Максимальное количество страниц: %d\n", *maxPages)
		fmt.Printf("Максимальная глубина: %d\n", *maxDepth)

		result, err := sc.CrawlSite(ctx)
		if err != nil {
			log.Fatalf("Ошибка при краулинге сайта: %v", err)
		}

		// Выводим статистику
		fmt.Printf("\nСтатистика краулинга:\n")
		fmt.Printf("Всего страниц: %d\n", result.TotalPages)
		fmt.Printf("Всего слов: %d\n", result.Statistics.TotalWordCount)
		fmt.Printf("Среднее количество слов на страницу: %d\n", result.Statistics.AverageWordCount)
		fmt.Printf("Время выполнения: %v\n", result.EndTime.Sub(result.StartTime))

		if len(result.Statistics.MissingMetaDesc) > 0 {
			fmt.Printf("\nСтраницы без мета-описания (%d):\n", len(result.Statistics.MissingMetaDesc))
			for _, url := range result.Statistics.MissingMetaDesc {
				fmt.Printf("  - %s\n", url)
			}
		}

		if len(result.Statistics.MissingTitle) > 0 {
			fmt.Printf("\nСтраницы без заголовка (%d):\n", len(result.Statistics.MissingTitle))
			for _, url := range result.Statistics.MissingTitle {
				fmt.Printf("  - %s\n", url)
			}
		}

		if len(result.Statistics.UniqueDomains) > 0 {
			fmt.Printf("\nВнешние домены:\n")
			for domain, count := range result.Statistics.UniqueDomains {
				fmt.Printf("  - %s: %d ссылок\n", domain, count)
			}
		}

		if len(result.Errors) > 0 {
			fmt.Printf("\nОшибки (%d):\n", len(result.Errors))
			for url, err := range result.Errors {
				fmt.Printf("  - %s: %v\n", url, err)
			}
		}
	}
}
