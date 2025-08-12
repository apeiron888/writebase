package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"write_base/internal/domain"
)

var ErrContentPolicyViolation = errors.New("generated content violates policy")

// --- small AI-side structs (JSON-friendly) ---
type aiParagraph struct {
	Text  string `json:"text"`
	Style string `json:"style,omitempty"`
}
type aiHeading struct {
	Text  string `json:"text"`
	Level int    `json:"level,omitempty"`
}
type aiImage struct {
	URL     string `json:"url,omitempty"`
	Alt     string `json:"alt,omitempty"`
	Caption string `json:"caption,omitempty"`
}
type aiCode struct {
	Language string `json:"language,omitempty"`
	Code     string `json:"code,omitempty"`
}
type aiList struct {
	Items []string `json:"items,omitempty"`
}
type aiDivider struct {
	Style string `json:"style,omitempty"`
}

type aiContent struct {
	Heading    *aiHeading `json:"heading,omitempty"`
	Paragraph  *aiParagraph `json:"paragraph,omitempty"`
	Image      *aiImage `json:"image,omitempty"`
	Code       *aiCode `json:"code,omitempty"`
	List       *aiList `json:"list,omitempty"`
	Divider    *aiDivider `json:"divider,omitempty"`
	VideoEmbed map[string]string `json:"video_embed,omitempty"` // flexible
}

type aiBlock struct {
	Type    string    `json:"type"`
	Order   int       `json:"order"`
	Content aiContent `json:"content"`
}

type aiArticle struct {
	Title        string    `json:"title,omitempty"`
	Excerpt      string    `json:"excerpt,omitempty"`
	Tags         []string  `json:"tags,omitempty"`
	ContentBlocks []aiBlock `json:"content_blocks"`
}

// extract first JSON object/array from a string
func extractJSON(s string) (string, bool) {
	firstObj := strings.IndexAny(s, "[{")
	if firstObj == -1 {
		return "", false
	}
	// naive: try to find last matching bracket for arrays or objects
	lastArr := strings.LastIndex(s, "]")
	lastObj := strings.LastIndex(s, "}")
	last := lastArr
	if lastObj > last {
		last = lastObj
	}
	if last == -1 || last <= firstObj {
		return "", false
	}
	return s[firstObj : last+1], true
}

