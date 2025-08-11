package repository

import (
	"time"
	"context"
	"write_base/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ArticleRepository struct {
	Collection *mongo.Collection
}

func NewArticleRepository(db *mongo.Database, collection string) domain.IArticleRepository {
	return &ArticleRepository{Collection: db.Collection(collection)}
}
//============================= Generals ============================================
func buildArticleFilterQuery(authorID string, filter domain.ArticleFilter) bson.M {
	q := bson.M{}

	// authorID takes precedence over filter.AuthorIDs
	if authorID != "" {
		q["author_id"] = authorID
	} else if len(filter.AuthorIDs) > 0 {
		q["author_id"] = bson.M{"$in": filter.AuthorIDs}
	}

	// statuses
	if len(filter.Statuses) > 0 {
		statusStrings := make([]string, len(filter.Statuses))
		for i, s := range filter.Statuses {
			statusStrings[i] = string(s)
		}
		q["status"] = bson.M{"$in": statusStrings}
	}

	// tags include
	if len(filter.Tags) > 0 {
		q["tags"] = bson.M{"$all": filter.Tags}
	}

	// exclude tags
	if len(filter.ExcludeTags) > 0 {
		// if tags already set to $all, we need to combine; safest is $and
		if _, ok := q["tags"]; ok {
			q = bson.M{"$and": []bson.M{
				{"tags": q["tags"]},
				{"tags": bson.M{"$nin": filter.ExcludeTags}},
			}}
		} else {
			q["tags"] = bson.M{"$nin": filter.ExcludeTags}
		}
	}

	// language
	if filter.Language != "" {
		q["language"] = filter.Language
	}

	// engagement thresholds
	if filter.MinViews > 0 {
		q["stats.views"] = bson.M{"$gte": filter.MinViews}
	}
	if filter.MinClaps > 0 {
		q["stats.claps"] = bson.M{"$gte": filter.MinClaps}
	}

	// published date ranges
	if filter.PublishedAfter != nil && filter.PublishedBefore != nil {
		q["timestamps.published_at"] = bson.M{
			"$gte": *filter.PublishedAfter,
			"$lte": *filter.PublishedBefore,
		}
	} else if filter.PublishedAfter != nil {
		q["timestamps.published_at"] = bson.M{"$gte": *filter.PublishedAfter}
	} else if filter.PublishedBefore != nil {
		q["timestamps.published_at"] = bson.M{"$lte": *filter.PublishedBefore}
	}

	return q
}


