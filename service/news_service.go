package service

import (
	"context"
	"news-retrieval/dto"
	"news-retrieval/model"
)

type NewsService interface {
	Search(ctx context.Context, query string, pgn model.Pagination) (dto.NewsResponseDTO, error)
	Category(ctx context.Context, category string, pgn model.Pagination) (dto.NewsResponseDTO, error)
	Source(ctx context.Context, source string, pgn model.Pagination) (dto.NewsResponseDTO, error)
	Score(ctx context.Context, min float64, pgn model.Pagination) (dto.NewsResponseDTO, error)
	Nearby(ctx context.Context, lat, lon, radius float64, opgn model.Pagination) (dto.NewsResponseDTO, error)

	// -------- Semantic (LLM-powered) --------
	GetNews(ctx context.Context, query string, opts model.SearchOptions, pgn model.Pagination) (dto.NewsResponseDTO, error)
}
