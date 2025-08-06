package repository

import (
	"time"
	"fmt"
	"errors"
	"write_base/internal/domain"
)

// =================== Full Article DTO (for full fetch) ===================
type ArticleDTO struct {
	ID            string              `bson:"_id"`
	Title         string              `bson:"title"`
	Slug          string              `bson:"slug"`
	AuthorID      string              `bson:"author_id"`
	ContentBlocks []ContentBlockDTO   `bson:"content_blocks"`
	Excerpt       string              `bson:"excerpt"`
	Language      string              `bson:"language"`
	Tags          []string            `bson:"tags"`
	Status        string              `bson:"status"`
	Stats         ArticleStatsDTO     `bson:"stats"`
	Timestamps    ArticleTimesDTO     `bson:"timestamps"`
}
var ErrArticleToDTO = errors.New("failed to convert Article to ArticleDTO")
// =================== Article List DTO (for list fetch) ===================
type ArticleListDTO struct {
	ID       string `bson:"_id"`
	Title    string `bson:"title"`
	Slug     string `bson:"slug"`
	AuthorID string `bson:"author_id"`
	Excerpt  string `bson:"excerpt"`
	Status   string `bson:"status"`
}
type ArticleStatsDTO struct {
	ViewsCount int `bson:"view_count"`
	ClapCount  int `bson:"clap_count"`
}
type ArticleTimesDTO struct {
	CreatedAt  time.Time `bson:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at"`
	PublishedAt *time.Time `bson:"published_at,omitempty"`
	ArchivedAt  *time.Time `bson:"archived_at,omitempty"`
}

type ContentBlockDTO struct {
	Type    string          `bson:"type"`
	Order   int             `bson:"order"`
	Content BlockContentDTO `bson:"content"`
}
type BlockContentDTO struct {
	Heading    *HeadingContentDTO    `bson:"heading,omitempty"`
	Paragraph  *ParagraphContentDTO  `bson:"paragraph,omitempty"`
	Image      *ImageContentDTO      `bson:"image,omitempty"`
	Code       *CodeContentDTO       `bson:"code,omitempty"`
	VideoEmbed *VideoEmbedContentDTO `bson:"video_embed,omitempty"`
	List       *ListContentDTO       `bson:"list,omitempty"`
	Divider    *DividerContentDTO    `bson:"divider,omitempty"`
}
type HeadingContentDTO struct {
	Text  string `bson:"text"`
	Level int    `bson:"level"` //h1, h2, h3, etc.
}
type ParagraphContentDTO struct {
	Text  string `bson:"text"`
	Style string `bson:"style"` //e.g., "normal", "bold", "italic"
}
type ImageContentDTO struct {
	URL      string `bson:"url"`
	Alt      string `bson:"alt"`
	Caption  string `bson:"caption"`
}
type CodeContentDTO struct {
	Language string `bson:"language"`
	Code     string `bson:"code"`
}
type VideoEmbedContentDTO struct {
	Provider string `bson:"provider"` // e.g., "youtube", "vimeo"
	URL      string `bson:"url"`
}
type ListContentDTO struct {
	Items []string `bson:"items"` // e.g., ["Item 1", "Item 2"]
}
type DividerContentDTO struct {
	Style string `bson:"style"` // e.g., "solid", "dashed", "dotted"
}

//===========================================================================//
//                      To Article Repo DTO Conversion                       //
//===========================================================================//


func ToArticleDTO(article *domain.Article) *ArticleDTO {
	return &ArticleDTO{
		ID:            article.ID,
		Title:         article.Title,
		Slug:          article.Slug,
		AuthorID:      article.AuthorID,
		ContentBlocks: ToContentBlockDTOs(article.ContentBlocks),
		Excerpt:       article.Excerpt,
		Language:      article.Language,
		Tags:          article.Tags,
		Status:       string(article.Status),
		Stats:        ToArticleStatsDTO(article.Stats),
		Timestamps:   ToArticleTimesDTO(article.Timestamps),
	}
}
func ToContentBlockDTOs(blocks []domain.ContentBlock) []ContentBlockDTO {
	var dtos []ContentBlockDTO
	for _, block := range blocks {
		dtos = append(dtos, ContentBlockDTO{
			Type:  string(block.Type),
			Order: block.Order,
			Content: BlockContentDTO{
				Heading:    ToHeadingContentDTO(block.Content.Heading),
				Paragraph:  ToParagraphContentDTO(block.Content.Paragraph),
				Image:      ToImageContentDTO(block.Content.Image),
				Code:       ToCodeContentDTO(block.Content.Code),
				VideoEmbed: ToVideoEmbedContentDTO(block.Content.VideoEmbed),
				List:       ToListContentDTO(block.Content.List),
				Divider:    ToDividerContentDTO(block.Content.Divider),
			},
		})
	}
	return dtos
}
func ToHeadingContentDTO(content *domain.HeadingContent) *HeadingContentDTO {
	if content == nil {
		return nil
	}
	return &HeadingContentDTO{
		Text:  content.Text,
		Level: content.Level,
	}
}
func ToParagraphContentDTO(content *domain.ParagraphContent) *ParagraphContentDTO {
	if content == nil {
		return nil
	}
	return &ParagraphContentDTO{
		Text:  content.Text,
		Style: content.Style,
	}
}
func ToImageContentDTO(content *domain.ImageContent) *ImageContentDTO {
	if content == nil {
		return nil
	}
	return &ImageContentDTO{
		URL:     content.URL,
		Alt:     content.Alt,
		Caption: content.Caption,
	}
}
func ToCodeContentDTO(content *domain.CodeContent) *CodeContentDTO {
	if content == nil {
		return nil
	}
	return &CodeContentDTO{
		Language: content.Language,
		Code:     content.Code,
	}
}
func ToVideoEmbedContentDTO(content *domain.VideoEmbedContent) *VideoEmbedContentDTO {
	if content == nil {
		return nil
	}
	return &VideoEmbedContentDTO{
		Provider: content.Provider,
		URL:      content.URL,
	}
}
func ToListContentDTO(content *domain.ListContent) *ListContentDTO {
	if content == nil {
		return nil
	}
	return &ListContentDTO{
		Items: content.Items,
	}
}
func ToDividerContentDTO(content *domain.DividerContent) *DividerContentDTO {
	if content == nil {
		return nil
	}
	return &DividerContentDTO{
		Style: content.Style,
	}
}
func ToArticleStatsDTO(stats domain.ArticleStats) ArticleStatsDTO {
	return ArticleStatsDTO{
		ViewsCount: stats.ViewCount,
		ClapCount:  stats.ClapCount,
	}
}
func ToArticleTimesDTO(times domain.ArticleTimes) ArticleTimesDTO {
	return ArticleTimesDTO{
		CreatedAt:  times.CreatedAt,
		UpdatedAt:  times.UpdatedAt,
		PublishedAt: times.PublishedAt,
		ArchivedAt:  times.ArchivedAt,
	}
}

