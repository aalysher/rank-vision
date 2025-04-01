package config

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Crawler  CrawlerConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type CrawlerConfig struct {
	MaxConcurrentRequests int
	UserAgent             string
	RequestDelay          int // в миллисекундах
}

func NewConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: ":8080",
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "postgres",
			DBName:   "rank_vision",
		},
		Crawler: CrawlerConfig{
			MaxConcurrentRequests: 10,
			UserAgent:             "RankVision Bot/1.0",
			RequestDelay:          1000,
		},
	}
}
