## 📰 Contextual News Data Retrieval System

Go + MongoDB + OpenAI (LLM) + Docker

A production-grade backend system for retrieving, ranking, and enriching news articles using explicit REST APIs and LLM-powered semantic search.

This service supports:

* Explicit REST APIs (search, category, source, score, nearby)

* LLM-powered semantic queries in natural English

* Multi-intent query resolution

* Intent-aware ranking strategies

* MongoDB text & geo-spatial search

* Pagination with metadata

* LLM-generated article summaries (post-pagination only)

* Debug logging for MongoDB & OpenAI

* Fully Dockerized setup (App + MongoDB)


## 📁 Folder Structure
```text
contextual-news-data-retrieval-system/
  config/          # Env config loader
  constant/        # App constants
  dto/             # Response DTOs
  errconst/        # Shared error definitions
  handler/         # HTTP handlers
  helper/          # Ranking, scoring, geo utils
  llm/             # OpenAI integration
  model/           # Domain models
  repository/      # MongoDB data access
  service/         # Business logic
  main.go
  Dockerfile
  docker-compose.yml
  .env.example
  README.md
```

## ⚙️ Setup & Run Instructions
```text
1️⃣ Clone the repository
git clone <your-repo-url>
cd contextual-news-data-retrieval-system

2️⃣ Create .env file
cp .env.example .env

3️⃣ Run with Docker Compose
docker compose up --build

Services available:

Component      URL
REST API       http://localhost:8080
MongoDB        mongodb://localhost:27017

4️⃣ View logs (recommended)
docker compose logs -f app
```

## 🌐 REST API Documentation
```text
Postman Collection added: import it.

Base path: /api/v1/news

1️⃣ Search News
GET /search
curl "http://localhost:8080/api/v1/news/search?query=ai&page=1&limit=5"

2️⃣ Category News
GET /category
curl "http://localhost:8080/api/v1/news/category?query=Technology&page=1&limit=5"

3️⃣ Source News
GET /source
curl "http://localhost:8080/api/v1/news/source?query=Reuters&page=1&limit=5"

4️⃣ Score-based News
GET /score
curl "http://localhost:8080/api/v1/news/score?query=0.8&page=1&limit=5"


Minimum score threshold: 0.7

5️⃣ Nearby News (Geo-based)
GET /nearby
curl "http://localhost:8080/api/v1/news/nearby?lat=12.97&lon=77.59&radius=5&page=1&limit=5"

6️⃣ Semantic Search (LLM-Powered)
GET /semantic
curl "http://localhost:8080/api/v1/news/semantic?query=high relevance tech news near bangalore"

```

## OUTPUT FORMAT 
```text
The LLM:

* Extracts multiple intents

* Generates structured filters

* Routes query to the correct backend logic

📦 Response Format
{
  "metadata": {
    "query": "tech news near bangalore",
    "page": 1,
    "limit": 5,
    "total": 42,
    "hasNext": true
  },
  "articles": [
    {
      "title": "Article title",
      "description": "Original description",
      "url": "https://example.com",
      "publication_date": "2025-03-24T11:08:11Z",
      "source_name": "Reuters",
      "category": ["Technology"],
      "relevance_score": 0.91,
      "llm_summary": "Short LLM-generated summary...",
      "latitude": 12.9716,
      "longitude": 77.5946,
      "distance": 2.3
    }
  ]
}
```

## 🧠 Intent-Aware Ranking Logic
```text
Intent:	   Ranking Strategy
search:	   relevance_score + text match score
category:  most recent publication
source:	   most recent publication
score:	   highest relevance_score
nearby:    shortest distance (Haversine)

LLM Usage Strategy

Used for:

* Intent extraction

* Filter extraction

* Article summarization

* Never used before pagination

* Deterministic JSON-only prompts

* Timeouts and fallbacks applied

* Pagination is a hard boundary before any LLM call.
```

## 🗄️ MongoDB Indexes
```text
db.news_articles.createIndex({ title: "text", description: "text" })
db.news_articles.createIndex({ category: 1 })
db.news_articles.createIndex({ sourceName: 1 })
db.news_articles.createIndex({ relevanceScore: -1 })
db.news_articles.createIndex({ location: "2dsphere" })
```


## 🧠 Architectural Overview
```text
Client
  |
  v
HTTP Handlers
  |
  v
Service Layer
  |   ├─ Intent Resolution (LLM)
  |   ├─ Ranking
  |   ├─ Pagination
  |   └─ LLM Enrichment
  |
  v
MongoDB
```

## 🧩 Design Decisions & Trade-offs

```text
1️⃣ MongoDB as Primary Store

* Flexible schema

* Strong text & geo queries

* Horizontal scalability

* Trade-off: Less transactional than RDBMS (acceptable here).

2️⃣ LLM Only at Service Layer

* Clear separation of concerns

* Prevents accidental overuse

* Cost-controlled enrichment

3️⃣ Ranking Before Pagination

* Correct global ordering

* Deterministic results

* Fair scoring

4️⃣ Pagination Before LLM

* Prevents token overuse

* Predictable cost

* Faster responses

```

## 🚀 Future Improvements

```text

* Async LLM summarization

* Caching ranked results

* Distributed rate limiting

* OpenAPI / Swagger docs

* Observability metrics
```