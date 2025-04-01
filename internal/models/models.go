package models

import "time"

// Domain представляет домен сайта
type Domain struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	URL       string    `json:"url" gorm:"uniqueIndex"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Backlink представляет обратную ссылку
type Backlink struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	SourceURL  string    `json:"source_url"`
	TargetURL  string    `json:"target_url"`
	AnchorText string    `json:"anchor_text"`
	FirstSeen  time.Time `json:"first_seen"`
	LastSeen   time.Time `json:"last_seen"`
	DomainID   uint      `json:"domain_id"`
	Domain     Domain    `json:"domain" gorm:"foreignKey:DomainID"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Keyword представляет ключевое слово
type Keyword struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Keyword   string    `json:"keyword"`
	Position  int       `json:"position"`
	Volume    int       `json:"volume"`
	DomainID  uint      `json:"domain_id"`
	Domain    Domain    `json:"domain" gorm:"foreignKey:DomainID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PageMetrics представляет метрики страницы
type PageMetrics struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	URL             string    `json:"url"`
	Title           string    `json:"title"`
	MetaDescription string    `json:"meta_description"`
	WordCount       int       `json:"word_count"`
	DomainID        uint      `json:"domain_id"`
	Domain          Domain    `json:"domain" gorm:"foreignKey:DomainID"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