// basic policy filtering — extend with a better/moderation service for production
var bannedPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\bsex\b`),
	regexp.MustCompile(`(?i)\bporn\b`),
	regexp.MustCompile(`(?i)\bxxx\b`),
	regexp.MustCompile(`(?i)\b18\+\b`),
	// add more patterns/slurs as needed
}

func violatesPolicyText(s string) bool {
	low := strings.ToLower(s)
	for _, re := range bannedPatterns {
		if re.MatchString(low) {
			return true
		}
	}
	return false
}

// normalize block type allowed set (make consistent with domain.BlockType)
var allowedBlockTypes = map[string]bool{
	"paragraph":    true,
	"heading":      true,
	"image":        true,
	"code":         true,
	"list":         true,
	"divider":      true,
	"video_embed":  true,
}

// normalize and validate ai blocks -> domain blocks
func aiBlocksToDomainStrict(aiBlocks []aiBlock) ([]domain.ContentBlock, error) {
	out := make([]domain.ContentBlock, 0, len(aiBlocks))
	for _, b := range aiBlocks {
		tt := strings.ToLower(strings.TrimSpace(b.Type))
		if !allowedBlockTypes[tt] {
			return nil, fmt.Errorf("unsupported block type: %s", b.Type)
		}
		// sanitize text fields & detect policy violation
		if b.Content.Paragraph != nil && violatesPolicyText(b.Content.Paragraph.Text) {
			return nil, ErrContentPolicyViolation
		}
		if b.Content.Heading != nil && violatesPolicyText(b.Content.Heading.Text) {
			return nil, ErrContentPolicyViolation
		}
		// image caption / alt check
		if b.Content.Image != nil {
			if violatesPolicyText(b.Content.Image.Caption) || violatesPolicyText(b.Content.Image.Alt) {
				return nil, ErrContentPolicyViolation
			}
		}
		// convert
		cb := domain.ContentBlock{
			Type:  domain.BlockType(tt),
			Order: b.Order,
			Content: domain.BlockContent{
				Paragraph: nil,
				Heading:   nil,
				Image:     nil,
				Code:      nil,
				List:      nil,
				Divider:   nil,
				VideoEmbed: nil,
			},
		}
		if b.Content.Paragraph != nil {
			cb.Content.Paragraph = &domain.ParagraphContent{
				Text:  strings.TrimSpace(b.Content.Paragraph.Text),
				Style: strings.TrimSpace(b.Content.Paragraph.Style),
			}
		}
		if b.Content.Heading != nil {
			cb.Content.Heading = &domain.HeadingContent{
				Text:  strings.TrimSpace(b.Content.Heading.Text),
				Level: b.Content.Heading.Level,
			}
		}
		if b.Content.Image != nil {
			cb.Content.Image = &domain.ImageContent{
				URL:     strings.TrimSpace(b.Content.Image.URL),
				Alt:     strings.TrimSpace(b.Content.Image.Alt),
				Caption: strings.TrimSpace(b.Content.Image.Caption),
			}
		}
		if b.Content.Code != nil {
			cb.Content.Code = &domain.CodeContent{
				Language: strings.TrimSpace(b.Content.Code.Language),
				Code:     strings.TrimSpace(b.Content.Code.Code),
			}
		}
		if b.Content.List != nil {
			cb.Content.List = &domain.ListContent{
				Items: b.Content.List.Items,
			}
		}
		if b.Content.Divider != nil {
			cb.Content.Divider = &domain.DividerContent{
				Style: b.Content.Divider.Style,
			}
		}
		if len(b.Content.VideoEmbed) > 0 {
			cb.Content.VideoEmbed = &domain.VideoEmbedContent{
				Provider: strings.TrimSpace(b.Content.VideoEmbed["provider"]),
				URL:      strings.TrimSpace(b.Content.VideoEmbed["url"]),
			}
		}
		out = append(out, cb)
	}
	// optional: sort by Order to ensure correct order
	// sort.Slice(out, func(i, j int) bool { return out[i].Order < out[j].Order })
	return out, nil
}

// GenerateContentForArticle instructs the AI to produce content following the article structure.
// Returns updated article or ErrContentPolicyViolation if the produced content breaks the rules.
func (u *ArticleUsecase) GenerateContentForArticle(ctx context.Context, article *domain.Article, instructions string) (*domain.Article, error) {
	c, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// Provide the AI with current content block snapshot for context
	// Build minimal context JSON
	type promptBlock struct {
		Type    string      `json:"type"`
		Order   int         `json:"order"`
		Content interface{} `json:"content"`
	}
	promptBlocks := make([]promptBlock, 0, len(article.ContentBlocks))
	for _, b := range article.ContentBlocks {
		cb := promptBlock{Type: string(b.Type), Order: b.Order}
		// only include textual parts for context
		if b.Content.Paragraph != nil {
			cb.Content = map[string]string{"paragraph": b.Content.Paragraph.Text}
		} else if b.Content.Heading != nil {
			cb.Content = map[string]string{"heading": b.Content.Heading.Text}
		} else {
			cb.Content = nil
		}
		promptBlocks = append(promptBlocks, cb)
	}
	promptContext, _ := json.MarshalIndent(promptBlocks, "", "  ")

	// Build a strict prompt
	prompt := fmt.Sprintf(`You are an article assistant. Below is the article metadata and current content blocks.
Instructions: %s

Article metadata (may be empty):
title: %q
excerpt: %q
language: %q
tags: %v

Context blocks (JSON):
%s

OUTPUT: Reply with **ONLY** a single JSON object (no explanation) that may include:
- "title": string (optional)
- "excerpt": string (optional)
- "tags": [string] (optional)
- "content_blocks": [ { "type": "...", "order": 1, "content": { "paragraph": {"text":"..."}, ... } }, ... ]

