package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RawArticle struct {
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	URL            string   `json:"url"`
	Publication    string   `json:"publication_date"`
	Source         string   `json:"source_name"`
	Category       []string `json:"category"`
	RelevanceScore float64  `json:"relevance_score"`
	Latitude       float64  `json:"latitude"`
	Longitude      float64  `json:"longitude"`
}

func main() {
	ctx := context.Background()

	// 1️⃣ Read dataset
	data, err := os.ReadFile("news_data.json")
	if err != nil {
		log.Fatal("Failed to read file:", err)
	}

	var articles []RawArticle
	if err := json.Unmarshal(data, &articles); err != nil {
		log.Fatal("Invalid JSON:", err)
	}

	// 2️⃣ Mongo connection
	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI("mongodb://localhost:27017"),
	)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("newsdb")
	col := db.Collection("news_articles")

	// 3️⃣ CREATE INDEXES (IDEMPOTENT)
	createIndexes(ctx, col)

	// 4️⃣ Insert / Upsert data
	for _, a := range articles {

		doc := bson.M{
			"_id":             a.ID,
			"title":           a.Title,
			"description":     a.Description,
			"url":             a.URL,
			"publicationDate": parseTime(a.Publication),
			"sourceName":      a.Source,
			"category":        a.Category,
			"relevanceScore":  a.RelevanceScore,
			"location": bson.M{
				"type":        "Point",
				"coordinates": []float64{a.Longitude, a.Latitude},
			},
		}

		_, err := col.UpdateOne(
			ctx,
			bson.M{"_id": a.ID},
			bson.M{"$setOnInsert": doc},
			options.Update().SetUpsert(true),
		)

		if err != nil {
			log.Println("Skipping article:", a.ID, err)
		}
	}

	log.Println("Seeding completed (indexes + data)")
}

func createIndexes(ctx context.Context, col *mongo.Collection) {

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "title", Value: "text"},
				{Key: "description", Value: "text"},
			},
		},
		{
			Keys: bson.D{{Key: "location", Value: "2dsphere"}},
		},
		{
			Keys: bson.D{{Key: "sourceName", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "category", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "relevanceScore", Value: -1}},
		},
	}

	_, err := col.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		log.Fatal("Index creation failed:", err)
	}

	log.Println("Indexes ensured")
}

func parseTime(value string) time.Time {
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05", // dataset format
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, value); err == nil {
			return t
		}
	}

	log.Println("invalid timestamp, using now():", value)
	return time.Now()
}
