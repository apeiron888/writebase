package domain

import "context"
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

// Extend IAIUsecase for content generation
type IAIUsecase interface {
	GetSuggestions(ctx context.Context, req *SuggestionRequest) (*SuggestionResponse, error)
	GenerateContent(ctx context.Context, req *GenerateContentRequest) (*GenerateContentResponse, error)
}