//===============================================================================//
//                                CRUD                                           //
//===============================================================================//
// =============================== Article Create ================================
func (ar *ArticleRepository) Create(ctx context.Context, article *domain.Article) error {
	articleDTO := ToArticleDTO(article)
	if articleDTO == nil {
		return domain.ErrInternalServer 
	}
	if _, err := ar.Collection.InsertOne(ctx, articleDTO); err != nil {
		return domain.ErrInternalServer
	}
	return nil
}
// =============================== Article Update ================================
func (ar *ArticleRepository) Update(ctx context.Context, article *domain.Article) error {
	articleDTO := ToArticleDTO(article)
	if articleDTO == nil {
		return domain.ErrInternalServer
	}
	filter := bson.M{"_id":articleDTO.ID}
	if _, err := ar.Collection.UpdateOne(ctx,filter,bson.M{"$set": articleDTO}); err != nil {
		return domain.ErrInternalServer
	}
	return nil
}
// =============================== Article Delete ================================
func (ar *ArticleRepository) Delete(ctx context.Context, articleID string) error {
	filter := bson.M{"_id":articleID}
	update := bson.M{"$set": bson.M{"status": string(domain.StatusDeleted)}}
	if _, err := ar.Collection.UpdateOne(ctx,filter,update); err != nil {
		return domain.ErrInternalServer
	}
	return nil
}
// =============================== Article Restore ================================
func (ar *ArticleRepository) Restore(ctx context.Context, articleID string) error {
	filter := bson.M{"_id":articleID}
	update := bson.M{"$set": bson.M{"status": string(domain.StatusDraft)}}
	if _, err := ar.Collection.UpdateOne(ctx,filter,update); err != nil {
		return domain.ErrInternalServer
	}
	return nil
}
//===============================================================================//
//                   Article State Management                                    //
//===============================================================================//
// ======================== Article Publish =======================================
func (ar *ArticleRepository) Publish(ctx context.Context, articleID string, publishAt time.Time) error {
	_, err := ar.Collection.UpdateOne(ctx, bson.M{"_id": articleID}, bson.M{"$set": bson.M{"status": string(domain.StatusPublished), "timestamps.published_at": publishAt}})
	if err != nil {
		return err
	}
	return nil
}
// ======================== Article Unpublish =====================================
func (r *ArticleRepository) Unpublish(ctx context.Context, articleID string) error {
	_, err := r.Collection.UpdateOne(ctx, bson.M{"_id": articleID}, bson.M{"$set": bson.M{"status": string(domain.StatusDraft)}})
	if err != nil {
		return err
	}
	return nil
}
// ======================== Article Archive =======================================
func (r *ArticleRepository) Archive(ctx context.Context, articleID string, archiveAt time.Time) error {
	_, err := r.Collection.UpdateOne(ctx, bson.M{"_id": articleID}, bson.M{"$set": bson.M{"status": string(domain.StatusArchived), "timestamps.archived_at": archiveAt}})
	if err != nil {
		return err
	}
	return nil
}
// ======================== Article Unarchive =====================================
func (r *ArticleRepository) Unarchive(ctx context.Context, articleID string) error {
	_, err := r.Collection.UpdateOne(ctx, bson.M{"_id": articleID}, bson.M{"$set": bson.M{"status": string(domain.StatusDraft)}})
	if err != nil {
		return err
	}
	return nil
}
//===============================================================================//
//                               Retrieve                                        //
//===============================================================================//
// =============================== Article GetByID ================================
func (ar *ArticleRepository) GetByID(ctx context.Context, articleID string) (*domain.Article, error) {
	var articleDTO ArticleDTO
	err := ar.Collection.FindOne(ctx, bson.M{"_id": articleID}).Decode(&articleDTO)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrArticleNotFound
		}
		return nil, err
	}
	return FromArticleDTO(&articleDTO), nil
}
// =============================== Article Stats ================================
func (ar *ArticleRepository) GetStats(ctx context.Context, articleID string) (*domain.ArticleStats, error) {
	opts :=  options.FindOne().SetProjection(bson.M{"stats": 1})
	var articleStatsDTO ArticleStatsDTO
	if err := ar.Collection.FindOne(ctx, bson.M{"_id":articleID,"status":domain.StatusPublished},opts).Decode(&articleStatsDTO);err!=nil {
				if err == mongo.ErrNoDocuments {
			return nil, domain.ErrArticleNotFound
		}
		return nil, err
	}
    stats := FromArticleStatsDTO(articleStatsDTO)
    return &stats, nil
}
// =================== All Article Stats of Author ================================
func (ar *ArticleRepository) GetAllArticleStats(ctx context.Context, userID string) ([]domain.ArticleStats,int, error) {
	opts :=  options.Find().SetProjection(bson.M{"stats": 1})
	filter := bson.M{"author_id":userID,"status":domain.StatusPublished}
	articlestats := []domain.ArticleStats{}
	cursor,err := ar.Collection.Find(ctx,filter,opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var articleStatsDTO ArticleStatsDTO
		if err := cursor.Decode(&articleStatsDTO); err!=nil {
			return nil,0,err
		}
		stats := FromArticleStatsDTO(articleStatsDTO)
		articlestats = append(articlestats, stats)
	}
	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}
	total, err := ar.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	if total == 0 {
    return nil, 0, domain.ErrArticleNotFound
	}
	return articlestats, int(total), nil
}
// ========================== Article Get By Slug ==================================
func (ar *ArticleRepository) GetBySlug(ctx context.Context, slug string) (*domain.Article, error) {
	var articleDTO ArticleDTO
	err := ar.Collection.FindOne(ctx, bson.M{"slug": slug}).Decode(&articleDTO)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrArticleNotFound
		}
		return nil, err
	}
	return FromArticleDTO(&articleDTO), nil
}

