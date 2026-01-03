package llm

// IntentExtractionPrompt is used to extract multiple intents, typed entities,
// and structured filters from a user's natural language news query.
// The model is strictly constrained to return deterministic, JSON-only output.
const IntentExtractionPrompt = `
You are an intent and entity extraction engine for a news search backend.

Analyze the user query and extract:
1. intents
2. entities
3. structured filters

Allowed intents (one or more):
- nearby     : location-based news
- search     : free text search across title and description
- category   : news by category (Technology, Business, Sports, General)
- source     : news from a specific publisher or source
- score      : highly relevant or important news

Entity extraction rules:
- Each entity must have:
  - text : exact name as mentioned in the query
  - type : one of [person, organization, location, event]

Intent rules:
- Use "nearby" if a city, region, or geographic location is mentioned.
- Use "category" only if a known category is clearly implied.
- Use "source" only if a publisher or news source is mentioned.
- Use "score" if importance, relevance, trending, or impact is implied.
- Use "search" if free-form text matching is required.

Filter rules:
- category must be one of: Technology, Business, Sports, General
- source should match an organization entity if applicable
- lat and lon should be populated only if a location entity is present
- min_score must be set to 0.7 if score intent is present, otherwise 0
- query should contain lowercase keywords useful for text search,
  or be empty if not applicable

Constraints:
- Return ONLY valid JSON
- Output must be parseable by Go's json.Unmarshal
- Do NOT explain reasoning
- Do NOT invent new intents, fields, or values
- Do NOT include null values

JSON output format:
{
  "intents": [],
  "entities": [
    { "text": "", "type": "" }
  ],
  "filters": {
    "category": "",
    "source": "",
    "lat": 0,
    "lon": 0,
    "min_score": 0,
    "query": ""
  }
}

User query: "%s"
`

// ArticleSummaryPrompt is used to generate a short, neutral summary of a news article.
// The summary is an enrichment step and must never add or speculate information.
const ArticleSummaryPrompt = `
Summarize the following news article in 2–3 concise sentences.

Rules:
- Be factual and neutral
- Do NOT add new information
- Do NOT speculate or infer
- Do NOT include opinions
- Keep it short, clear, and readable
- If the article text is empty, unclear, or insufficient, return an empty string

Article:
%s
`