IMPORTANT POLICY: DO NOT produce sexual content, pornography, explicit adult material, hate slurs, or otherwise offensive content. If any content would violate this policy, either refuse by returning an empty content_blocks array or replace offending text with "[filtered]". Output must be valid JSON only.`, strings.TrimSpace(instructions),
		article.Title, article.Excerpt, article.Language, article.Tags, string(promptContext),
	)

	aiResp, err := u.AIClient.GenerateContent(c, prompt)
	if err != nil {
		return nil, err
	}

	// Extract JSON object from aiResp
	jsonCandidate := aiResp
	if s, ok := extractJSON(aiResp); ok {
		jsonCandidate = s
	}

	// Try to decode into aiArticle
	var parsed aiArticle
	if err := json.Unmarshal([]byte(jsonCandidate), &parsed); err != nil {
		// If the AI returned text (not JSON), treat the whole text as a single paragraph and attempt to sanitize
		plain := strings.TrimSpace(aiResp)
		if plain == "" {
			return nil, domain.ErrInternalServer
		}
		if violatesPolicyText(plain) {
			return nil, ErrContentPolicyViolation
		}
		// replace first paragraph or append
		found := false
		for i := range article.ContentBlocks {
			if article.ContentBlocks[i].Type == domain.BlockParagraph {
				if article.ContentBlocks[i].Content.Paragraph == nil {
					article.ContentBlocks[i].Content.Paragraph = &domain.ParagraphContent{}
				}
				article.ContentBlocks[i].Content.Paragraph.Text = plain
				found = true
				break
			}
		}
		if !found {
			article.ContentBlocks = append(article.ContentBlocks, domain.ContentBlock{
				Type:  domain.BlockParagraph,
				Order: len(article.ContentBlocks) + 1,
				Content: domain.BlockContent{
					Paragraph: &domain.ParagraphContent{Text: plain},
				},
			})
		}
		article.Timestamps.UpdatedAt = time.Now()
		return article, nil
	}

	// Validate parsed blocks
	if len(parsed.ContentBlocks) == 0 {
		// nothing useful; reject
		return nil, ErrContentPolicyViolation
	}

	// Convert and validate each aiBlock to domain.ContentBlock with policy checks
	converted, convErr := aiBlocksToDomainStrict(parsed.ContentBlocks)
	if convErr != nil {
		// If it's a policy violation, bubble up
		if convErr == ErrContentPolicyViolation {
			return nil, ErrContentPolicyViolation
		}
		return nil, convErr
	}

	// Attach parsed metadata (if present) after basic policy checks
	if parsed.Title != "" {
		if violatesPolicyText(parsed.Title) {
			return nil, ErrContentPolicyViolation
		}
		article.Title = parsed.Title
	}
	if parsed.Excerpt != "" {
		if violatesPolicyText(parsed.Excerpt) {
			return nil, ErrContentPolicyViolation
		}
		article.Excerpt = parsed.Excerpt
	}
	if len(parsed.Tags) > 0 {
		// simple tag cleanup (trim/limit)
		cleanTags := make([]string, 0, len(parsed.Tags))
		for _, t := range parsed.Tags {
			t = strings.TrimSpace(t)
			if t == "" {
				continue
			}
			if violatesPolicyText(t) {
				return nil, ErrContentPolicyViolation
			}
			cleanTags = append(cleanTags, t)
			if len(cleanTags) >= 5 { // domain had max 5 tags
				break
			}
		}
		if len(cleanTags) > 0 {
			article.Tags = cleanTags
		}
	}

	// Replace article content blocks with converted ones
	article.ContentBlocks = converted
	article.Timestamps.UpdatedAt = time.Now()
	return article, nil
}


func (u *ArticleUsecase) GenerateSlugForTitle(ctx context.Context, title string) (string, error) {
	c, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	if strings.TrimSpace(title) == "" {
		return "", domain.ErrArticleInvalidSlug
	}
	generatedText, err := u.AIClient.GenerateSlug(c, title)
	if err != nil {
		// If AI fails, fall back to the simple generator
		return u.Utils.GenerateSlug(title), nil
	}
	// 1. Remove potential quotes
	cleanedSlug := strings.Trim(generatedText, "\"")

	// 2. Replace spaces and invalid characters with hyphens
	reg := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	cleanedSlug = reg.ReplaceAllString(cleanedSlug, "-")

	// 3. Remove leading/trailing hyphens and convert to lowercase
	cleanedSlug = strings.Trim(cleanedSlug, "-")
	cleanedSlug = strings.ToLower(cleanedSlug)

	return cleanedSlug, nil
}
