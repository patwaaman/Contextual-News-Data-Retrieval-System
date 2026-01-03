package helper

import (
	"news-retrieval/constant"
	"news-retrieval/model"
	util "news-retrieval/utility"
	"sort"
)

var (
	allowedIntent = map[string]bool{
		"search":   true,
		"category": true,
		"source":   true,
		"score":    true,
		"nearby":   true,
	}

	intentPriority = []string{
		"nearby",
		"search",
		"category",
		"source",
		"score",
	}
)

func NormalizeIntents(intents []string) []string {
	var res []string
	for _, i := range intents {
		if allowedIntent[i] {
			res = append(res, i)
		}
	}

	if len(res) == 0 {
		return []string{"search"}
	}
	return res
}

func PickPrimaryIntent(intents []string) string {
	for _, p := range intentPriority {
		for _, i := range intents {
			if p == i {
				return p
			}
		}
	}
	return "search"
}

func DefaultRadius(r float64) float64 {
	if r > 0 {
		return r
	}
	return constant.DefaultRadius
}

func RankArticles(articles []model.NewsArticle, primary string) {
	switch primary {

	case "nearby":
		sort.Slice(articles, func(i, j int) bool {
			return articles[i].Distance < articles[j].Distance
		})

	case "score":
		sort.Slice(articles, func(i, j int) bool {
			return articles[i].RelevanceScore > articles[j].RelevanceScore
		})

	case "search":
		sort.Slice(articles, func(i, j int) bool {
			scoreI := articles[i].RelevanceScore + articles[i].TextScore
			scoreJ := articles[j].RelevanceScore + articles[j].TextScore
			return scoreI > scoreJ
		})

	default: // category, source
		sort.Slice(articles, func(i, j int) bool {
			return articles[i].Publication.After(articles[j].Publication)
		})
	}
}

func ApplySecondaryFilters(articles []model.NewsArticle, f model.Filters) []model.NewsArticle {

	var out []model.NewsArticle

	for _, a := range articles {

		if f.Source != "" && a.Source != f.Source {
			continue
		}
		if f.Category != "" && !util.Contains(a.Category, f.Category) {
			continue
		}
		if f.MinScore > 0 && a.RelevanceScore < f.MinScore {
			continue
		}

		out = append(out, a)
	}

	return out
}

func ComputeDistances(articles []model.NewsArticle, lat, lon float64) {
	for i := range articles {
		if len(articles[i].Location.Coordinates) == 2 {
			aLon := articles[i].Location.Coordinates[0]
			aLat := articles[i].Location.Coordinates[1]
			articles[i].Distance = util.ComputeHaversineDistance(lat, lon, aLat, aLon)
		}
	}
}

func ComputeTextScores(articles []model.NewsArticle, query string) {
	for i := range articles {
		articles[i].TextScore = util.ComputeTextScore(articles[i], query)
	}
}
