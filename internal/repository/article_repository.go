package repository

import (
	"context"
	"time"
	"write_base/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ArticleRepository struct {
	Collection *mongo.Collection
}

func NewArticleRepository(db *mongo.Database, collection_name string) *ArticleRepository {
	return &ArticleRepository{
		Collection: db.Collection(collection_name),
	}
}

// ===========================================================================//
//	Article Lifecycle (Author only)                                           //
// ===========================================================================//
func (r *ArticleRepository) Insert(ctx context.Context, article *domain.Article) (*domain.Article, error) {
	articleDTO := ToArticleDTO(article)
	_, err := r.Collection.InsertOne(ctx, articleDTO)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, domain.ErrArticleAlreadyExists
		}
		return nil, err
	}
	return article, nil
}
func (r *ArticleRepository) Update(ctx context.Context, article *domain.Article) (*domain.Article, error) {
	articleDTO := ToArticleDTO(article)
	res, err := r.Collection.UpdateOne(ctx, bson.M{"_id": article.ID}, bson.M{"$set": articleDTO})
	if err != nil {
		return nil, err
	}
	if res.MatchedCount == 0 {
		return nil, domain.ErrArticleNotFound
	}
	return article, nil
}
func (r *ArticleRepository) Delete(ctx context.Context, articleID string) error {
	res, err := r.Collection.UpdateOne(ctx, bson.M{"_id": articleID}, bson.M{"$set": bson.M{"status": string(domain.StatusDeleted)}})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return domain.ErrArticleNotFound
	}
	return nil
}

func (r *ArticleRepository) Restore(ctx context.Context, articleID string) (*domain.Article, error) {
	res, err := r.Collection.UpdateOne(ctx, bson.M{"_id": articleID}, bson.M{"$set": bson.M{"status": string(domain.StatusDraft)}})
	if err != nil {
		return nil, err
	}
	if res.MatchedCount == 0 {
		return nil, domain.ErrArticleNotFound
	}
	var articleDTO ArticleDTO
	err = r.Collection.FindOne(ctx, bson.M{"_id": articleID}).Decode(&articleDTO)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrArticleNotFound
		}
		return nil, err
	}
	return FromArticleDTO(&articleDTO), nil
}

// ===========================================================================//
//	Article Retrieval                                                         //
// ===========================================================================//
func (r *ArticleRepository) GetByID(ctx context.Context, articleID string) (*domain.Article, error) {
	var articleDTO ArticleDTO
	err := r.Collection.FindOne(ctx, bson.M{"_id": articleID}).Decode(&articleDTO)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrArticleNotFound
		}
		return nil, err
	}
	return FromArticleDTO(&articleDTO), nil
}
func (r *ArticleRepository) GetBySlug(ctx context.Context, slug string) (*domain.Article, error) {
	var articleDTO ArticleDTO
	err := r.Collection.FindOne(ctx, bson.M{"slug": slug, "status": string(domain.StatusPublished)}).Decode(&articleDTO)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrArticleNotFound
		}
		return nil, err
	}
	return FromArticleDTO(&articleDTO), nil
}

// ===========================================================================//
//	Article State Management (Author only)                                    //
// ===========================================================================//
func (r *ArticleRepository) Publish(ctx context.Context, articleID string, publishAt time.Time) error {
	_, err := r.Collection.UpdateOne(ctx, bson.M{"_id": articleID}, bson.M{"$set": bson.M{"status": string(domain.StatusPublished), "timestamps.published_at": publishAt}})
	if err != nil {
		return err
	}
	return nil
}
func (r *ArticleRepository) Unpublish(ctx context.Context, articleID string) error {
	_, err := r.Collection.UpdateOne(ctx, bson.M{"_id": articleID}, bson.M{"$set": bson.M{"status": string(domain.StatusDraft)}})
	if err != nil {
		return err
	}
	return nil
}
func (r *ArticleRepository) Archive(ctx context.Context, articleID string, archiveAt time.Time) error {
	_, err := r.Collection.UpdateOne(ctx, bson.M{"_id": articleID}, bson.M{"$set": bson.M{"status": string(domain.StatusArchived), "timestamps.archived_at": archiveAt}})
	if err != nil {
		return err
	}
	return nil
}
func (r *ArticleRepository) Unarchive(ctx context.Context, articleID string) error {
	_, err := r.Collection.UpdateOne(ctx, bson.M{"_id": articleID}, bson.M{"$set": bson.M{"status": string(domain.StatusDraft)}})
	if err != nil {
		return err
	}
	return nil
}

