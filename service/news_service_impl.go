package service

import (
	"context"
	"fmt"
	"news-retrieval/dto"
	"news-retrieval/helper"
	"news-retrieval/model"
	"news-retrieval/repository"
)

type NewsServiceImpl struct {
	repo repository.NewsRepository
	llm  LLMService
}

func NewNewsService(repo repository.NewsRepository, llm LLMService) NewsService {
	return &NewsServiceImpl{repo: repo, llm: llm}
}

func (s *NewsServiceImpl) Search(ctx context.Context, query string, pgn model.Pagination) (dto.NewsResponseDTO, error) {
	articles, err := s.repo.Search(ctx, query)
	if err != nil {
		return dto.NewsResponseDTO{}, err
	}

	return buildPaginatedResponse(query, articles, pgn), nil
}

func (s *NewsServiceImpl) Category(ctx context.Context, category string, pgn model.Pagination) (dto.NewsResponseDTO, error) {

	articles, err := s.repo.FindByCategory(ctx, category)
	if err != nil {
		return dto.NewsResponseDTO{}, err
	}

	return buildPaginatedResponse(category, articles, pgn), nil
}

func (s *NewsServiceImpl) Source(ctx context.Context, source string, pgn model.Pagination) (dto.NewsResponseDTO, error) {
	articles, err := s.repo.FindBySource(ctx, source)
	if err != nil {
		return dto.NewsResponseDTO{}, err
	}

	return buildPaginatedResponse(source, articles, pgn), nil
}

func (s *NewsServiceImpl) Score(ctx context.Context, score float64, pgn model.Pagination) (dto.NewsResponseDTO, error) {
	articles, err := s.repo.FindByScore(ctx, score)
	if err != nil {
		return dto.NewsResponseDTO{}, err
	}

	// Use score as query context for metadata
	query := fmt.Sprintf("relevance_score >= %.2f", score)

	return buildPaginatedResponse(query, articles, pgn), nil
}

func (s *NewsServiceImpl) Nearby(ctx context.Context, lat, lon, radius float64, pgn model.Pagination) (dto.NewsResponseDTO, error) {
	articles, err := s.repo.FindNearby(ctx, lat, lon, radius)
	if err != nil {
		return dto.NewsResponseDTO{}, err
	}

	query := fmt.Sprintf("nearby(lat=%.4f, lon=%.4f, radius=%.1fkm)", lat, lon, radius)

	return buildPaginatedResponse(query, articles, pgn), nil
}

func (s *NewsServiceImpl) GetNews(ctx context.Context, q string, opts model.SearchOptions, pgn model.Pagination) (dto.NewsResponseDTO, error) {

	llmRes := s.llm.ParseQuery(ctx, q)
	primary := helper.PickPrimaryIntent(llmRes.Intents)

	var (
		res []model.NewsArticle
		err error
	)

	switch primary {

	case "nearby":
		lat, lon := llmRes.Filters.Lat, llmRes.Filters.Lon
		if lat != 0 && lon != 0 {
			res, err = s.repo.FindNearby(
				ctx,
				lat,
				lon,
				helper.DefaultRadius(opts.Radius),
			)
			break
		}
		res, err = fallbackSearch(ctx, s.repo, q, llmRes)

	case "category":
		if llmRes.Filters.Category != "" {
			res, err = s.repo.FindByCategory(ctx, llmRes.Filters.Category)
			break
		}
		res, err = fallbackSearch(ctx, s.repo, q, llmRes)

	case "source":
		if llmRes.Filters.Source != "" {
			res, err = s.repo.FindBySource(ctx, llmRes.Filters.Source)
			break
		}
		res, err = fallbackSearch(ctx, s.repo, q, llmRes)

	case "score":
		res, err = s.repo.FindByScore(ctx, llmRes.Filters.MinScore)

	default: // search
		res, err = fallbackSearch(ctx, s.repo, q, llmRes)
	}

	if err != nil {
		return dto.NewsResponseDTO{}, err
	}

	res = helper.ApplySecondaryFilters(res, llmRes.Filters)

	if primary == "nearby" {
		helper.ComputeDistances(res, llmRes.Filters.Lat, llmRes.Filters.Lon)
	}
	if primary == "search" {
		helper.ComputeTextScores(res, llmRes.Filters.Query)
	}

	helper.RankArticles(res, primary)

	for i := range res {
		res[i].Description = s.llm.Summarize(ctx, res[i].Description)
	}

	return buildPaginatedResponse(q, res, pgn), nil
}

func fallbackSearch(ctx context.Context, repo repository.NewsRepository, q string, llmRes model.LLMResult) ([]model.NewsArticle, error) {

	// Prefer LLM-derived query if present
	if llmRes.Filters.Query != "" {
		return repo.Search(ctx, llmRes.Filters.Query)
	}

	// Final fallback to raw query
	return repo.Search(ctx, q)
}

func buildPaginatedResponse(query string, articles []model.NewsArticle, pgn model.Pagination) dto.NewsResponseDTO {

	page := pgn.Page
	limit := pgn.Limit
	total := int64(len(articles))

	start := (page - 1) * limit
	if start >= len(articles) {
		start = len(articles)
	}

	end := start + limit
	if end > len(articles) {
		end = len(articles)
	}

	paged := articles[start:end]

	return dto.BuildNewsResponse(
		query,
		page,
		limit,
		total,
		paged,
	)
}