//===========================================================================//
//                      From Article Repo DTO Conversion                     //
//===========================================================================//


func FromArticleListDTO(dto *ArticleListDTO) *domain.Article {
	   if dto == nil || dto.ID == "" {
			   fmt.Printf("FromArticleListDTO: got nil or empty ID DTO: %+v\n", dto)
			   return nil
	   }
	   return &domain.Article{
			   ID:       dto.ID,
			   Title:    dto.Title,
			   Slug:     dto.Slug,
			   AuthorID: dto.AuthorID,
			   Excerpt:  dto.Excerpt,
			   Status:   domain.ArticleStatus(dto.Status),
	   }
}
func FromArticleDTO(dto *ArticleDTO) *domain.Article {
	return &domain.Article{
		ID:            dto.ID,
		Title:         dto.Title,
		Slug:          dto.Slug,
		AuthorID:      dto.AuthorID,
		ContentBlocks: FromContentBlockDTOs(dto.ContentBlocks),
		Excerpt:       dto.Excerpt,
		Language:      dto.Language,
		Tags:          dto.Tags,
		Status:       domain.ArticleStatus(dto.Status),
		Stats:        FromArticleStatsDTO(dto.Stats),
		Timestamps:   FromArticleTimesDTO(dto.Timestamps),
	}
}
func FromContentBlockDTOs(dtos []ContentBlockDTO) []domain.ContentBlock {
	var blocks []domain.ContentBlock
	for _, dto := range dtos {
		blocks = append(blocks, domain.ContentBlock{
			Type:  domain.BlockType(dto.Type),
			Order: dto.Order,
			Content: domain.BlockContent{
				Heading:    FromHeadingContentDTO(dto.Content.Heading),
				Paragraph:  FromParagraphContentDTO(dto.Content.Paragraph),
				Image:      FromImageContentDTO(dto.Content.Image),
				Code:       FromCodeContentDTO(dto.Content.Code),
				VideoEmbed: FromVideoEmbedContentDTO(dto.Content.VideoEmbed),
				List:       FromListContentDTO(dto.Content.List),
				Divider:    FromDividerContentDTO(dto.Content.Divider),
			},
		})
	}
	return blocks
}
func FromHeadingContentDTO(dto *HeadingContentDTO) *domain.HeadingContent {
	if dto == nil {
		return nil
	}
	return &domain.HeadingContent{
		Text:  dto.Text,
		Level: dto.Level,
	}
}
func FromParagraphContentDTO(dto *ParagraphContentDTO) *domain.ParagraphContent {
	if dto == nil {
		return nil
	}
	return &domain.ParagraphContent{
		Text:  dto.Text,
		Style: dto.Style,
	}
}
func FromImageContentDTO(dto *ImageContentDTO) *domain.ImageContent {
	if dto == nil {
		return nil
	}
	return &domain.ImageContent{
		URL:     dto.URL,
		Alt:     dto.Alt,
		Caption: dto.Caption,
	}
}
func FromCodeContentDTO(dto *CodeContentDTO) *domain.CodeContent {
	if dto == nil {
		return nil
	}
	return &domain.CodeContent{
		Language: dto.Language,
		Code:     dto.Code,
	}	
}
func FromVideoEmbedContentDTO(dto *VideoEmbedContentDTO) *domain.VideoEmbedContent {
	if dto == nil {
		return nil
	}
	return &domain.VideoEmbedContent{
		Provider: dto.Provider,
		URL:      dto.URL,
	}
}	
func FromListContentDTO(dto *ListContentDTO) *domain.ListContent {
	if dto == nil {
		return nil
	}
	return &domain.ListContent{
		Items: dto.Items,
	}
}
func FromDividerContentDTO(dto *DividerContentDTO) *domain.DividerContent {
	if dto == nil {
		return nil
	}
	return &domain.DividerContent{
		Style: dto.Style,
	}
}
func FromArticleStatsDTO(dto ArticleStatsDTO) domain.ArticleStats {
	return domain.ArticleStats{
		ViewCount: dto.ViewsCount,
		ClapCount: dto.ClapCount,
	}
}
func FromArticleTimesDTO(dto ArticleTimesDTO) domain.ArticleTimes {
	return domain.ArticleTimes{
		CreatedAt:  dto.CreatedAt,
		UpdatedAt:  dto.UpdatedAt,
		PublishedAt: dto.PublishedAt,
		ArchivedAt:  dto.ArchivedAt,
	}
}