// ===========================================================================//
//                    	Article operations                                    //
// ===========================================================================//
func (r *ArticleRepository) ListAuthorArticles(ctx context.Context, authorID string, pag domain.Pagination) ([]domain.Article, int, error) {
	var articles []domain.Article
	opts := options.Find().SetSkip(int64(pag.Page * pag.PageSize)).SetLimit(int64(pag.PageSize))
	cursor, err := r.Collection.Find(ctx, bson.M{"author_id": authorID}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var articleDTO ArticleListDTO
		if err := cursor.Decode(&articleDTO); err != nil {
			return nil, 0, err
		}
		article := FromArticleListDTO(&articleDTO)
		articles = append(articles, *article)
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	total, err := r.Collection.CountDocuments(ctx, bson.M{"author_id": authorID})
	if err != nil {
		return nil, 0, err
	}

	return articles, int(total), nil
}
func (r *ArticleRepository) FilterAuthorArticles(ctx context.Context, authorID string, filter domain.ArticleFilter, pag domain.Pagination) ([]domain.Article, int, error) {
	query := buildArticleFilterQuery(authorID, filter)
	opts := options.Find().SetSkip(int64(pag.Page * pag.PageSize)).SetLimit(int64(pag.PageSize))
	cursor, err := r.Collection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	var articles []domain.Article
	for cursor.Next(ctx) {
		var articleDTO ArticleListDTO
		if err := cursor.Decode(&articleDTO); err != nil {
			return nil, 0, err
		}
		article := FromArticleListDTO(&articleDTO)
		if article == nil {
			return nil, 0, ErrArticleToDTO
		}
		articles = append(articles, *article)
	}
	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}
	total, err := r.Collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	return articles, int(total), nil
}

func (r *ArticleRepository) Filter(ctx context.Context, filter domain.ArticleFilter, pag domain.Pagination) ([]domain.Article, int, error) {
	query := buildArticleFilterQuery("", filter)
	opts := options.Find().SetSkip(int64(pag.Page * pag.PageSize)).SetLimit(int64(pag.PageSize))
	cursor, err := r.Collection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	var articles []domain.Article
	for cursor.Next(ctx) {
		var articleDTO ArticleListDTO
		if err := cursor.Decode(&articleDTO); err != nil {
			return nil, 0, err
		}
		article := FromArticleListDTO(&articleDTO)
		if article == nil {
			return nil, 0, ErrArticleToDTO
		}
		articles = append(articles, *article)
	}
	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}
	total, err := r.Collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	return articles, int(total), nil
}

// Helper to build a BSON query from ArticleFilter
func buildArticleFilterQuery(authorID string, filter domain.ArticleFilter) bson.M {
	query := bson.M{}
	if authorID != "" {
		query["author_id"] = authorID
	}
	if len(filter.AuthorIDs) > 0 {
		query["author_id"] = bson.M{"$in": filter.AuthorIDs}
	}
	if len(filter.Statuses) > 0 {
		statuses := make([]string, len(filter.Statuses))
		for i, s := range filter.Statuses {
			statuses[i] = string(s)
		}
		query["status"] = bson.M{"$in": statuses}
	}
	if len(filter.Tags) > 0 {
		query["tags"] = bson.M{"$in": filter.Tags}
	}
	if len(filter.ExcludeTags) > 0 {
		query["tags"] = bson.M{"$nin": filter.ExcludeTags}
	}
	if filter.Language != "" {
		query["language"] = filter.Language
	}
	if filter.MinViews > 0 {
		query["stats.view_count"] = bson.M{"$gte": filter.MinViews}
	}
	if filter.MinClaps > 0 {
		query["stats.clap_count"] = bson.M{"$gte": filter.MinClaps}
	}
	if filter.PublishedAfter != nil && !filter.PublishedAfter.IsZero() {
		query["timestamps.published_at"] = bson.M{"$gte": filter.PublishedAfter}
	}
	if filter.PublishedBefore != nil && !filter.PublishedBefore.IsZero() {
		if q, ok := query["timestamps.published_at"].(bson.M); ok {
			q["$lte"] = filter.PublishedBefore
			query["timestamps.published_at"] = q
		} else {
			query["timestamps.published_at"] = bson.M{"$lte": filter.PublishedBefore}
		}
	}
	return query
}

// ===========================================================================//
//	                          Article Lists                                   //
// ===========================================================================//
func (r *ArticleRepository) ListByAuthor(ctx context.Context, authorID string, pag domain.Pagination) ([]domain.Article, int, error) {
	var articles []domain.Article
	opts := options.Find().SetSkip(int64(pag.Page * pag.PageSize)).SetLimit(int64(pag.PageSize))
	query := bson.M{"author_id": authorID, "status": string(domain.StatusPublished)}
	cursor, err := r.Collection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var articleDTO ArticleListDTO
		if err := cursor.Decode(&articleDTO); err != nil {
			return nil, 0, err
		}
		article := FromArticleListDTO(&articleDTO)
		articles = append(articles, *article)
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	total, err := r.Collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	return articles, int(total), nil
}

