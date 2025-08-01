// ArticleRepository provides CRUD and discovery operations for Article entities in MongoDB.
//
// Usage:
//
//	repo := NewArticleRepository(db, "articles")
//
// Methods:
//   - Create(ctx, article)
//   - GetByID(ctx, articleID)
//   - Update(ctx, article)
//   - Delete(ctx, articleID)
//   - List(ctx, filter, pagination)
//   - GetTrending(ctx, limit)
//   - GetRelated(ctx, articleID, tags, limit)
//   - ViewArticle(ctx, articleID, userID, ipAddress)
//   - ClapArticle(ctx, articleID, userID, count)
//   - UnclapArticle(ctx, articleID, userID, count)
//   - Search(ctx, query, pagination)
//   - UpdateStatus(ctx, articleID, status)
package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"write_base/internal/domain"
)

// =======================================================================================
// General
type ArticleRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewArticleRepository(db *mongo.Database, collectionName string) *ArticleRepository {
	return &ArticleRepository{
		db:         db,
		collection: db.Collection(collectionName),
	}
}

//=======================================================================================

// =======================================================================================
// Data Transfer Object (DTO) for Article
type ArticleDSO struct {
	ID            string                `bson:"_id"`
	Title         string                `bson:"title"`
	Slug          string                `bson:"slug"`
	AuthorID      string                `bson:"author_id"`
	ContentBlocks []domain.ContentBlock `bson:"content_blocks"`
	Excerpt       string                `bson:"excerpt"`
	CoverImage    string                `bson:"cover_image"`
	Language      string                `bson:"language"`
	Tags          []domain.Tag          `bson:"tags"`
	Status        domain.ArticleStatus  `bson:"status"`
	CreatedAt     time.Time             `bson:"created_at"`
	PublishedAt   time.Time             `bson:"published_at"`
	UpdatedAt     time.Time             `bson:"updated_at"`
	ViewCount     int                   `bson:"view_count"`
	ClapCount     int                   `bson:"clap_count"`
}

func (a *ArticleDSO) ToDomain() *domain.Article {
	return &domain.Article{
		ID:            a.ID,
		Title:         a.Title,
		Slug:          a.Slug,
		AuthorID:      a.AuthorID,
		ContentBlocks: a.ContentBlocks,
		Excerpt:       a.Excerpt,
		CoverImage:    a.CoverImage,
		Language:      a.Language,
		Tags:          a.Tags,
		Status:        a.Status,
		CreatedAt:     a.CreatedAt,
		PublishedAt:   a.PublishedAt,
		UpdatedAt:     a.UpdatedAt,
		ViewCount:     a.ViewCount,
		ClapCount:     a.ClapCount,
	}
}
func (a *ArticleDSO) FromDomain(article *domain.Article) {
	a.ID = article.ID
	a.Title = article.Title
	a.Slug = article.Slug
	a.AuthorID = article.AuthorID
	a.ContentBlocks = article.ContentBlocks
	a.Excerpt = article.Excerpt
	a.CoverImage = article.CoverImage
	a.Language = article.Language
	a.Tags = article.Tags
	a.Status = article.Status
	a.CreatedAt = article.CreatedAt
	a.PublishedAt = article.PublishedAt
	a.UpdatedAt = article.UpdatedAt
	a.ViewCount = article.ViewCount
	a.ClapCount = article.ClapCount
}

//=======================================================================================

// =======================================================================================
// ------------------CRUD operations for Article------------------------------------------
func (r *ArticleRepository) Create(ctx context.Context, article *domain.Article) error {
	articleDSO := &ArticleDSO{}
	articleDSO.FromDomain(article)
	_, err := r.collection.InsertOne(ctx, articleDSO)
	return err
}

func (r *ArticleRepository) GetByID(ctx context.Context, articleID string) (*domain.Article, error) {
	var articleDSO ArticleDSO
	err := r.collection.FindOne(ctx, bson.M{"_id": articleID}).Decode(&articleDSO)
	if err != nil {
		return nil, err
	}
	return articleDSO.ToDomain(), nil
}

func (r *ArticleRepository) GetBySlug(ctx context.Context, slug string) (*domain.Article, error) {
	var articleDSO ArticleDSO
	err := r.collection.FindOne(ctx, bson.M{"slug": slug}).Decode(&articleDSO)
	if err != nil {
		return nil, err
	}
	return articleDSO.ToDomain(), nil
}

func (r *ArticleRepository) Update(ctx context.Context, article *domain.Article) error {
	articleDSO := &ArticleDSO{}
	articleDSO.FromDomain(article)
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": article.ID}, bson.M{"$set": articleDSO})
	return err
}

