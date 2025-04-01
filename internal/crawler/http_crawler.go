package crawler

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// HTTPCrawler реализует интерфейс Crawler
type HTTPCrawler struct {
	client *http.Client
	config *Config
}

// NewHTTPCrawler создает новый экземпляр HTTPCrawler
func NewHTTPCrawler(config *Config) *HTTPCrawler {
	if config == nil {
		config = DefaultConfig()
	}

	client := &http.Client{
		Timeout: config.Timeout,
	}

	return &HTTPCrawler{
		client: client,
		config: config,
	}
}

// CrawlPage реализует метод интерфейса Crawler
func (c *HTTPCrawler) CrawlPage(ctx context.Context, targetURL string) (*CrawlerResult, error) {
	if !c.IsValidURL(targetURL) {
		return nil, ErrInvalidURL
	}

	req, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.config.UserAgent)

	var resp *http.Response
	for i := 0; i < c.config.MaxRetries; i++ {
		resp, err = c.client.Do(req)
		if err == nil {
			break
		}
		time.Sleep(c.config.RequestDelay)
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrUnexpectedStatusCode
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	result := &CrawlerResult{
		URL: targetURL,
	}

	// Извлекаем заголовок и мета-описание
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "title":
				if n.FirstChild != nil {
					result.Title = n.FirstChild.Data
				}
			case "meta":
				for _, attr := range n.Attr {
					if attr.Key == "name" && attr.Val == "description" {
						for _, descAttr := range n.Attr {
							if descAttr.Key == "content" {
								result.MetaDescription = descAttr.Val
								break
							}
						}
					}
				}
			case "a":
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						result.Links = append(result.Links, attr.Val)
						break
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	// Подсчитываем количество слов
	text := extractText(doc)
	result.WordCount = len(strings.Fields(text))

	return result, nil
}

// IsValidURL реализует метод интерфейса Crawler
func (c *HTTPCrawler) IsValidURL(targetURL string) bool {
	if targetURL == "" {
		return false
	}
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return false
	}
	return parsedURL.Scheme != "" && parsedURL.Host != ""
}

// extractText извлекает текст из HTML-документа
func extractText(n *html.Node) string {
	var text string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			text += n.Data
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	return text
}
