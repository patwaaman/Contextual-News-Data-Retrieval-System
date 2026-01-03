package model

type Entity struct {
	Text string `json:"text"`
	Type string `json:"type"` // person | organization | location | event
}

type Filters struct {
	Category string  `json:"category"`
	Source   string  `json:"source"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	MinScore float64 `json:"min_score"`
	Query    string  `json:"query"`
}

type LLMResult struct {
	Intents  []string `json:"intents"`
	Entities []Entity `json:"entities"`
	Filters  Filters  `json:"filters"`
}

type SearchOptions struct {
	Lat    *float64
	Lon    *float64
	Radius float64 // kilometers
}

type Pagination struct {
	Page  int
	Limit int
}
