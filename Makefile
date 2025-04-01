# Переменные
URL ?= https://example.com
TIMEOUT = 30s

# User-Agent'ы для тестирования
DESKTOP_UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
MOBILE_UA = "Mozilla/5.0 (iPhone; CPU iPhone OS 17_1_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1"
BOT_UA = "RankVision Bot/1.0 (+https://rank-vision.com; bot@rank-vision.com) Compatible/Go-http-client/1.1"

# Компиляция
.PHONY: build
build:
	go build -o crawler cmd/crawler/main.go

# Тестирование с разными User-Agent
.PHONY: desktop mobile bot all
desktop: build
	@echo "Тестирование с User-Agent для десктопа..."
	./crawler -url $(URL) -user-agent $(DESKTOP_UA) -timeout $(TIMEOUT)

mobile: build
	@echo "Тестирование с User-Agent для мобильного устройства..."
	./crawler -url $(URL) -user-agent $(MOBILE_UA) -timeout $(TIMEOUT)

bot: build
	@echo "Тестирование с User-Agent для бота..."
	./crawler -url $(URL) -user-agent $(BOT_UA) -timeout $(TIMEOUT)

all: desktop mobile bot

# Запуск тестов
.PHONY: test
test:
	go test ./internal/crawler/...

# Очистка
.PHONY: clean
clean:
	rm -f crawler

# Справка
.PHONY: help
help:
	@echo "Использование:"
	@echo "  make desktop URL=https://example.com  - запуск с User-Agent для десктопа"
	@echo "  make mobile URL=https://example.com   - запуск с User-Agent для мобильного"
	@echo "  make bot URL=https://example.com      - запуск с User-Agent для бота"
	@echo "  make all URL=https://example.com      - запуск всех тестов"
	@echo "  make test                            - запуск unit-тестов"
	@echo "  make clean                           - удаление скомпилированного файла"
	@echo ""
	@echo "Параметры:"
	@echo "  URL     - URL для краулинга (по умолчанию: https://example.com)"
	@echo "  TIMEOUT - таймаут запроса (по умолчанию: 30s)" 