package dto

type SuggestionRequestDTO struct {
	Prompt string `json:"prompt"`
}

type SuggestionResponseDTO struct {
	Suggestions  []string `json:"suggestions"`
	Improvements []string `json:"improvements"`
}

type GenerateContentRequestDTO struct {
	Prompt string `json:"prompt"`
}

type GenerateContentResponseDTO struct {
	Content string `json:"content"`
}
