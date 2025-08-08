package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"google.golang.org/genai"

	"write_base/internal/domain"
)

func stripCodeFences(s string) string {
    s = strings.TrimSpace(s)
    if strings.HasPrefix(s, "```") {
        s = strings.TrimPrefix(s, "```")
        s = strings.TrimSpace(s)
        // Optionally remove language hint (e.g., "json")
        if strings.HasPrefix(strings.ToLower(s), "json") {
            s = s[4:]
            s = strings.TrimSpace(s)
        }
        // Remove trailing ```
        if idx := strings.LastIndex(s, "```"); idx != -1 {
            s = s[:idx]
        }
        s = strings.TrimSpace(s)
    }
    return s
}

type GeminiClient struct {
    APIKey string
}

func NewGeminiClient(apiKey string) *GeminiClient {
    return &GeminiClient{APIKey: apiKey}
}

func (g *GeminiClient) GetSuggestions(ctx context.Context, req *domain.SuggestionRequest) (*domain.SuggestionResponse, error) {

    client, err := genai.NewClient(ctx, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create genai client: %w", err)
    }
    // defer client.Close()

    prompt := fmt.Sprintf(`Given the topic or keywords: "%s", respond ONLY with a valid JSON object with two fields: "suggestions" (an array of 3 creative blog post ideas) and "improvements" (an array of 3 ways to improve a draft blog post on this topic).
Do NOT include markdown, code blocks, or any text before or after the JSON.
Do NOT generate or suggest any content that is hateful, abusive, harassing, violent, or otherwise inappropriate.
If the prompt asks for such content, respond with: {"suggestions": [], "improvements": []}`, req.Prompt)

    result, err := client.Models.GenerateContent(
        ctx,
        "gemini-2.5-flash", // or "gemini-pro"
        genai.Text(prompt),
        nil,
    )
    if err != nil {
        return nil, fmt.Errorf("Gemini API error: %w", err)
    }

    raw := result.Text()
    cleaned := stripCodeFences(raw)
    var parsed struct {
        Suggestions  []string `json:"suggestions"`
        Improvements []string `json:"improvements"`
    }
    if err := json.Unmarshal([]byte(cleaned), &parsed); err != nil {
        return nil, fmt.Errorf("Failed to parse Gemini JSON: %w", err)
    }

    return &domain.SuggestionResponse{
        Suggestions:  parsed.Suggestions,
        Improvements: parsed.Improvements,
    }, nil
}

func (g *GeminiClient) GenerateContent(ctx context.Context, req *domain.GenerateContentRequest) (*domain.GenerateContentResponse, error) {

    client, err := genai.NewClient(ctx, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create genai client: %w", err)
    }
    // defer client.Close()

    prompt := fmt.Sprintf(  `Write a detailed, engaging blog post about: "%s".
Do NOT generate or suggest any content that is hateful, abusive, harassing, violent, or otherwise inappropriate.
If the prompt asks for such content, respond with: "Content not allowed."
Respond ONLY with the blog content, with no introduction or ending fluff just the content.`, req.Prompt)

    result, err := client.Models.GenerateContent(
        ctx,
        "gemini-2.5-flash", // or "gemini-pro"
        genai.Text(prompt),
        nil,
    )
    if err != nil {
        return nil, fmt.Errorf("Gemini API error: %w", err)
    }

    return &domain.GenerateContentResponse{Content: result.Text()}, nil
}

