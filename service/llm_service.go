package service

import (
	"context"
	"news-retrieval/model"
)

type LLMService interface {
	ParseQuery(ctx context.Context, q string) model.LLMResult
	Summarize(ctx context.Context, text string) string
}
