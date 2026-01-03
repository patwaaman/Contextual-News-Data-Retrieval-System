package dto

import "news-retrieval/model"

func ToNewsArticleDTO(a model.NewsArticle) NewsArticleDTO {
	dto := NewsArticleDTO{
		Title:           a.Title,
		Description:     a.Description,
		URL:             a.URL,
		PublicationDate: a.Publication,
		SourceName:      a.Source,
		Category:        a.Category,
		RelevanceScore:  a.RelevanceScore,
		LLMSummary:      a.Description, // already summarized
	}

	if len(a.Location.Coordinates) == 2 {
		dto.Longitude = a.Location.Coordinates[0]
		dto.Latitude = a.Location.Coordinates[1]
	}

	return dto
}

func BuildNewsResponse(query string, page, limit int, total int64, articles []model.NewsArticle) NewsResponseDTO {

	dtoArticles := make([]NewsArticleDTO, 0, len(articles))
	for _, a := range articles {
		dtoArticles = append(dtoArticles, ToNewsArticleDTO(a))
	}

	return NewsResponseDTO{
		Metadata: ResponseMetadata{
			Query:   query,
			Page:    page,
			Limit:   limit,
			Total:   total,
			HasNext: int64(page*limit) < total,
		},
		Articles: dtoArticles,
	}
}
