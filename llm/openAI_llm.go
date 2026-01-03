package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"news-retrieval/constant"
	"news-retrieval/helper"
	"news-retrieval/model"
	"news-retrieval/utility"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAILLMService struct {
	client      *openai.Client
	debugOpenAI bool
}

func NewOpenAILLM(key string, debugFlag bool) *OpenAILLMService {
	return &OpenAILLMService{client: openai.NewClient(key), debugOpenAI: debugFlag}
}

func (l *OpenAILLMService) ParseQuery(ctx context.Context, q string) model.LLMResult {
	start := time.Now()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	safeQ := strings.ReplaceAll(q, `"`, `\"`)
	prompt := fmt.Sprintf(IntentExtractionPrompt, safeQ)

	logOpenAI(
		l.debugOpenAI,
		"parse_query",
		fmt.Sprintf("prompt=%q", util.Truncate(q, 120)),
	)

	req := openai.ChatCompletionRequest{
		Model: constant.OpenAIModel,
		Messages: []openai.ChatCompletionMessage{
			{Role: constant.UserRole, Content: prompt},
		},
		Temperature: 0,
	}

	resp, err := l.client.CreateChatCompletion(ctx, req)
	if err != nil || len(resp.Choices) == 0 {
		logOpenAI(
			l.debugOpenAI,
			"parse_query",
			fmt.Sprintf("fallback=search err=%v latency=%s", err, time.Since(start)),
		)

		return model.LLMResult{
			Intents: []string{"search"},
			Filters: model.Filters{Query: q},
		}
	}

	content := strings.TrimSpace(resp.Choices[0].Message.Content)
	if content == "" {
		logOpenAI(
			l.debugOpenAI,
			"parse_query",
			"empty_response fallback=search",
		)

		return model.LLMResult{
			Intents: []string{"search"},
			Filters: model.Filters{Query: q},
		}
	}

	var out model.LLMResult
	if err := json.Unmarshal([]byte(content), &out); err != nil {
		logOpenAI(
			l.debugOpenAI,
			"parse_query",
			fmt.Sprintf("invalid_json fallback=search err=%v", err),
		)

		return model.LLMResult{
			Intents: []string{"search"},
			Filters: model.Filters{Query: q},
		}
	}

	out.Intents = helper.NormalizeIntents(out.Intents)

	logOpenAI(
		l.debugOpenAI,
		"parse_query",
		fmt.Sprintf(
			"intents=%v latency=%s",
			out.Intents,
			time.Since(start),
		),
	)

	return out
}

func (l *OpenAILLMService) Summarize(ctx context.Context, text string) string {
	start := time.Now()

	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	if len(strings.TrimSpace(text)) < 30 {
		logOpenAI(
			l.debugOpenAI,
			"summarize",
			"skipped short_text",
		)
		return text
	}

	logOpenAI(
		l.debugOpenAI,
		"summarize",
		fmt.Sprintf("input_len=%d", len(text)),
	)

	safeQ := strings.ReplaceAll(text, `"`, `\"`)
	prompt := fmt.Sprintf(ArticleSummaryPrompt, safeQ)

	resp, err := l.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: constant.OpenAIModel,
		Messages: []openai.ChatCompletionMessage{
			{Role: constant.UserRole, Content: prompt},
		},
		Temperature: 0.2,
		MaxTokens:   120,
	})

	if err != nil || len(resp.Choices) == 0 {
		logOpenAI(
			l.debugOpenAI,
			"summarize",
			fmt.Sprintf("fallback original err=%v latency=%s", err, time.Since(start)),
		)
		return text
	}

	summary := strings.TrimSpace(resp.Choices[0].Message.Content)
	if summary == "" {
		logOpenAI(
			l.debugOpenAI,
			"summarize",
			"empty_response fallback original",
		)
		return text
	}

	summary = strings.Join(strings.Fields(summary), " ")

	logOpenAI(
		l.debugOpenAI,
		"summarize",
		fmt.Sprintf("summary_len=%d latency=%s", len(summary), time.Since(start)),
	)

	return summary
}

func logOpenAI(debug bool, op string, msg string) {
	if !debug {
		return
	}
	log.Printf("[openai][%s] %s", op, msg)
}
