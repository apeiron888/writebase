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