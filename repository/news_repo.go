package repository

import (
	"context"
	"news-retrieval/model"
)

type NewsRepository interface {
	Search(ctx context.Context, query string) ([]model.NewsArticle, error)
	FindByCategory(ctx context.Context, category string) ([]model.NewsArticle, error)
	FindBySource(ctx context.Context, source string) ([]model.NewsArticle, error)
	FindByScore(ctx context.Context, min float64) ([]model.NewsArticle, error)
	FindNearby(ctx context.Context, lat, lon, km float64) ([]model.NewsArticle, error)
}
