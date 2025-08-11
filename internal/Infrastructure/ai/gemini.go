package ai

import (
	"context"
	"fmt"
	"log"
	"strings"
	"write_base/internal/domain"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiClient struct {
	client *genai.GenerativeModel
}

// Ensure GeminiClient implements the IAI interface
var _ domain.IAI = (*GeminiClient)(nil)

func NewGeminiClient(apiKey string) *GeminiClient {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Failed to create Gemini client: %v", err)
	}
	model := client.GenerativeModel("gemini-2.0-flash")
	
	// Configure the model
	model.SetTemperature(0.7)
	model.SetTopK(40)
	model.SetTopP(0.95)
	
	return &GeminiClient{client: model}
}

func (g *GeminiClient) GenerateContent(ctx context.Context, prompt string) (string, error) {
	resp, err := g.client.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}
	
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content generated")
	}
	
	content := string(resp.Candidates[0].Content.Parts[0].(genai.Text))
	return content, nil
}

func (g *GeminiClient) GenerateSlug(ctx context.Context, title string) (string, error) {
	prompt := fmt.Sprintf("Generate a URL-friendly slug for the following title. Return only the slug, nothing else: %s", title)
	resp, err := g.client.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}
	
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content generated")
	}
	
	slug := string(resp.Candidates[0].Content.Parts[0].(genai.Text))
	
	// Clean up the slug
	slug = strings.TrimSpace(slug)
	slug = strings.ToLower(slug)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "--", "-")
	
	// Remove special characters
	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, slug)
	
	return slug, nil
}

func (g *GeminiClient) EditContent(ctx context.Context, content string, instructions string) (string, error) {
	prompt := fmt.Sprintf("Edit the following content based on these instructions: '%s'.\n\nContent:\n%s\n\nReturn only the edited content, nothing else.", instructions, content)
	resp, err := g.client.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}
	
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content generated")
	}
	
	editedContent := string(resp.Candidates[0].Content.Parts[0].(genai.Text))
	return editedContent, nil
}

func (g *GeminiClient) SummarizeContent(ctx context.Context, content string, maxWords int) (string, error) {
	prompt := fmt.Sprintf("Summarize the following content in a maximum of %d words. Return only the summary, nothing else:\n\n%s", maxWords, content)
	resp, err := g.client.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}
	
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content generated")
	}
	
	summary := string(resp.Candidates[0].Content.Parts[0].(genai.Text))
	return summary, nil
}

func (g *GeminiClient) TranslateContent(ctx context.Context, content string, targetLanguage string) (string, error) {
	prompt := fmt.Sprintf("Translate the following content to %s. Return only the translated content, nothing else:\n\n%s", targetLanguage, content)
	resp, err := g.client.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}
	
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content generated")
	}
	
	translatedContent := string(resp.Candidates[0].Content.Parts[0].(genai.Text))
	return translatedContent, nil
}