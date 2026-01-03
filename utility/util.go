package util

import (
	"math"
	"news-retrieval/model"
	"strings"
)

func ComputeHaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth radius in km

	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180

	lat1 = lat1 * math.Pi / 180
	lat2 = lat2 * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

func ComputeTextScore(article model.NewsArticle, query string) float64 {
	if query == "" {
		return 0
	}

	text := strings.ToLower(article.Title + " " + article.Description)
	terms := strings.Fields(strings.ToLower(query))

	score := 0.0
	for _, t := range terms {
		if strings.Contains(text, t) {
			score += 1.0
		}
	}
	return score
}

func Contains(list []string, v string) bool {
	for _, s := range list {
		if s == v {
			return true
		}
	}
	return false
}

func Truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
