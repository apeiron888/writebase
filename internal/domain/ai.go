package domain

import (
	"context"
)

type IAI interface {
	GenerateContent(ctx context.Context, prompt string) (string, error)
	GenerateSlug(ctx context.Context, title string) (string, error)


	// EditContent(ctx context.Context, content string, instructions string) (string, error)
	// SummarizeContent(ctx context.Context, content string, maxWords int) (string, error)
	// TranslateContent(ctx context.Context, content string, targetLanguage string) (string, error)
}

type SuggestionRequest struct {
	Prompt string
}
type SuggestionResponse struct {
	Suggestions  []string
	Improvements []string
}
type GenerateContentRequest struct {
	Prompt string
}
type GenerateContentResponse struct {
	Content string
}
type IAIUsecase interface {
	GetSuggestions(ctx context.Context, req *SuggestionRequest) (*SuggestionResponse, error)
	GenerateContent(ctx context.Context, req *GenerateContentRequest) (*GenerateContentResponse, error)
}