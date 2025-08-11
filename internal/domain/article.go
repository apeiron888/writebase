package domain

import (
	"context"
	"time"
)

//========================================================================\\
//                              CONSTANTS                                 ||
//========================================================================//

const (
	DefaultTimeout    = 10 * time.Second
	AITImeout         = 20 * time.Second
	MaxContentLength  = 10000
	MaxTagsPerArticle = 10
	MaxArticleLength  = 5000
	MaxSlugLength     = 100
	MaxTitleLength    = 200
	MaxExcerptLength  = 300
	MaxContentBlocks  = 50
)

//========================================================================\\
//                               MODELS                                   ||
//========================================================================//

type Article struct {
	ID            string
	Title         string
	Slug          string
	AuthorID      string
	ContentBlocks []ContentBlock
	Excerpt       string
	Language      string
	Tags          []string
	Status        ArticleStatus
	Stats         ArticleStats
	Timestamps    ArticleTimes
}

type ArticleStatus string

const (
	StatusDraft     ArticleStatus = "draft"
	StatusScheduled ArticleStatus = "scheduled"
	StatusPublished ArticleStatus = "published"
	StatusArchived  ArticleStatus = "archived"
	StatusDeleted   ArticleStatus = "deleted"
)

type ContentBlock struct {
	Type    BlockType
	Order   int
	Content BlockContent
}

// BlockType defines content block categories
type BlockType string

const (
	BlockHeading    BlockType = "heading"
	BlockParagraph  BlockType = "paragraph"
	BlockImage      BlockType = "image"
	BlockCode       BlockType = "code"
	BlockVideoEmbed BlockType = "video_embed"
	BlockDivider    BlockType = "divider"
	BlockList       BlockType = "list"
)

type BlockContent struct {
	Heading    *HeadingContent
	Paragraph  *ParagraphContent
	Image      *ImageContent
	Code       *CodeContent
	VideoEmbed *VideoEmbedContent
	List       *ListContent
	Divider    *DividerContent
}
type HeadingContent struct {
	Text  string
	Level int //h1, h2, h3, etc.
}
type ParagraphContent struct {
	Text  string
	Style string // e.g., "normal", "italic", "bold"
}
type ImageContent struct {
	URL     string
	Alt     string
	Caption string
}
type CodeContent struct {
	Code     string // Source code content
	Language string
}
type VideoEmbedContent struct {
	Provider string // e.g., "youtube", "vimeo"
	URL      string
}
type ListContent struct {
	Items []string
}
type DividerContent struct {
	Style string // e.g., "solid", "dashed"
}

type ArticleStats struct {
	ViewCount int
	ClapCount int
}

// Timestamps
type ArticleTimes struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	PublishedAt *time.Time
	ArchivedAt  *time.Time
}
// ===========================================================================//
//	                      Article Filter and Pagination                       //
// ===========================================================================//
type ArticleFilter struct {
	// Core filters
	AuthorIDs   []string
	Statuses    []ArticleStatus
	Tags        []string
	ExcludeTags []string
	// Content filters
	Language string
	// Engagement filters
	MinViews int
	MinClaps int
	// Date ranges
	PublishedAfter  *time.Time
	PublishedBefore *time.Time
}
// Pagination controls result slicing
type Pagination struct {
	Page      int
	PageSize  int
	SortField string
	SortOrder string
}
func (p *Pagination) ValidatePagination() {
    // Set default page if not provided or invalid
    if p.Page <= 0 {
        p.Page = 1
    }
    // Set default page size if not provided or invalid
    if p.PageSize <= 0 {
        p.PageSize = 10
    }
    // Limit page size to maximum allowed value
    if p.PageSize > 100 {
        p.PageSize = 100
    }
    // Validate sort order
    if p.SortOrder != "" && p.SortOrder != "asc" && p.SortOrder != "desc" {
        p.SortOrder = "desc" // Default to descending
    }
}


