package service

import (
	"context"
	"fmt"
	"log"
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
	helper.ComputeTextScores(articles, query)
	helper.RankArticles(articles, "search")

	paged, total := paginate(articles, pgn)

	for i := range paged {
		paged[i].Description = s.llm.Summarize(ctx, paged[i].Description)
	}

	for i, a := range articles {
		log.Printf(
			"[rank][search] #%d title=%q relevance=%.2f textScore=%.2f final=%.2f",
			i+1,
			a.Title,
			a.RelevanceScore,
			a.TextScore,
			a.RelevanceScore+a.TextScore,
		)
	}

	return dto.BuildNewsResponse(
		query,
		pgn.Page,
		pgn.Limit,
		total,
		paged,
	), nil
}

func (s *NewsServiceImpl) Category(ctx context.Context, category string, pgn model.Pagination) (dto.NewsResponseDTO, error) {

	articles, err := s.repo.FindByCategory(ctx, category)
	if err != nil {
		return dto.NewsResponseDTO{}, err
	}
	helper.RankArticles(articles, "category")

	paged, total := paginate(articles, pgn)

	for i := range paged {
		paged[i].Description = s.llm.Summarize(ctx, paged[i].Description)
	}

	return dto.BuildNewsResponse(
		category,
		pgn.Page,
		pgn.Limit,
		total,
		paged,
	), nil
}

func (s *NewsServiceImpl) Source(ctx context.Context, source string, pgn model.Pagination) (dto.NewsResponseDTO, error) {
	articles, err := s.repo.FindBySource(ctx, source)
	if err != nil {
		return dto.NewsResponseDTO{}, err
	}
	helper.RankArticles(articles, "source")

	paged, total := paginate(articles, pgn)

	for i := range paged {
		paged[i].Description = s.llm.Summarize(ctx, paged[i].Description)
	}

	return dto.BuildNewsResponse(
		source,
		pgn.Page,
		pgn.Limit,
		total,
		paged,
	), nil
}

func (s *NewsServiceImpl) Score(ctx context.Context, score float64, pgn model.Pagination) (dto.NewsResponseDTO, error) {
	articles, err := s.repo.FindByScore(ctx, score)
	if err != nil {
		return dto.NewsResponseDTO{}, err
	}
	helper.RankArticles(articles, "score")

	paged, total := paginate(articles, pgn)

	for i := range paged {
		paged[i].Description = s.llm.Summarize(ctx, paged[i].Description)
	}

	query := fmt.Sprintf("relevance_score >= %.2f", score)

	return dto.BuildNewsResponse(
		query,
		pgn.Page,
		pgn.Limit,
		total,
		paged,
	), nil
}

func (s *NewsServiceImpl) Nearby(ctx context.Context, lat, lon, radius float64, pgn model.Pagination) (dto.NewsResponseDTO, error) {
	articles, err := s.repo.FindNearby(ctx, lat, lon, radius)
	if err != nil {
		return dto.NewsResponseDTO{}, err
	}
	helper.ComputeDistances(articles, lat, lon)
	helper.RankArticles(articles, "nearby")

	paged, total := paginate(articles, pgn)

	for i := range paged {
		paged[i].Description = s.llm.Summarize(ctx, paged[i].Description)
	}

	for i, a := range articles {
		log.Printf(
			"[rank][nearby] #%d title=%q distance=%.2fkm",
			i+1,
			a.Title,
			a.Distance,
		)
	}

	query := fmt.Sprintf("nearby(lat=%.4f, lon=%.4f, radius=%.1fkm)", lat, lon, radius)

	return dto.BuildNewsResponse(
		query,
		pgn.Page,
		pgn.Limit,
		total,
		paged,
	), nil
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

	for i, a := range res {
		log.Printf(
			"[rank][%s] #%d title=%q relevance=%.2f textScore=%.2f final=%.2f distance=%.2fkm",
			primary,
			i+1,
			a.Title,
			a.RelevanceScore,
			a.TextScore,
			a.RelevanceScore+a.TextScore,
			a.Distance,
		)
	}

	paged, total := paginate(res, pgn)

	for i := range paged {
		paged[i].Description = s.llm.Summarize(ctx, paged[i].Description)
	}

	return dto.BuildNewsResponse(
		q,
		pgn.Page,
		pgn.Limit,
		total,
		paged,
	), nil
}

func fallbackSearch(ctx context.Context, repo repository.NewsRepository, q string, llmRes model.LLMResult) ([]model.NewsArticle, error) {

	// Prefer LLM-derived query if present
	if llmRes.Filters.Query != "" {
		return repo.Search(ctx, llmRes.Filters.Query)
	}

	// Final fallback to raw query
	return repo.Search(ctx, q)
}

func paginate(articles []model.NewsArticle, pgn model.Pagination) ([]model.NewsArticle, int64) {

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

	return articles[start:end], total
}
