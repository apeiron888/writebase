package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"write_base/internal/domain"
)

type GeminiClient struct {
   APIKey string
}

func NewGeminiClient(apiKey string) *GeminiClient {
	return &GeminiClient{APIKey: apiKey}
}

// GetSuggestions generates suggestions and improvements using Gemini API with strict JSON output
func (g *GeminiClient) GetSuggestions(ctx context.Context, req *domain.SuggestionRequest) (*domain.SuggestionResponse, error) {
	   // Prompt: Only JSON, no fluff, both suggestions and improvements
	   prompt := fmt.Sprintf(`Given the topic or keywords: "%s", respond with a JSON object with two fields: "suggestions" (an array of 3 creative blog post ideas) and "improvements" (an array of 3 ways to improve a draft blog post on this topic). Respond ONLY with the JSON object, no extra text.`, req.Prompt)

	   url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=" + g.APIKey

	   body, _ := json.Marshal(map[string]interface{}{
			   "contents": []map[string]interface{}{
					   {"parts": []map[string]string{{"text": prompt}}},
			   },
	   })

	   httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	   if err != nil {
			   return nil, err
	   }
	   httpReq.Header.Set("Content-Type", "application/json")

	   resp, err := http.DefaultClient.Do(httpReq)
	   if err != nil {
			   return nil, err
	   }
	   defer resp.Body.Close()

	   if resp.StatusCode != http.StatusOK {
			   return nil, fmt.Errorf("Gemini API error: %s", resp.Status)
	   }

	   // Parse Gemini API response
	   var geminiResp struct {
			   Candidates []struct {
					   Content struct {
							   Parts []struct {
									   Text string `json:"text"`
							   } `json:"parts"`
					   } `json:"content"`
			   } `json:"candidates"`
	   }
	   data, _ := ioutil.ReadAll(resp.Body)
	   if err := json.Unmarshal(data, &geminiResp); err != nil {
			   return nil, err
	   }

	   if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
			   return nil, fmt.Errorf("No response from Gemini")
	   }

	   // The actual JSON is in the text field
	   var parsed struct {
			   Suggestions  []string `json:"suggestions"`
			   Improvements []string `json:"improvements"`
	   }
	   if err := json.Unmarshal([]byte(geminiResp.Candidates[0].Content.Parts[0].Text), &parsed); err != nil {
			   return nil, fmt.Errorf("Failed to parse Gemini JSON: %w", err)
	   }

	   // Return both suggestions and improvements
	   return &domain.SuggestionResponse{Suggestions: parsed.Suggestions, Improvements: parsed.Improvements}, nil
}

// GenerateContent generates a full blog post using Gemini API
func (g *GeminiClient) GenerateContent(ctx context.Context, req *domain.GenerateContentRequest) (*domain.GenerateContentResponse, error) {
	prompt := fmt.Sprintf(`Write a detailed, engaging blog post about: "%s". Respond ONLY with the blog content, no extra text.`, req.Prompt)

	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=" + g.APIKey

	body, _ := json.Marshal(map[string]interface{}{
		"contents": []map[string]interface{}{
			{"parts": []map[string]string{{"text": prompt}}},
		},
	})

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Gemini API error: %s", resp.Status)
	}

	var geminiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	data, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(data, &geminiResp); err != nil {
		return nil, err
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("No response from Gemini")
	}

	content := geminiResp.Candidates[0].Content.Parts[0].Text
	return &domain.GenerateContentResponse{Content: content}, nil
}