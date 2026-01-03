package model

import "time"

type NewsArticle struct {
	ID             string    `bson:"_id,omitempty" json:"id"`
	Title          string    `bson:"title" json:"title"`
	Description    string    `bson:"description" json:"description"`
	URL            string    `bson:"url" json:"url"`
	Publication    time.Time `bson:"publicationDate" json:"publication_date"`
	Source         string    `bson:"sourceName" json:"source_name"`
	Category       []string  `bson:"category" json:"category"`
	RelevanceScore float64   `bson:"relevanceScore" json:"relevance_score"`
	Location       struct {
		Type        string    `bson:"type"`
		Coordinates []float64 `bson:"coordinates"`
	} `bson:"location" json:"location"`

	// Computed (not persisted)
	Distance  float64 `bson:"-" json:"distance,omitempty"`
	TextScore float64 `bson:"-" json:"-"`
	Summary   string  `bson:"-" json:"summary,omitempty"`
}
