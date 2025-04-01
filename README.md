# Rank Vision

A powerful web crawler and SEO analysis tool written in Go, inspired by Ahrefs and Screaming Frog. This tool helps you analyze websites, gather SEO metrics, and identify potential improvements.

## Features

- **Multi-page Crawling**: Crawl entire websites with configurable depth and page limits
- **SEO Analysis**: Collect important SEO metrics including:
  - Meta descriptions
  - Page titles
  - Word count
  - Internal and external links
  - Broken links detection
- **Concurrent Processing**: Multiple workers for efficient crawling
- **Configurable Settings**: Customize request delays, timeouts, and other parameters
- **Detailed Statistics**: Get comprehensive reports about your website's structure and content

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/rank-vision.git
cd rank-vision
```

2. Install dependencies:
```bash
go mod download
```

3. Build the project:
```bash
go build -o crawler cmd/crawler/main.go
```

## Usage

Run the crawler with the following command:

```bash
./crawler -url <website-url> [options]
```

### Available Options

- `-url`: Target website URL (required)
- `-max-pages`: Maximum number of pages to crawl (default: 100)
- `-timeout`: Maximum time to spend crawling (default: 30s)
- `-request-delay`: Delay between requests (default: 1s)

### Example

```bash
./crawler -url https://example.com -max-pages 50 -timeout 60s
```

## Project Structure

```
rank-vision/
├── cmd/
│   └── crawler/
│       └── main.go
├── internal/
│   ├── crawler/
│   │   ├── crawler.go
│   │   └── site_crawler.go
│   └── parser/
│       └── html_parser.go
├── go.mod
└── README.md
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Inspired by Ahrefs and Screaming Frog
- Built with Go
- Uses standard library HTML parsing 