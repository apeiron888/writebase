package controller

import (
	"time"
	"write_base/internal/domain"
)

type CreateArticleRequest struct {
	Title         string            `json:"title" validate:"required,min=1,max=100"`
	Slug          string            `json:"slug" validate:"required,min=1,max=100"`
	ContentBlocks []ContentBlockDTO `json:"content_blocks" validate:"required,min=1,dive"`
	Excerpt       string            `json:"excerpt" validate:"required,min=1,max=250"`
	Language      string            `json:"language" validate:"required,len=2"`
	Tags          []string          `json:"tags" validate:"required,min=1,max=5,dive,min=1,max=20"`
}

type UpdateArticleRequest struct {
	Title         string           `json:"title,omitempty" validate:"omitempty,min=1,max=100"`
	Slug          string           `json:"slug,omitempty" validate:"omitempty,min=1,max=100"`
	ContentBlocks []ContentBlockDTO `json:"content_blocks,omitempty" validate:"omitempty,min=1,dive"`
	Excerpt       string           `json:"excerpt,omitempty" validate:"omitempty,min=1,max=300"`
	Language      string           `json:"language,omitempty" validate:"omitempty,len=2"`
	Tags          []string          `json:"tags,omitempty" validate:"omitempty,min=1,max=5,dive,min=1,max=10"`
	Status        string           `json:"status,omitempty" validate:"omitempty,oneof=draft published archived"`
}

type ArticleResponse struct {
	ID            string            `json:"id"`
	Title         string            `json:"title"`
	Slug          string            `json:"slug"`
	AuthorID      string            `json:"author_id"`
	ContentBlocks []ContentBlockDTO `json:"content_blocks,omitempty"`
	Excerpt       string            `json:"excerpt"`
	Language      string            `json:"language"`
	Tags          []string          `json:"tags,omitempty"`
	Status        string            `json:"status"`
	Stats         ArticleStatsDTO   `json:"stats"`
	Timestamps    ArticleTimesDTO   `json:"timestamps"`
}

type ArticleListResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	AuthorID  string    `json:"author_id"`
	Excerpt   string    `json:"excerpt"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type ArticleStatsDTO struct {
	ViewsCount int `json:"view_count"`
	ClapCount  int `json:"clap_count"`
}

type ArticleTimesDTO struct {
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	ArchivedAt  *time.Time `json:"archived_at,omitempty"`
}

type ContentBlockDTO struct {
	Type    string          `json:"type" validate:"required,oneof=heading paragraph image code video_embed list divider"`
	Order   int             `json:"order" validate:"required,min=0"`
	Content BlockContentDTO `json:"content" validate:"required"`
}

type BlockContentDTO struct {
	Heading    *HeadingContentDTO    `json:"heading,omitempty"`
	Paragraph  *ParagraphContentDTO  `json:"paragraph,omitempty"`
	Image      *ImageContentDTO      `json:"image,omitempty"`
	Code       *CodeContentDTO       `json:"code,omitempty"`
	VideoEmbed *VideoEmbedContentDTO `json:"video_embed,omitempty"`
	List       *ListContentDTO       `json:"list,omitempty"`
	Divider    *DividerContentDTO    `json:"divider,omitempty"`
}
type HeadingContentDTO struct {
	Text  string `json:"text"`
	Level int    `json:"level"` //h1, h2, h3, etc.
}
type ParagraphContentDTO struct {
	Text  string `json:"text"`
	Style string `json:"style"` //e.g., "normal", "bold", "italic"
}
type ImageContentDTO struct {
	URL     string `json:"url"`
	Alt     string `json:"alt"`
	Caption string `json:"caption"`
}
type CodeContentDTO struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}
type VideoEmbedContentDTO struct {
	Provider string `json:"provider"` // e.g., "youtube", "vimeo"
	URL      string `json:"url"`
}
type ListContentDTO struct {
	Items []string `json:"items"` // e.g., ["Item 1", "Item 2"]
}
type DividerContentDTO struct {
	Style string `json:"style"` // e.g., "solid", "dashed", "dotted"
}

//===========================================================================//
//                        Article Conversion                               //
//=========================================================================//

func mapContentBlocks(dtos []ContentBlockDTO) []domain.ContentBlock {
	var blocks []domain.ContentBlock

	for _, dto := range dtos {
		block := domain.ContentBlock{
			Type:  domain.BlockType(dto.Type),
			Order: dto.Order,
			Content: domain.BlockContent{
				Paragraph:  mapParagraph(dto.Content.Paragraph),
				Heading:    mapHeading(dto.Content.Heading),
				Image:      mapImage(dto.Content.Image),
				Code:       mapCode(dto.Content.Code),
				VideoEmbed: mapVideo(dto.Content.VideoEmbed),
				List:       mapList(dto.Content.List),
				Divider:    mapDivider(dto.Content.Divider),
			},
		}
		blocks = append(blocks, block)
	}

	return blocks
}

func mapParagraph(p *ParagraphContentDTO) *domain.ParagraphContent {
	if p == nil {
		return nil
	}
	return &domain.ParagraphContent{
		Text:  p.Text,
		Style: p.Style,
	}
}

func mapHeading(h *HeadingContentDTO) *domain.HeadingContent {
	if h == nil {
		return nil
	}
	return &domain.HeadingContent{
		Text:  h.Text,
		Level: h.Level,
	}
}

func mapImage(i *ImageContentDTO) *domain.ImageContent {
	if i == nil {
		return nil
	}
	return &domain.ImageContent{
		URL:     i.URL,
		Alt:     i.Alt,
		Caption: i.Caption,
	}
}

func mapCode(c *CodeContentDTO) *domain.CodeContent {
	if c == nil {
		return nil
	}
	return &domain.CodeContent{
		Code:     c.Code,
		Language: c.Language,
	}
}

func mapVideo(v *VideoEmbedContentDTO) *domain.VideoEmbedContent {
	if v == nil {
		return nil
	}
	return &domain.VideoEmbedContent{
		Provider: v.Provider,
		URL:      v.URL,
	}
}

func mapList(l *ListContentDTO) *domain.ListContent {
	if l == nil {
		return nil
	}
	return &domain.ListContent{
		Items: l.Items,
	}
}

func mapDivider(d *DividerContentDTO) *domain.DividerContent {
	if d == nil {
		return nil
	}
	return &domain.DividerContent{
		Style: d.Style,
	}
}
