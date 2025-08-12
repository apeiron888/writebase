package controller

import (
	"time"
	"write_base/internal/domain"
)

type ArticleRequest struct {
	Title         string            `json:"title" validate:"required,min=1,max=100"`
	Slug          string            `json:"slug" validate:"min=1,max=100"`
	ContentBlocks []ContentBlockDTO `json:"content_blocks" validate:"required,min=1,dive"`
	Excerpt       string            `json:"excerpt" validate:"min=1,max=250"`
	Language      string            `json:"language" validate:"required,len=2"`
	Tags          []string          `json:"tags" validate:"min=1,max=5"`
}

type ArticleUpdateRequest struct {
	ID            string            `json:"id" validate:"required"`
	Title         string            `json:"title" validate:"required,min=1,max=100"`
	Slug          string            `json:"slug" validate:"min=1,max=100"`
	ContentBlocks []ContentBlockDTO `json:"content_blocks" validate:"required,min=1,dive"`
	Excerpt       string            `json:"excerpt" validate:"min=1,max=250"`
	Language      string            `json:"language" validate:"required,len=2"`
	Tags          []string          `json:"tags" validate:"min=1,max=5"`
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
	ID       string `json:"id"`
	Title    string `json:"title"`
	Slug     string `json:"slug"`
	AuthorID string `json:"author_id"`
	Excerpt  string `json:"excerpt"`
	Status   string `json:"status"`
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

type PaginationRequest struct {
	Page     int `form:"page" binding:"min=1"`
	PageSize int `form:"page_size" binding:"min=1,max=100"`
    SortField string `form:"sort_field"`
    SortOrder string `form:"sort_order"`
}



// =======================================================//
//
//	Conversion                        //
//
// =======================================================//
func (aur *ArticleUpdateRequest) ToDomain() *domain.Article {
	return &domain.Article{
		ID:            aur.ID,
		Title:         aur.Title,
		Slug:          aur.Slug,
		Excerpt:       aur.Excerpt,
		Language:      aur.Language,
		Tags:          aur.Tags,
		ContentBlocks: mapContentBlocks(aur.ContentBlocks),
	}
}

func (ar *ArticleRequest) ToDomain() *domain.Article {
	return &domain.Article{
		Title:         ar.Title,
		Slug:          ar.Slug,
		Excerpt:       ar.Excerpt,
		Language:      ar.Language,
		Tags:          ar.Tags,
		ContentBlocks: mapContentBlocks(ar.ContentBlocks),
	}
}

func (ar *ArticleResponse) ToDTO(article *domain.Article) {
	ar.ID = article.ID
	ar.Title = article.Title
	ar.Slug = article.Slug
	ar.AuthorID = article.AuthorID
	ar.ContentBlocks = toContentBlockDTOs(article.ContentBlocks)
	ar.Excerpt = article.Excerpt
	ar.Language = article.Language
	ar.Tags = article.Tags
	ar.Status = string(article.Status)
	ar.Stats = ArticleStatsDTO{ViewsCount: article.Stats.ViewCount, ClapCount: article.Stats.ClapCount}
	ar.Timestamps = ArticleTimesDTO{
		CreatedAt:   article.Timestamps.CreatedAt,
		UpdatedAt:   article.Timestamps.UpdatedAt,
		PublishedAt: article.Timestamps.PublishedAt,
		ArchivedAt:  article.Timestamps.ArchivedAt,
	}
}

func (alr *ArticleListResponse) ToListDTO(article domain.Article) {
	alr.ID = article.ID
	alr.Title = article.Title
	alr.Slug = article.Slug
	alr.AuthorID = article.AuthorID
	alr.Excerpt = article.Excerpt
	alr.Status = string(article.Status)
}

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

func toContentBlockDTOs(blocks []domain.ContentBlock) []ContentBlockDTO {
	var dtos []ContentBlockDTO

	for _, block := range blocks {
		dto := ContentBlockDTO{
			Type:  string(block.Type),
			Order: block.Order,
			Content: BlockContentDTO{
				Heading:    toHeadingContentDTO(block.Content.Heading),
				Paragraph:  toParagraphContentDTO(block.Content.Paragraph),
				Image:      toImageContentDTO(block.Content.Image),
				Code:       toCodeContentDTO(block.Content.Code),
				VideoEmbed: toVideoEmbedContentDTO(block.Content.VideoEmbed),
				List:       toListContentDTO(block.Content.List),
				Divider:    toDividerContentDTO(block.Content.Divider),
			},
		}
		dtos = append(dtos, dto)
	}

	return dtos
}

func toHeadingContentDTO(content *domain.HeadingContent) *HeadingContentDTO {
	if content == nil {
		return nil
	}
	return &HeadingContentDTO{
		Text:  content.Text,
		Level: content.Level,
	}
}

func toParagraphContentDTO(content *domain.ParagraphContent) *ParagraphContentDTO {
	if content == nil {
		return nil
	}
	return &ParagraphContentDTO{
		Text:  content.Text,
		Style: content.Style,
	}
}

func toImageContentDTO(content *domain.ImageContent) *ImageContentDTO {
	if content == nil {
		return nil
	}
	return &ImageContentDTO{
		URL:     content.URL,
		Alt:     content.Alt,
		Caption: content.Caption,
	}
}

func toCodeContentDTO(content *domain.CodeContent) *CodeContentDTO {
	if content == nil {
		return nil
	}
	return &CodeContentDTO{
		Language: content.Language,
		Code:     content.Code,
	}
}

func toVideoEmbedContentDTO(content *domain.VideoEmbedContent) *VideoEmbedContentDTO {
	if content == nil {
		return nil
	}
	return &VideoEmbedContentDTO{
		Provider: content.Provider,
		URL:      content.URL,
	}
}

func toListContentDTO(content *domain.ListContent) *ListContentDTO {
	if content == nil {
		return nil
	}
	return &ListContentDTO{
		Items: content.Items,
	}
}

func toDividerContentDTO(content *domain.DividerContent) *DividerContentDTO {
	if content == nil {
		return nil
	}
	return &DividerContentDTO{
		Style: content.Style,
	}
}