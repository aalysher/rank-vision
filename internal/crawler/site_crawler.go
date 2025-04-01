package crawler

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"
)

// SiteCrawler представляет краулер для всего сайта
type SiteCrawler struct {
	baseURL     *url.URL
	config      *Config
	results     map[string]*CrawlerResult
	queue       chan string
	visited     map[string]bool
	visitedLock sync.RWMutex
	maxPages    int
	maxDepth    int
}

// SiteCrawlerResult представляет результат краулинга всего сайта
type SiteCrawlerResult struct {
	BaseURL    string
	TotalPages int
	TotalLinks int
	StartTime  time.Time
	EndTime    time.Time
	Pages      map[string]*CrawlerResult
	Errors     map[string]error
	Statistics SiteStatistics
}

// SiteStatistics содержит статистику по сайту
type SiteStatistics struct {
	TotalWordCount   int
	AverageWordCount int
	UniqueDomains    map[string]int
	BrokenLinks      []string
	MissingMetaDesc  []string
	MissingTitle     []string
}

// NewSiteCrawler создает новый экземпляр SiteCrawler
func NewSiteCrawler(baseURL string, config *Config, maxPages, maxDepth int) (*SiteCrawler, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	if config == nil {
		config = DefaultConfig()
	}

	return &SiteCrawler{
		baseURL:  parsedURL,
		config:   config,
		results:  make(map[string]*CrawlerResult),
		queue:    make(chan string, 1000),
		visited:  make(map[string]bool),
		maxPages: maxPages,
		maxDepth: maxDepth,
	}, nil
}

// CrawlSite запускает краулинг всего сайта
func (sc *SiteCrawler) CrawlSite(ctx context.Context) (*SiteCrawlerResult, error) {
	log.Printf("Начинаем краулинг сайта: %s", sc.baseURL.String())

	startTime := time.Now()
	result := &SiteCrawlerResult{
		BaseURL:   sc.baseURL.String(),
		StartTime: startTime,
		Pages:     make(map[string]*CrawlerResult),
		Errors:    make(map[string]error),
		Statistics: SiteStatistics{
			UniqueDomains: make(map[string]int),
		},
	}

	// Запускаем воркеры для обработки URL
	var wg sync.WaitGroup
	numWorkers := 5 // Количество параллельных воркеров
	log.Printf("Запускаем %d воркеров", numWorkers)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go sc.worker(ctx, &wg, result)
	}

	// Добавляем начальный URL в очередь
	initialURL := sc.baseURL.String()
	log.Printf("Добавляем начальный URL в очередь: %s", initialURL)
	sc.queue <- initialURL

	// Ждем завершения всех воркеров
	done := make(chan bool)
	go func() {
		wg.Wait()
		close(done)
	}()

	// Ждем завершения или таймаута
	select {
	case <-done:
		log.Printf("Все воркеры завершили работу")
	case <-ctx.Done():
		log.Printf("Краулинг прерван по таймауту")
		return result, ctx.Err()
	}

	result.EndTime = time.Now()
	result.TotalPages = len(result.Pages)
	sc.calculateStatistics(result)

	log.Printf("Краулинг завершен. Обработано страниц: %d", result.TotalPages)
	return result, nil
}

// isValidInternalURL проверяет, является ли URL внутренним
func (sc *SiteCrawler) isValidInternalURL(link string) bool {
	// Обрабатываем относительные URL
	if strings.HasPrefix(link, "/") {
		link = sc.baseURL.Scheme + "://" + sc.baseURL.Host + link
	}

	parsedURL, err := url.Parse(link)
	if err != nil {
		return false
	}

	// Если URL относительный, добавляем базовый домен
	if parsedURL.Host == "" {
		parsedURL.Host = sc.baseURL.Host
		parsedURL.Scheme = sc.baseURL.Scheme
	}

	// Проверяем, что URL принадлежит тому же домену
	if parsedURL.Host != sc.baseURL.Host {
		return false
	}

	// Проверяем, что это не файл
	ext := strings.ToLower(path.Ext(parsedURL.Path))
	if ext != "" && ext != ".html" && ext != ".htm" && ext != "/" {
		return false
	}

	return true
}

// worker обрабатывает URL из очереди
func (sc *SiteCrawler) worker(ctx context.Context, wg *sync.WaitGroup, result *SiteCrawlerResult) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case currentURL, ok := <-sc.queue:
			if !ok {
				return
			}

			log.Printf("Обработка URL: %s", currentURL)

			// Проверяем, не посещали ли мы уже этот URL
			sc.visitedLock.RLock()
			visited := sc.visited[currentURL]
			sc.visitedLock.RUnlock()
			if visited {
				log.Printf("URL уже был посещен: %s", currentURL)
				continue
			}

			// Проверяем, не превысили ли мы лимит страниц
			if result.TotalPages >= sc.maxPages {
				log.Printf("Достигнут лимит страниц (%d)", sc.maxPages)
				return
			}

			// Краулим страницу
			crawler := NewHTTPCrawler(sc.config)
			pageResult, err := crawler.CrawlPage(ctx, currentURL)
			if err != nil {
				log.Printf("Ошибка при краулинге %s: %v", currentURL, err)
				result.Errors[currentURL] = err
				continue
			}

			log.Printf("Успешно обработан URL %s. Найдено ссылок: %d", currentURL, len(pageResult.Links))

			// Сохраняем результат
			result.Pages[currentURL] = pageResult
			result.TotalPages++

			// Обрабатываем найденные ссылки
			for _, link := range pageResult.Links {
				// Обрабатываем относительные URL
				if strings.HasPrefix(link, "/") {
					link = fmt.Sprintf("%s://%s%s", sc.baseURL.Scheme, sc.baseURL.Host, link)
				}

				if sc.isValidInternalURL(link) {
					log.Printf("Найдена валидная внутренняя ссылка: %s", link)
					sc.visitedLock.Lock()
					if !sc.visited[link] {
						sc.visited[link] = true
						select {
						case sc.queue <- link:
							log.Printf("Добавлена новая ссылка в очередь: %s", link)
						default:
							log.Printf("Очередь переполнена, пропускаем URL: %s", link)
						}
					}
					sc.visitedLock.Unlock()
				} else {
					log.Printf("Пропущена невалидная или внешняя ссылка: %s", link)
				}
			}

			// Добавляем задержку между запросами
			time.Sleep(sc.config.RequestDelay)
		}
	}
}

// calculateStatistics вычисляет статистику по сайту
func (sc *SiteCrawler) calculateStatistics(result *SiteCrawlerResult) {
	for pageURL, page := range result.Pages {
		// Подсчет слов
		result.Statistics.TotalWordCount += page.WordCount

		// Проверка мета-описания
		if page.MetaDescription == "" {
			result.Statistics.MissingMetaDesc = append(result.Statistics.MissingMetaDesc, pageURL)
		}

		// Проверка заголовка
		if page.Title == "" {
			result.Statistics.MissingTitle = append(result.Statistics.MissingTitle, pageURL)
		}

		// Подсчет внешних ссылок
		for _, link := range page.Links {
			parsedURL, err := url.Parse(link)
			if err != nil {
				continue
			}
			if parsedURL.Host != sc.baseURL.Host {
				result.Statistics.UniqueDomains[parsedURL.Host]++
			}
		}
	}

	// Вычисляем среднее количество слов
	if result.TotalPages > 0 {
		result.Statistics.AverageWordCount = result.Statistics.TotalWordCount / result.TotalPages
	}
}