//===========================================================================//
//                           Article Lists                                   //
//===========================================================================//
//======================= List user articles ==================================
func (ar *ArticleRepository) ListByAuthor(ctx context.Context, authorID string, pag domain.Pagination) ([]domain.Article, int, error) {
	var articles []domain.Article
	
	// Build query
	query := bson.M{"author_id": authorID, "status": domain.StatusPublished}
	
	// Build options with pagination and sorting
	opts := options.Find().
		SetSkip(int64((pag.Page - 1) * pag.PageSize)).
		SetLimit(int64(pag.PageSize))
	
	// Add sorting if specified
	if pag.SortField != "" {
		sortOrder := 1
		if pag.SortOrder == "desc" {
			sortOrder = -1
		}
		opts = opts.SetSort(bson.D{{Key: pag.SortField, Value: sortOrder}})
	} else {
		opts = opts.SetSort(bson.D{{Key: "timestamps.created_at", Value: -1}})
	}
	
	cursor, err := ar.Collection.Find(ctx, query, opts)
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

	total, err := ar.Collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	if total<=0 {
		return nil, 0, domain.ErrArticleNotFound
	}

	return articles, int(total), nil
}
//======================= List Trending articles ==================================
func (ar *ArticleRepository) FindTrending(ctx context.Context, windowDays int, pag domain.Pagination) ([]domain.Article, int, error) {
	var articles []domain.Article
	windowAgo := time.Now().AddDate(0, 0, -windowDays)
	query := bson.M{
		"status":                  string(domain.StatusPublished),
		"timestamps.published_at": bson.M{"$gte": windowAgo},
	}
	
	// Build options with pagination and sorting
	opts := options.Find().
		SetSkip(int64((pag.Page - 1) * pag.PageSize)).
		SetLimit(int64(pag.PageSize)).
		SetSort(bson.D{{Key: "stats.view_count", Value: -1}})
	
	if pag.SortField != "" {
		sortOrder := 1 // ascending by default
		if pag.SortOrder == "desc" {
			sortOrder = -1
		}
		opts = opts.SetSort(bson.D{{Key: pag.SortField, Value: sortOrder}})
	}

	cursor, err := ar.Collection.Find(ctx, query, opts)
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

	total, err := ar.Collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	if total<=0{
		return nil,0,domain.ErrArticleNotFound
	}

	return articles, int(total), nil
}
//======================= List New articles =======================================
func (ar *ArticleRepository) FindNewArticles(ctx context.Context, pag domain.Pagination) ([]domain.Article, int, error) {
    var articles []domain.Article

    filter := bson.M{
        "status": domain.StatusPublished,
    }

    opts := options.Find().
        SetSkip(int64((pag.Page - 1) * pag.PageSize)).
        SetLimit(int64(pag.PageSize)).
        SetSort(bson.D{{Key: "timestamps.published_at", Value: -1}})

    if pag.SortField != "" {
        sortOrder := 1
        if pag.SortOrder == "desc" {
            sortOrder = -1
        }
        opts = opts.SetSort(bson.D{{Key: pag.SortField, Value: sortOrder}})
    }

    cursor, err := ar.Collection.Find(ctx, filter, opts)
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

    total, err := ar.Collection.CountDocuments(ctx, filter)
    if err != nil {
        return nil, 0, err
    }
	if total<=0{
		return nil,0,domain.ErrArticleNotFound
	}

    return articles, int(total), nil
}


//======================= List Popular articles ===================================

func (ar *ArticleRepository) FindPopularArticles(ctx context.Context, pag domain.Pagination) ([]domain.Article, int, error) {
    var articles []domain.Article

    filter := bson.M{
        "status": domain.StatusPublished,
    }

    opts := options.Find().
		SetSkip(int64((pag.Page - 1) * pag.PageSize)).
		SetLimit(int64(pag.PageSize)).
        SetSort(bson.D{{Key: "stats.views_total", Value: -1}})
	
	if pag.SortField != "" {
		sortOrder := 1 // ascending by default
		if pag.SortOrder == "desc" {
			sortOrder = -1
		}
		opts = opts.SetSort(bson.D{{Key: pag.SortField, Value: sortOrder}})
	}

    cursor, err := ar.Collection.Find(ctx, filter, opts)
    if err != nil {
        return nil, 0, err
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var articleDTO ArticleListDTO
        if err := cursor.Decode(&articleDTO); err != nil {
            return nil, 0,  ErrArticleToDTO
        }
        article := FromArticleListDTO(&articleDTO)
        articles = append(articles, *article)
    }

    total, err := ar.Collection.CountDocuments(ctx, filter)
    if err != nil {
        return nil, 0, err
    }
	if total<=0{
		return nil,0,domain.ErrArticleNotFound
	}

    return articles, int(total), nil
}