// ===========================================================================//
//	                     Usecase Interface                                    //
// ===========================================================================//
type IArticleUsecase interface {
	CreateArticle(ctx context.Context, userID string, input *Article) (string, error)
	UpdateArticle(ctx context.Context, userID string, input *Article) error
	DeleteArticle(ctx context.Context, articleID,userID string) error
	RestoreArticle(ctx context.Context, userID, articleID string)  error

	GetArticleByID(ctx context.Context, articleID, userID string) (*Article, error)
	GetArticleBySlug(ctx context.Context, slug string, clientIP string) (*Article, error)
	GetArticleStats(ctx context.Context, articleID, userID string) (*ArticleStats, error)
	GetAllArticleStats(ctx context.Context, userID string) ([]ArticleStats,int, error)

	PublishArticle(ctx context.Context, articleID, userID string) (*Article, error)
	UnpublishArticle(ctx context.Context, articleID, userID string) (*Article, error)
	ArchiveArticle(ctx context.Context, articleID, userID string) (*Article, error)
	UnarchiveArticle(ctx context.Context, articleID, userID string) (*Article, error)

	ListArticlesByAuthor(ctx context.Context, userID, authorID string, pag Pagination) ([]Article, int, error)
	GetTrendingArticles(ctx context.Context, userID string, pag Pagination) ([]Article, int, error)
	GetNewArticles(ctx context.Context, userID string, pag Pagination) ([]Article, int, error)
	GetPopularArticles(ctx context.Context, userID string, pag Pagination) ([]Article, int, error) 

	FilterAuthorArticles(ctx context.Context, callerID, authorID string, filter ArticleFilter, pag Pagination) ([]Article, int, error)

	FilterArticles(ctx context.Context, filter ArticleFilter, pag Pagination) ([]Article, int, error) 

	SearchArticles(ctx context.Context, userID, query string, pag Pagination) ([]Article, int, error)

	ListArticlesByTags(ctx context.Context, userID string, tags []string, pag Pagination) ([]Article, int, error)

	EmptyTrash(ctx context.Context, userID string) error
	DeleteArticleFromTrash(ctx context.Context,articleID, userID string) error

	AdminListAllArticles(ctx context.Context, userID, userRole string, pag Pagination) ([]Article, int, error)
	AdminHardDeleteArticle(ctx context.Context, userID, userRole, articleID string) error

	AdminUnpublishArticle(ctx context.Context, userID, userRole, articleID string) (*Article, error) 

	AddClap(ctx context.Context, userID, articleID string) (ArticleStats, error)

	GenerateContentForArticle(ctx context.Context, article *Article, instructions string) (*Article, error)
	GenerateSlugForTitle(ctx context.Context, title string) (string, error)

}
// ===========================================================================//
//	                     Repository Interface                                 //
// ===========================================================================//
type IArticleRepository interface {
	Create(ctx context.Context, article *Article) error
	Update(ctx context.Context, article *Article) error
	Delete(ctx context.Context, articleID string) error
	Restore(ctx context.Context, articleID string) error

	GetByID(ctx context.Context, articleID string) (*Article, error)
	GetBySlug(ctx context.Context, slug string) (*Article, error)
	GetStats(ctx context.Context, articleID string) (*ArticleStats, error)
	GetAllArticleStats(ctx context.Context, userID string) ([]ArticleStats,int, error)

	Publish(ctx context.Context, articleID string, publishAt time.Time) error
	Unpublish(ctx context.Context, articleID string) error
	Archive(ctx context.Context, articleID string, archiveAt time.Time) error
	Unarchive(ctx context.Context, articleID string) error

	ListByAuthor(ctx context.Context, authorID string, pag Pagination) ([]Article, int, error)
	FindTrending(ctx context.Context, windowDays int, pag Pagination) ([]Article, int, error)
	FindNewArticles(ctx context.Context, pag Pagination) ([]Article, int, error)
	FindPopularArticles(ctx context.Context, pag Pagination) ([]Article, int, error)

	FilterAuthorArticles(ctx context.Context, authorID string, filter ArticleFilter, pag Pagination) ([]Article, int, error)

	Filter(ctx context.Context, filter ArticleFilter, pag Pagination) ([]Article, int, error)
	Search(ctx context.Context, query string, pag Pagination) ([]Article, int, error)

	ListByTags(ctx context.Context, tags []string, pag Pagination) ([]Article, int, error)

	EmptyTrash(ctx context.Context, userID string) error
	DeleteFromTrash(ctx context.Context, articleID, userID string) error

	AdminListAllArticles(ctx context.Context, pag Pagination) ([]Article, int, error)
	HardDelete(ctx context.Context, articleID string) error

	IncrementView(ctx context.Context, articleID string) error
	UpdateClapCount(ctx context.Context, articleID string, count int) error
}
// ===========================================================================//
//	                        Article Policy Interface                          //
// ===========================================================================//
type IPolicy interface {
	UserExists(userID string) bool
	ArticleCreateValid(input *Article) bool
	UserOwnsArticle(userID string, input *Article) bool
	CheckArticleChangesAndValid(oldArticle *Article, newArticle *Article) bool
	IsAdmin(userID string, userRole string) bool
}
// ===========================================================================//
//	                          Utils Interface                                 //
// ===========================================================================//
type IUtils interface {
	GenerateUUID() string
	GenerateSlug(title string) string
	GenerateShortUUID() string
	ValidateContent(blocks []ContentBlock) bool
}