func (r *ArticleRepository) Delete(ctx context.Context, articleID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": articleID})
	return err
}

//----------------------------------------------------------------------------------------

// ------------------Content Discovery Operations------------------------------------------
func (r *ArticleRepository) List(ctx context.Context, filter domain.ArticleFilter, pagination domain.Pagination) ([]*domain.Article, error) {
	var articles []*domain.Article
	opts := options.Find().SetSkip(int64((pagination.Page - 1) * pagination.PageSize)).SetLimit(int64(pagination.PageSize))
	if pagination.SortField != "" {
		opts.SetSort(bson.D{{Key: pagination.SortField, Value: pagination.SortOrder}})
	}
	query := bson.M{}
	if len(filter.AuthorIDs) > 0 {
		query["author_id"] = bson.M{"$in": filter.AuthorIDs}
	}
	if filter.Status != "" {
		query["status"] = filter.Status
	}
	if len(filter.Tags) > 0 {
		query["tags"] = bson.M{"$in": filter.Tags}
	}
	if filter.SearchQuery != "" {
		query["$text"] = bson.M{"$search": filter.SearchQuery}
	}
	if !filter.AfterDate.IsZero() {
		query["created_at"] = bson.M{"$gte": filter.AfterDate}
	}
	if !filter.BeforeDate.IsZero() {
		query["created_at"] = bson.M{"$lte": filter.BeforeDate}
	}
	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var articleDSO ArticleDSO
		if err := cursor.Decode(&articleDSO); err != nil {
			return nil, err
		}
		articles = append(articles, articleDSO.ToDomain())
	}
	return articles, nil
}

func (r *ArticleRepository) GetTrending(ctx context.Context, limit int) ([]*domain.Article, error) {
	var articles []*domain.Article
	opts := options.Find().SetLimit(int64(limit)).SetSort(bson.D{{Key: "view_count", Value: -1}})
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var articleDSO ArticleDSO
		if err := cursor.Decode(&articleDSO); err != nil {
			return nil, err
		}
		articles = append(articles, articleDSO.ToDomain())
	}
	return articles, nil
}

func (r *ArticleRepository) GetRelated(ctx context.Context, articleID string, tags []string, limit int) ([]*domain.Article, error) {
	var articles []*domain.Article
	opts := options.Find().SetLimit(int64(limit)).SetSort(bson.D{{Key: "view_count", Value: -1}})
	query := bson.M{
		"_id":       bson.M{"$ne": articleID},
		"tags.name": bson.M{"$in": tags},
	}
	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var articleDSO ArticleDSO
		if err := cursor.Decode(&articleDSO); err != nil {
			return nil, err
		}
		articles = append(articles, articleDSO.ToDomain())
	}
	return articles, nil
}

//----------------------------------------------------------------------------------------

// -------------------Engagement Operations------------------------------------------------
func (r *ArticleRepository) ViewArticle(ctx context.Context, articleID, userID, ipAddress string) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": articleID}, bson.M{
		"$inc": bson.M{"view_count": 1},
	})
	return err
}

func (r *ArticleRepository) ClapArticle(ctx context.Context, articleID, userID string, count int) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": articleID}, bson.M{
		"$inc": bson.M{"clap_count": count},
	})
	return err
}

func (r *ArticleRepository) UnclapArticle(ctx context.Context, articleID, userID string, count int) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": articleID}, bson.M{
		"$inc": bson.M{"clap_count": -count},
	})
	return err
}

//----------------------------------------------------------------------------------------

// -------------------Search Operations----------------------------------------------------
func (r *ArticleRepository) Search(ctx context.Context, query string, pagination domain.Pagination) ([]*domain.Article, error) {
	var articles []*domain.Article
	opts := options.Find().SetSkip(int64((pagination.Page - 1) * pagination.PageSize)).SetLimit(int64(pagination.PageSize))
	if pagination.SortField != "" {
		opts.SetSort(bson.D{{Key: pagination.SortField, Value: pagination.SortOrder}})
	}
	cursor, err := r.collection.Find(ctx, bson.M{"$text": bson.M{"$search": query}}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var articleDSO ArticleDSO
		if err := cursor.Decode(&articleDSO); err != nil {
			return nil, err
		}
		articles = append(articles, articleDSO.ToDomain())
	}
	return articles, nil
}

//----------------------------------------------------------------------------------------

// --------------------Admin Operations----------------------------------------------------
func (r *ArticleRepository) UpdateStatus(ctx context.Context, articleID string, status domain.ArticleStatus) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": articleID}, bson.M{
		"$set": bson.M{"status": status},
	})
	return err
}

//----------------------------------------------------------------------------------------

//=======================================================================================
