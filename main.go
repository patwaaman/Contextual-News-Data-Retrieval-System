package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"news-retrieval/config"
	"news-retrieval/handler"
	"news-retrieval/llm"
	"news-retrieval/repository"
	"news-retrieval/service"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Fatalf("mongo ping failed: %v", err)
	}

	db := mongoClient.Database("newsdb")

	newsRepo := repository.NewMongoNewsRepository(db, cfg.DebugMongo)

	llm := llm.NewOpenAILLM(cfg.OpenAIKey, cfg.DebugOpenAI)

	newsSvc := service.NewNewsService(newsRepo, llm)

	h := handler.NewNewsHandler(newsSvc)

	mux := http.NewServeMux()

	// ---- Versioned API base ----
	apiV1 := "/api/v1/news"

	// Explicit REST endpoints (NO LLM)
	mux.HandleFunc(apiV1+"/search", h.Search)     // ?query=
	mux.HandleFunc(apiV1+"/category", h.Category) // ?query=
	mux.HandleFunc(apiV1+"/source", h.Source)     // ?query=
	mux.HandleFunc(apiV1+"/score", h.Score)       // ?min=
	mux.HandleFunc(apiV1+"/nearby", h.Nearby)     // ?lat=&lon=&radius=

	// Semantic (LLM-powered)
	mux.HandleFunc(apiV1+"/semantic", h.Semantic) // ?query=&lat=&lon=&radius=

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start server
	go func() {
		log.Println("server started on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	log.Println("shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = server.Shutdown(shutdownCtx)
	_ = mongoClient.Disconnect(shutdownCtx)
}
