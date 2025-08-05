package domain

import (
	"context"
	"time"
)

//===========================================================================//
//                         Article Definition                                //
//===========================================================================//

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
// Article Status
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
	Text   string
	Style  string // e.g., "normal", "italic", "bold"
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
	URL     string
}
type ListContent struct {
	Items []string 
}
type DividerContent struct {
	Style string // e.g., "solid", "dashed"
}

type ArticleStats struct {
	ViewCount      int
	ClapCount      int
}
// Timestamps
type ArticleTimes struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	PublishedAt *time.Time
	ArchivedAt  *time.Time
}

//===========================================================================//
//                       Article Filter and Pagination                       //
//===========================================================================//
type ArticleFilter struct {
	// Core filters
	AuthorIDs     []string
	Statuses      []ArticleStatus
	Tags          []string
	ExcludeTags   []string
	// Content filters
	Language      string
	// Engagement filters
	MinViews      int
	MinClaps      int
	// Date ranges
	PublishedAfter *time.Time
	PublishedBefore *time.Time
}

// Pagination controls result slicing
type Pagination struct {
    Page      int
    PageSize  int
    SortField string
    SortOrder string
}

//===========================================================================//
//                         Article Policy Interface                          //
//===========================================================================//
type IArticlePolicy interface {
	UserExists(userID string) bool //check existence of userID
	UserOwnsArticle(userID string, article Article) bool //check ownership

	CanViewByID(userID string, userRole string, article Article) bool

	IsAdmin(userID string, userRole string) bool
}

//===========================================================================//
//                       Usecase Interface                                    //
//===========================================================================//

type IArticleUsecase interface {
	// Lifecycle
	CreateArticle(ctx context.Context, UserID string, input *Article) (*Article, error)
	UpdateArticle(ctx context.Context, UserID string, input *Article) (*Article, error)
	DeleteArticle(ctx context.Context, userID, articleID string) error
	RestoreArticle(ctx context.Context, userID, articleID string) (*Article, error)

	// Statistics
	GetArticleStats(ctx context.Context, userID, articleID string) (ArticleStats, error)

	// State transitions
	PublishArticle(ctx context.Context, userID, articleID string) (*Article, error)
	UnpublishArticle(ctx context.Context, userID, articleID string, force bool) (*Article, error)
	ArchiveArticle(ctx context.Context, userID, articleID string) (*Article, error)
	UnarchiveArticle(ctx context.Context, userID, articleID string) (*Article, error)

	// Retrieval
	GetArticleByID(ctx context.Context, viewerID, articleID, userRole string) (*Article, error)
	ViewArticleBySlug(ctx context.Context, slug, clientIP string) (*Article, error)

	// Listing
	ListUserArticles(ctx context.Context, userID, authorID string, pag Pagination) ([]Article, int, error)
	ListTrendingArticles(ctx context.Context, pag Pagination, windowDays int) ([]Article, int, error)
	ListArticlesByTag(ctx context.Context, userID, tag string, pag Pagination) ([]Article, int, error)
	SearchArticles(ctx context.Context, userID, query string, pag Pagination) ([]Article, int, error)
	FilterArticles(ctx context.Context, userID string, filter ArticleFilter, pag Pagination) ([]Article, int, error)

	// Engagement
	ClapArticle(ctx context.Context, articleID, userID string) (ArticleStats, error)

	// Trash
	EmptyTrash(ctx context.Context, userID string) error

	// Admin
	AdminListAllArticles(ctx context.Context, userID, userRole string, filter ArticleFilter, pag Pagination) ([]Article, int, error)
	AdminHardDeleteArticle(ctx context.Context, userID, userRole, articleID string) error
	AdminUnpublishArticle(ctx context.Context, userID, userRole, articleID string) (*Article, error)
}

//===========================================================================//
//                      Repository Interface                                 //
//===========================================================================//

type IArticleRepository interface {
	// CRUD operations
	Insert(ctx context.Context, article *Article) (*Article, error)
	Update(ctx context.Context, article *Article) (*Article, error)
	Delete(ctx context.Context, articleID string) error
	Restore(ctx context.Context, articleID string) (*Article, error)

	// Retrieval
	GetByID(ctx context.Context, articleID string) (*Article, error)
	GetBySlug(ctx context.Context, slug string) (*Article, error)

	// State transition helpers
	Publish(ctx context.Context, articleID string, publishAt time.Time) error
	Unpublish(ctx context.Context, articleID string) error
	Archive(ctx context.Context, articleID string, archiveAt time.Time) error
	Unarchive(ctx context.Context, articleID string) error

	// Author operations
	ListAuthorArticles(ctx context.Context, authorID string, pag Pagination) ([]Article, int, error)
	FilterAuthorArticles(ctx context.Context, authorID string, filter ArticleFilter, pag Pagination) ([]Article, int, error)

	// Lists & search for users
	ListByAuthor(ctx context.Context, authorID string, pag Pagination) ([]Article, int, error)
	ListTrending(ctx context.Context, pag Pagination, windowDays int) ([]Article, int, error)
	ListByTag(ctx context.Context, tag string, pag Pagination) ([]Article, int, error)
	Search(ctx context.Context, query string, pag Pagination) ([]Article, int, error)
	Filter(ctx context.Context, filter ArticleFilter, pag Pagination) ([]Article, int, error)

	// Engagement updates
	IncrementView(ctx context.Context, articleID string) error
	IncrementClap(ctx context.Context, articleID string) error

	// Trash
	EmptyTrash(ctx context.Context) error

	// Admin
	AdminListAllArticles(ctx context.Context, pag Pagination) ([]Article, int, error)
	HardDelete(ctx context.Context, articleID string) error
}
