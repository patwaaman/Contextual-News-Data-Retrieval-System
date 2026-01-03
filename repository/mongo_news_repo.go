package repository

import (
	"context"
	"log"
	"news-retrieval/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoNewsRepository struct {
	col        *mongo.Collection
	debugMongo bool
}

func NewMongoNewsRepository(db *mongo.Database, debug bool) NewsRepository {
	return &MongoNewsRepository{col: db.Collection("news_articles"), debugMongo: debug}
}

func (r *MongoNewsRepository) Search(ctx context.Context, q string) ([]model.NewsArticle, error) {
	filter := bson.M{"$text": bson.M{"$search": q}}

	logMongoFilter(r.debugMongo, "search", filter)
	return r.find(ctx, filter)
}

func (r *MongoNewsRepository) FindByCategory(ctx context.Context, c string) ([]model.NewsArticle, error) {
	filter := bson.M{"category": c}

	logMongoFilter(r.debugMongo, "category", filter)
	return r.find(ctx, filter)
}

func (r *MongoNewsRepository) FindBySource(ctx context.Context, s string) ([]model.NewsArticle, error) {
	filter := bson.M{"sourceName": s}

	logMongoFilter(r.debugMongo, "source", filter)
	return r.find(ctx, filter)
}

func (r *MongoNewsRepository) FindByScore(ctx context.Context, min float64) ([]model.NewsArticle, error) {
	filter := bson.M{
		"relevanceScore": bson.M{
			"$gte": min,
		},
	}

	logMongoFilter(r.debugMongo, "score", filter)
	return r.find(ctx, filter)
}

func (r *MongoNewsRepository) FindNearby(
	ctx context.Context,
	lat, lon, km float64,
) ([]model.NewsArticle, error) {

	filter := bson.M{
		"location": bson.M{
			"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{lon, lat}, // Mongo expects [lon, lat]
				},
				"$maxDistance": int64(km * 1000), // km → meters
			},
		},
	}

	logMongoFilter(r.debugMongo, "nearby", filter)
	return r.find(ctx, filter)
}

func (r *MongoNewsRepository) find(ctx context.Context, filter bson.M) ([]model.NewsArticle, error) {
	cur, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var res []model.NewsArticle
	err = cur.All(ctx, &res)
	return res, err
}

func logMongoFilter(debug bool, op string, filter bson.M) {
	if !debug {
		return
	}

	b, err := bson.MarshalExtJSON(filter, true, true)
	if err != nil {
		log.Printf("[mongo][%s] failed to marshal filter: %v", op, err)
		return
	}

	log.Printf("[mongo][%s] filter=%s", op, string(b))
}