//======================== Filter Author Repository =================================================
func (ar *ArticleRepository) FilterAuthorArticles(ctx context.Context, authorID string, filter domain.ArticleFilter, pag domain.Pagination) ([]domain.Article, int, error) {
	// Defensive pagination defaults
	if pag.Page < 1 {
		pag.Page = 1
	}
	if pag.PageSize <= 0 {
		pag.PageSize = 20
	}
	// build query (buildArticleFilterQuery should include author_id when authorID != "")
	query := buildArticleFilterQuery(authorID, filter)

	// Build options with pagination and sorting
	opts := options.Find().
		SetSkip(int64((pag.Page - 1) * pag.PageSize)).
		SetLimit(int64(pag.PageSize))

	// Add sorting if specified
	if pag.SortField != "" {
		sortOrder := 1 // ascending by default
		if pag.SortOrder == "desc" {
			sortOrder = -1
		}
		opts = opts.SetSort(bson.D{{Key: pag.SortField, Value: sortOrder}})
	} else {
		// Default sorting by created_at descending
		opts = opts.SetSort(bson.D{{Key: "timestamps.created_at", Value: -1}})
	}

	cursor, err := ar.Collection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var articles []domain.Article
	for cursor.Next(ctx) {
		var articleDTO ArticleListDTO
		if err := cursor.Decode(&articleDTO); err != nil {
			return nil, 0, domain.ErrInternalServer
		}
		article := FromArticleListDTO(&articleDTO)
		if article == nil {
			return nil, 0, domain.ErrInternalServer
		}
		articles = append(articles, *article)
	}
	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	total, err := ar.Collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return nil, 0, domain.ErrArticleNotFound
	}

	return articles, int(total), nil
}
//======================== Filter Repository =================================================
func (ar *ArticleRepository) Filter(ctx context.Context, filter domain.ArticleFilter, pag domain.Pagination) ([]domain.Article, int, error) {
	// Defensive pagination defaults
	if pag.Page < 1 {
		pag.Page = 1
	}
	if pag.PageSize <= 0 {
		pag.PageSize = 20
	}

	// build query; empty author ID means global filter
	query := buildArticleFilterQuery("", filter)

	// Build options with pagination and sorting
	opts := options.Find().
		SetSkip(int64((pag.Page - 1) * pag.PageSize)).
		SetLimit(int64(pag.PageSize))

	// Add sorting if specified
	if pag.SortField != "" {
		sortOrder := 1 // ascending by default
		if pag.SortOrder == "desc" {
			sortOrder = -1
		}
		opts = opts.SetSort(bson.D{{Key: pag.SortField, Value: sortOrder}})
	} else {
		// Default sorting by created_at descending
		opts = opts.SetSort(bson.D{{Key: "timestamps.created_at", Value: -1}})
	}

	cursor, err := ar.Collection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var articles []domain.Article
	for cursor.Next(ctx) {
		var articleDTO ArticleListDTO
		if err := cursor.Decode(&articleDTO); err != nil {
			return nil, 0, domain.ErrInternalServer
		}
		article := FromArticleListDTO(&articleDTO)
		if article == nil {
			return nil, 0, domain.ErrInternalServer
		}
		articles = append(articles, *article)
	}
	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	total, err := ar.Collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return nil, 0, domain.ErrArticleNotFound
	}

	return articles, int(total), nil
}
//=========================== Search ==============================================
func (ar *ArticleRepository) Search(ctx context.Context, query string, pag domain.Pagination) ([]domain.Article, int, error) {
    articles := []domain.Article{}

    filter := bson.M{
        "$or": []bson.M{
            {"title": bson.M{"$regex": query, "$options": "i"}},
            {"content": bson.M{"$regex": query, "$options": "i"}},
        },
    }

    opts := options.Find().
        SetSkip(int64((pag.Page - 1) * pag.PageSize)).
        SetLimit(int64(pag.PageSize))

    // Add sorting if specified
    if pag.SortField != "" {
        sortOrder := 1
        if pag.SortOrder == "desc" {
            sortOrder = -1
        }
        opts = opts.SetSort(bson.D{{Key: pag.SortField, Value: sortOrder}})
    } else {
        // Default sort by title ascending (or any default field)
        opts = opts.SetSort(bson.D{{Key: "title", Value: 1}})
    }

    cursor, err := ar.Collection.Find(ctx, filter, opts)
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

    total, err := ar.Collection.CountDocuments(ctx, filter)
    if err != nil {
        return nil, 0, err
    }
    if total == 0 {
        return nil, 0, domain.ErrArticleNotFound
    }

    return articles, int(total), nil
}
//======================== List By Tags =======================================
func (r *ArticleRepository) ListByTags(ctx context.Context, tags []string, pag domain.Pagination) ([]domain.Article, int, error) {
	query := bson.M{
		"tags":   bson.M{"$in": tags},
		"status": string(domain.StatusPublished),
	}

	opts := options.Find().
		SetSkip(int64((pag.Page - 1) * pag.PageSize)).
		SetLimit(int64(pag.PageSize)).
		SetSort(bson.D{{Key: "timestamps.created_at", Value: -1}})

	cursor, err := r.Collection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var articles []domain.Article
	for cursor.Next(ctx) {
		var dto ArticleListDTO
		if err := cursor.Decode(&dto); err != nil {
			return nil, 0, err
		}
		articles = append(articles, *FromArticleListDTO(&dto))
	}

	total, err := r.Collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	return articles, int(total), nil
}
// ===========================================================================//
//	                           Trash Management                               //
// ===========================================================================//
func (r *ArticleRepository) EmptyTrash(ctx context.Context, userID string) error {
	res, err := r.Collection.DeleteMany(ctx, bson.M{"status": string(domain.StatusDeleted), "author_id": userID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return domain.ErrArticleNotFound
	}
	return nil
}
//=========================== Delete Article From Trash =========================
func (r *ArticleRepository) DeleteFromTrash(ctx context.Context, articleID, userID string) error {
	res, err := r.Collection.DeleteMany(ctx, bson.M{"_id":articleID, "status": string(domain.StatusDeleted), "author_id": userID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return domain.ErrArticleNotFound
	}
	return nil
}

// ===========================================================================//
//                            Admin Operations                                //
// ===========================================================================//
func (r *ArticleRepository) AdminListAllArticles(ctx context.Context, pag domain.Pagination) ([]domain.Article, int, error) {
	articles := []domain.Article{}
	
	// Build options with pagination and sorting
	opts := options.Find().
		SetSkip(int64((pag.Page - 1) * pag.PageSize)).
		SetLimit(int64(pag.PageSize))
	
	// Add sorting if specified
	if pag.SortField != "" {
		sortOrder := 1 // ascending by default
		if pag.SortOrder == "desc" {
			sortOrder = -1
		}
		opts = opts.SetSort(bson.D{{Key: pag.SortField, Value: sortOrder}})
	} else {
		// Default sorting by created_at descending
		opts = opts.SetSort(bson.D{{Key: "timestamps.created_at", Value: -1}})
	}
	
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
//================================ Hard Delete (Admin) ====================================
func (r *ArticleRepository) HardDelete(ctx context.Context, articleID string) error {
	_, err := r.Collection.DeleteOne(ctx, bson.M{"_id": articleID})
	if err != nil {
		return err
	}
	return nil
}


//================================ Increment View ===========================================
func (r *ArticleRepository) IncrementView(ctx context.Context, articleID string) error {
	_, err := r.Collection.UpdateOne(ctx, 
		bson.M{"_id": articleID}, 
		bson.M{"$inc": bson.M{"stats.view_count": 1}})
	return err
}

// Add this method to ArticleRepository
func (r *ArticleRepository) UpdateClapCount(ctx context.Context, articleID string, count int) error {
	_, err := r.Collection.UpdateOne(ctx, 
		bson.M{"_id": articleID}, 
		bson.M{"$set": bson.M{"stats.clap_count": count}})
	return err
}