func (r *ArticleRepository) ListTrending(ctx context.Context, pag domain.Pagination, windowDays int) ([]domain.Article, int, error) {
	var articles []domain.Article
	windowAgo := time.Now().AddDate(0, 0, -windowDays)
	query := bson.M{
		"status":                  string(domain.StatusPublished),
		"timestamps.published_at": bson.M{"$gte": windowAgo},
	}
	opts := options.Find().
		SetSkip(int64(pag.Page * pag.PageSize)).
		SetLimit(int64(pag.PageSize)).
		SetSort(bson.D{{Key: "stats.view_count", Value: -1}})

	cursor, err := r.Collection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var articleDTO ArticleListDTO
		if err := cursor.Decode(&articleDTO); err != nil {
			return nil, 0, err
		}
		article := FromArticleListDTO(&articleDTO)
		articles = append(articles, *article)
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	total, err := r.Collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	return articles, int(total), nil
}
func (r *ArticleRepository) ListByTag(ctx context.Context, tag string, pag domain.Pagination) ([]domain.Article, int, error) {
	var articles []domain.Article
	opts := options.Find().SetSkip(int64(pag.Page * pag.PageSize)).SetLimit(int64(pag.PageSize))
	query := bson.M{"tags": tag, "status": string(domain.StatusPublished)}
	cursor, err := r.Collection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var articleDTO ArticleListDTO
		if err := cursor.Decode(&articleDTO); err != nil {
			return nil, 0, err
		}
		article := FromArticleListDTO(&articleDTO)
		articles = append(articles, *article)
	}
	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	total, err := r.Collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	return articles, int(total), nil
}
func (r *ArticleRepository) Search(ctx context.Context, query string, pag domain.Pagination) ([]domain.Article, int, error) {
	articles := []domain.Article{}
	filter := bson.M{"$text": bson.M{"$search": query}}
	opts := options.Find().SetSkip(int64(pag.Page * pag.PageSize)).SetLimit(int64(pag.PageSize))
	cursor, err := r.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var articleDTO ArticleListDTO
		if err := cursor.Decode(&articleDTO); err != nil {
			return nil, 0, err
		}
		article := FromArticleListDTO(&articleDTO)
		articles = append(articles, *article)
	}
	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}
	total, err := r.Collection.CountDocuments(ctx, bson.M{"$text": bson.M{"$search": query}})
	if err != nil {
		return nil, 0, err
	}
	return articles, int(total), nil
}

// ===========================================================================//
//                    	    Engagement updates                                //
// ===========================================================================//
func (r *ArticleRepository) IncrementView(ctx context.Context, articleID string) error {
	_, err := r.Collection.UpdateOne(ctx, bson.M{"_id": articleID, "status": string(domain.StatusPublished)}, bson.M{"$inc": bson.M{"stats.view_count": 1}})
	return err
}
func (r *ArticleRepository) IncrementClap(ctx context.Context, articleID string) error {
	_, err := r.Collection.UpdateOne(ctx, bson.M{"_id": articleID, "status": string(domain.StatusPublished)}, bson.M{"$inc": bson.M{"stats.clap_count": 1}})
	return err
}

// ===========================================================================//
//	                           Trash Management                               //
// ===========================================================================//
func (r *ArticleRepository) EmptyTrash(ctx context.Context) error {
	_, err := r.Collection.DeleteMany(ctx, bson.M{"status": string(domain.StatusDeleted)})
	return err
}

// ===========================================================================//
//	                            Admin Operations                              //
// ===========================================================================//
func (r *ArticleRepository) AdminListAllArticles(ctx context.Context, pag domain.Pagination) ([]domain.Article, int, error) {
	articles := []domain.Article{}
	opts := options.Find().SetSkip(int64(pag.Page * pag.PageSize)).SetLimit(int64(pag.PageSize))
	cursor, err := r.Collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var articleDTO ArticleListDTO
		if err := cursor.Decode(&articleDTO); err != nil {
			return nil, 0, err
		}
		article := FromArticleListDTO(&articleDTO)
		articles = append(articles, *article)
	}
	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}
	total, err := r.Collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}
	return articles, int(total), nil

}
func (r *ArticleRepository) HardDelete(ctx context.Context, articleID string) error {
	_, err := r.Collection.DeleteOne(ctx, bson.M{"article_id": articleID})
	if err != nil {
		return err
	}
	return nil
}
