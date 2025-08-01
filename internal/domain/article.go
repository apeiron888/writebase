package domain

import (
	"time"
	"context"
)

type Article struct {
	ID             string
    Title          string
    Slug           string
    AuthorID       string 
    ContentBlocks []ContentBlock
    Excerpt        string
    CoverImage     string
    Language       string
    Tags           []Tag
    Status         ArticleStatus
    CreatedAt      time.Time
    PublishedAt    time.Time
    UpdatedAt      time.Time
    ViewCount      int
    LikeCount      int
}
// Embedded Content Block Types
type ContentBlockType string
const (
	TextBlockType      ContentBlockType = "text"
	ImageBlockType     ContentBlockType = "image"
	CodeBlockType      ContentBlockType = "code"
	ParagraphBlockType ContentBlockType = "paragraph"
)
type BaseBlock struct {
	Order int `json:"order"`
}
type TextBlock struct {
	BaseBlock
	Text  string
	Style string
}
type ImageBlock struct {
	BaseBlock
	URL     string
	AltText string
	Caption string
}
type CodeBlock struct {
	BaseBlock
	Code     string
	Language string
}
type ContentBlock interface {
	GetType() ContentBlockType
	GetOrder() int
}

type ArticleStatus string
const (
	StatusDraft     ArticleStatus = "draft"
	StatusPublished ArticleStatus = "published"
	StatusArchived  ArticleStatus = "archived"
)

type Tag struct {
    ID          string
    Name        string
    CreatedAt   time.Time
}

// ArticleFilter is used to filter articles based on various criteria
type ArticleFilter struct {
	AuthorIDs   []string
	Status      ArticleStatus
	Tags        []string
	SearchQuery string
	MinClaps    int
	MinViews    int
	AfterDate   time.Time
	BeforeDate  time.Time
}


// Pagination controls result slicing
type Pagination struct {
    Page      int
    PageSize  int
    SortField string
    SortOrder string
}

type IArticleUsecase interface {
	// Core Operations
	CreateArticle(ctx context.Context, article Article) error
	GetArticle(ctx context.Context, articleID string) (Article, error)
	UpdateArticle(ctx context.Context, article Article) error
	DeleteArticle(ctx context.Context, userID, articleID string) error
	
	// Content Discovery
	ListArticles(ctx context.Context, filter ArticleFilter, pagination Pagination) ([]Article, error)
	GetTrendingArticles(ctx context.Context, limit int) ([]Article, error)
	GetRelatedArticles(ctx context.Context, articleID string, limit int) ([]Article, error)
	SearchArticles(ctx context.Context, query string, pagination Pagination) ([]Article, error)
	
	// Engagement
	ViewArticle(ctx context.Context, articleID, userID, ipAddress string) error
	ClapArticle(ctx context.Context, articleID, userID string, count int) error
	// ReportArticle(ctx context.Context, report ArticleReport) error
	
	// Content Management
	PublishArticle(ctx context.Context, userID, articleID string) error
	UnpublishArticle(ctx context.Context, userID, articleID string) error
	ScheduleArticle(ctx context.Context, userID, articleID string, publishAt time.Time) error
}

type IArticleRepository interface {
	// CRUD Operations
	Create(ctx context.Context, article *Article) error
	GetByID(ctx context.Context, id string) (*Article, error)
	Update(ctx context.Context, article *Article) error
	Delete(ctx context.Context, id string) error
	
	// Content Discovery
	List(ctx context.Context, filter ArticleFilter, pagination Pagination) ([]*Article, error)
	GetTrending(ctx context.Context, limit int) ([]*Article, error)
	GetRelated(ctx context.Context, articleID string, limit int) ([]*Article, error)
	
	// Metrics
	IncrementViews(ctx context.Context, articleID string) error
	IncrementClaps(ctx context.Context, articleID string, count int) error
	
	// Search
	Search(ctx context.Context, query string, pagination Pagination) ([]*Article, error)
	
	// Admin Operation
	UpdateStatus(ctx context.Context, articleID string, status ArticleStatus) error
}