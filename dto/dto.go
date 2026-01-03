package dto

import "time"

type ResponseMetadata struct {
	Query   string `json:"query"`
	Page    int    `json:"page"`
	Limit   int    `json:"limit"`
	Total   int64  `json:"total"`
	HasNext bool   `json:"has_next"`
}

type NewsArticleDTO struct {
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	URL             string    `json:"url"`
	PublicationDate time.Time `json:"publication_date"`
	SourceName      string    `json:"source_name"`
	Category        []string  `json:"category"`
	RelevanceScore  float64   `json:"relevance_score"`
	LLMSummary      string    `json:"llm_summary"`
	Latitude        float64   `json:"latitude"`
	Longitude       float64   `json:"longitude"`
}

type NewsResponseDTO struct {
	Metadata ResponseMetadata `json:"metadata"`
	Articles []NewsArticleDTO `json:"articles"`
}
