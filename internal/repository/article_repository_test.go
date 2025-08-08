package repository_test

import (
	"context"
	"testing"
	"time"

	"write_base/internal/domain"
	"write_base/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	testDBName         = "test_db"
	testCollectionName = "test_articles"
)

func setupTestDB(t *testing.T) (*mongo.Database, func()) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err)

	db := client.Database(testDBName)
	// Ensure text index for Search tests
	coll := db.Collection(testCollectionName)
	_, err = coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "title", Value: "text"}, {Key: "excerpt", Value: "text"}},
		Options: options.Index().SetDefaultLanguage("none").SetLanguageOverride("none"),
	})
	require.NoError(t, err)
	cleanup := func() {
		err := db.Drop(ctx)
		assert.NoError(t, err)
		err = client.Disconnect(ctx)
		assert.NoError(t, err)
	}
	return db, cleanup
}

func TestArticleRepository_Insert(t *testing.T) {
	testCases := []struct {
		name        string
		article     *domain.Article
		preInsert   bool // For duplicate test
		expectError bool
	}{
		{
			name: "Success - valid article",
			article: &domain.Article{
				ID:       "article_1",
				Title:    "Test Article",
				Slug:     "test-article",
				AuthorID: "author_1",
				Status:   domain.StatusDraft,
				Timestamps: domain.ArticleTimes{
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			expectError: false,
		},
		{
			name: "Failure - duplicate ID",
			article: &domain.Article{
				ID:       "duplicate_id",
				Title:    "Duplicate Article",
				AuthorID: "author_1",
				Timestamps: domain.ArticleTimes{
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			preInsert:   true,
			expectError: true,
		},
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := repository.NewArticleRepository(db, testCollectionName)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			// Pre-insert for duplicate test
			if tc.preInsert {
				_, err := repo.Insert(ctx, tc.article)
				require.NoError(t, err)
			}

			_, err := repo.Insert(ctx, tc.article)

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				fetched, err := repo.GetByID(ctx, tc.article.ID)
				require.NoError(t, err)
				assert.Equal(t, tc.article.ID, fetched.ID)
			}
		})
	}
}

func TestArticleRepository_Update(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repository.NewArticleRepository(db, testCollectionName)
	ctx := context.Background()

	// Setup
	article := &domain.Article{
		ID:       "article_update",
		Title:    "Original Title",
		AuthorID: "author_1",
		Status:   domain.StatusDraft,
		Timestamps: domain.ArticleTimes{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	_, err := repo.Insert(ctx, article)
	require.NoError(t, err)

	testCases := []struct {
		name        string
		article     *domain.Article
		expectError bool
		expectedErr error
	}{
		{
			name: "Success - update title",
			article: &domain.Article{
				ID:       "article_update",
				Title:    "Updated Title",
				AuthorID: "author_1",
				Status:   domain.StatusDraft,
				Timestamps: domain.ArticleTimes{
					CreatedAt: article.Timestamps.CreatedAt,
					UpdatedAt: time.Now(),
				},
			},
			expectError: false,
		},
		{
			name: "Failure - non-existent ID",
			article: &domain.Article{
				ID:       "non_existent",
				Title:    "Non-existent Article",
				AuthorID: "author_1",
			},
			expectError: true,
			expectedErr: domain.ErrArticleNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := repo.Update(ctx, tc.article)

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				updated, err := repo.GetByID(ctx, tc.article.ID)
				require.NoError(t, err)
				assert.Equal(t, tc.article.Title, updated.Title)
			}
		})
	}
}

func TestArticleRepository_Delete(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repository.NewArticleRepository(db, testCollectionName)
	ctx := context.Background()

	// Setup
	article := &domain.Article{
		ID:       "article_delete",
		Title:    "Article to Delete",
		AuthorID: "author_1",
		Status:   domain.StatusDraft,
		Timestamps: domain.ArticleTimes{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	_, err := repo.Insert(ctx, article)
	require.NoError(t, err)

	testCases := []struct {
		name        string
		articleID   string
		expectError bool
		expectedErr error
	}{
		{
			name:        "Success - delete existing article",
			articleID:   "article_delete",
			expectError: false,
		},
		{
			name:        "Failure - delete non-existent article",
			articleID:   "non_existent",
			expectError: true,
			expectedErr: domain.ErrArticleNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.Delete(ctx, tc.articleID)

			if tc.expectError {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
				deleted, err := repo.GetByID(ctx, tc.articleID)
				require.NoError(t, err)
				assert.Equal(t, domain.StatusDeleted, deleted.Status)
			}
		})
	}
}

func TestArticleRepository_Restore(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repository.NewArticleRepository(db, testCollectionName)
	ctx := context.Background()

	// Setup
	article := &domain.Article{
		ID:       "article_restore",
		Title:    "Article to Restore",
		AuthorID: "author_1",
		Status:   domain.StatusDraft,
		Timestamps: domain.ArticleTimes{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	_, err := repo.Insert(ctx, article)
	require.NoError(t, err)

	err = repo.Delete(ctx, article.ID)
	require.NoError(t, err)

	testCases := []struct {
		name        string
		articleID   string
		expectError bool
		expectedErr error
	}{
		{
			name:        "Success - restore deleted article",
			articleID:   "article_restore",
			expectError: false,
		},
		{
			name:        "Failure - restore non-existent article",
			articleID:   "non_existent",
			expectError: true,
			expectedErr: domain.ErrArticleNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := repo.Restore(ctx, tc.articleID)

			if tc.expectError {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
				restored, err := repo.GetByID(ctx, tc.articleID)
				require.NoError(t, err)
				assert.Equal(t, domain.StatusDraft, restored.Status)
			}
		})
	}
}

func TestArticleRepository_GetByID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := repository.NewArticleRepository(db, testCollectionName)
	ctx := context.Background()

	article := &domain.Article{
		ID:       "getbyid_1",
		Title:    "GetByID Article",
		Slug:     "getbyid-article",
		AuthorID: "author_1",
		Status:   domain.StatusDraft,
		Timestamps: domain.ArticleTimes{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	_, err := repo.Insert(ctx, article)
	require.NoError(t, err)

	t.Run("Success - found", func(t *testing.T) {
		got, err := repo.GetByID(ctx, article.ID)
		require.NoError(t, err)
		assert.Equal(t, article.ID, got.ID)
	})

	t.Run("Failure - not found", func(t *testing.T) {
		_, err := repo.GetByID(ctx, "not_found")
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrArticleNotFound)
	})
}

func TestArticleRepository_GetBySlug(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := repository.NewArticleRepository(db, testCollectionName)
	ctx := context.Background()

	article := &domain.Article{
		ID:       "getbyslug_1",
		Title:    "GetBySlug Article",
		Slug:     "getbyslug-article",
		AuthorID: "author_1",
		Status:   domain.StatusPublished,
		Timestamps: domain.ArticleTimes{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	_, err := repo.Insert(ctx, article)
	require.NoError(t, err)

	t.Run("Success - found", func(t *testing.T) {
		got, err := repo.GetBySlug(ctx, article.Slug)
		require.NoError(t, err)
		assert.Equal(t, article.ID, got.ID)
	})

	t.Run("Failure - not found", func(t *testing.T) {
		_, err := repo.GetBySlug(ctx, "not_found")
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrArticleNotFound)
	})
}

func TestArticleRepository_Publish_Unpublish(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := repository.NewArticleRepository(db, testCollectionName)
	ctx := context.Background()

	article := &domain.Article{
		ID:       "pubunpub_1",
		Title:    "Publish Unpublish Article",
		AuthorID: "author_1",
		Status:   domain.StatusDraft,
		Timestamps: domain.ArticleTimes{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	_, err := repo.Insert(ctx, article)
	require.NoError(t, err)

	t.Run("Publish", func(t *testing.T) {
		err := repo.Publish(ctx, article.ID, time.Now())
		require.NoError(t, err)
		got, err := repo.GetByID(ctx, article.ID)
		require.NoError(t, err)
		assert.Equal(t, domain.StatusPublished, got.Status)
	})

	t.Run("Unpublish", func(t *testing.T) {
		err := repo.Unpublish(ctx, article.ID)
		require.NoError(t, err)
		got, err := repo.GetByID(ctx, article.ID)
		require.NoError(t, err)
		assert.Equal(t, domain.StatusDraft, got.Status)
	})
}

func TestArticleRepository_Archive_Unarchive(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := repository.NewArticleRepository(db, testCollectionName)
	ctx := context.Background()

	article := &domain.Article{
		ID:       "archive_1",
		Title:    "Archive Article",
		AuthorID: "author_1",
		Status:   domain.StatusDraft,
		Timestamps: domain.ArticleTimes{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	_, err := repo.Insert(ctx, article)
	require.NoError(t, err)

	t.Run("Archive", func(t *testing.T) {
		err := repo.Archive(ctx, article.ID, time.Now())
		require.NoError(t, err)
		got, err := repo.GetByID(ctx, article.ID)
		require.NoError(t, err)
		assert.Equal(t, domain.StatusArchived, got.Status)
	})

	t.Run("Unarchive", func(t *testing.T) {
		err := repo.Unarchive(ctx, article.ID)
		require.NoError(t, err)
		got, err := repo.GetByID(ctx, article.ID)
		require.NoError(t, err)
		assert.Equal(t, domain.StatusDraft, got.Status)
	})
}

func TestArticleRepository_ListAuthorArticles(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := repository.NewArticleRepository(db, testCollectionName)
	ctx := context.Background()

	// Insert articles
	articles := []*domain.Article{
		{ID: "a1", AuthorID: "author1", Title: "A1", Status: domain.StatusDraft, Timestamps: domain.ArticleTimes{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
		{ID: "a2", AuthorID: "author1", Title: "A2", Status: domain.StatusDraft, Timestamps: domain.ArticleTimes{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
		{ID: "a3", AuthorID: "author2", Title: "A3", Status: domain.StatusDraft, Timestamps: domain.ArticleTimes{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
	}
	for _, a := range articles {
		_, err := repo.Insert(ctx, a)
		require.NoError(t, err)
	}

	pag := domain.Pagination{Page: 0, PageSize: 10}
	got, total, err := repo.ListAuthorArticles(ctx, "author1", pag)
	require.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, got, 2)
}

func TestArticleRepository_FilterAuthorArticles(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := repository.NewArticleRepository(db, testCollectionName)
	ctx := context.Background()

	// Insert articles
	articles := []*domain.Article{
		{ID: "fa1", AuthorID: "author1", Title: "FA1", Slug: "fa1", Excerpt: "Excerpt FA1", Status: domain.StatusDraft, Tags: []string{"go"}, Language: "en", Timestamps: domain.ArticleTimes{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
		{ID: "fa2", AuthorID: "author1", Title: "FA2", Slug: "fa2", Excerpt: "Excerpt FA2", Status: domain.StatusPublished, Tags: []string{"mongo"}, Language: "en", Timestamps: domain.ArticleTimes{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
		{ID: "fa3", AuthorID: "author2", Title: "FA3", Slug: "fa3", Excerpt: "Excerpt FA3", Status: domain.StatusDraft, Tags: []string{"go"}, Language: "fr", Timestamps: domain.ArticleTimes{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
	}
	for _, a := range articles {
		_, err := repo.Insert(ctx, a)
		require.NoError(t, err)
	}

	pag := domain.Pagination{Page: 0, PageSize: 10}
	// Convert []domain.ArticleStatus to []string
	statusStrings := []string{string(domain.StatusDraft)}
	statuses := make([]domain.ArticleStatus, len(statusStrings))
	for i, s := range statusStrings {
		statuses[i] = domain.ArticleStatus(s)
	}
	filter := domain.ArticleFilter{
		Statuses: statuses,
		Tags:     []string{"go"},
		Language: "en",
	}
	got, total, err := repo.FilterAuthorArticles(ctx, "author1", filter, pag)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, got, 1)
	assert.Equal(t, "fa1", got[0].ID)
}

func TestArticleRepository_ListByAuthor(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := repository.NewArticleRepository(db, testCollectionName)
	ctx := context.Background()

	// Insert published and draft articles
	articles := []*domain.Article{
		{ID: "lba1", AuthorID: "author1", Title: "LBA1", Status: domain.StatusPublished, Timestamps: domain.ArticleTimes{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
		{ID: "lba2", AuthorID: "author1", Title: "LBA2", Status: domain.StatusDraft, Timestamps: domain.ArticleTimes{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
	}
	for _, a := range articles {
		_, err := repo.Insert(ctx, a)
		require.NoError(t, err)
	}

	pag := domain.Pagination{Page: 0, PageSize: 10}
	got, total, err := repo.ListByAuthor(ctx, "author1", pag)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, got, 1)
	assert.Equal(t, "lba1", got[0].ID)
}

func TestArticleRepository_ListTrending(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := repository.NewArticleRepository(db, testCollectionName)
	ctx := context.Background()

	// Insert trending and non-trending articles
	now := time.Now()
	articles := []*domain.Article{
		{ID: "lt1", AuthorID: "author1", Title: "LT1", Status: domain.StatusPublished, Stats: domain.ArticleStats{ViewCount: 100}, Timestamps: domain.ArticleTimes{PublishedAt: ptrTime(now.AddDate(0, 0, -1)), CreatedAt: now, UpdatedAt: now}},
		{ID: "lt2", AuthorID: "author1", Title: "LT2", Status: domain.StatusPublished, Stats: domain.ArticleStats{ViewCount: 50}, Timestamps: domain.ArticleTimes{PublishedAt: ptrTime(now.AddDate(0, 0, -2)), CreatedAt: now, UpdatedAt: now}},
		{ID: "lt3", AuthorID: "author1", Title: "LT3", Status: domain.StatusDraft, Stats: domain.ArticleStats{ViewCount: 200}, Timestamps: domain.ArticleTimes{PublishedAt: ptrTime(now.AddDate(0, 0, -1)), CreatedAt: now, UpdatedAt: now}},
	}
	for _, a := range articles {
		_, err := repo.Insert(ctx, a)
		require.NoError(t, err)
	}

	pag := domain.Pagination{Page: 0, PageSize: 10}
	got, total, err := repo.ListTrending(ctx, pag, 7)
	require.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, got, 2)
	assert.Equal(t, "lt1", got[0].ID)
}

func ptrTime(t time.Time) *time.Time {
	return &t
}

func TestArticleRepository_ListByTag(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := repository.NewArticleRepository(db, testCollectionName)
	ctx := context.Background()

	articles := []*domain.Article{
		{ID: "tag1", AuthorID: "author1", Title: "Tag1", Status: domain.StatusPublished, Tags: []string{"go"}, Timestamps: domain.ArticleTimes{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
		{ID: "tag2", AuthorID: "author1", Title: "Tag2", Status: domain.StatusPublished, Tags: []string{"mongo"}, Timestamps: domain.ArticleTimes{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
	}
	for _, a := range articles {
		_, err := repo.Insert(ctx, a)
		require.NoError(t, err)
	}

	pag := domain.Pagination{Page: 0, PageSize: 10}
	got, total, err := repo.ListByTag(ctx, "go", pag)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, got, 1)
	assert.Equal(t, "tag1", got[0].ID)
}

func TestArticleRepository_Search(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := repository.NewArticleRepository(db, testCollectionName)
	ctx := context.Background()

	articles := []*domain.Article{
		{ID: "search1", AuthorID: "author1", Title: "Go Mongo", Status: domain.StatusPublished, Timestamps: domain.ArticleTimes{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
		{ID: "search2", AuthorID: "author1", Title: "Python", Status: domain.StatusPublished, Timestamps: domain.ArticleTimes{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
	}
	for _, a := range articles {
		_, err := repo.Insert(ctx, a)
		require.NoError(t, err)
	}

	pag := domain.Pagination{Page: 0, PageSize: 10}
	got, total, err := repo.Search(ctx, "Go", pag)
	require.NoError(t, err)
	assert.True(t, total >= 1)
	assert.True(t, len(got) >= 1)
}

func TestArticleRepository_Filter(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := repository.NewArticleRepository(db, testCollectionName)
	ctx := context.Background()

	articles := []*domain.Article{
		{ID: "filter1", AuthorID: "author1", Title: "Filter1", Status: domain.StatusPublished, Language: "en", Timestamps: domain.ArticleTimes{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
		{ID: "filter2", AuthorID: "author1", Title: "Filter2", Status: domain.StatusPublished, Language: "fr", Timestamps: domain.ArticleTimes{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
	}
	for _, a := range articles {
		_, err := repo.Insert(ctx, a)
		require.NoError(t, err)
	}

	pag := domain.Pagination{Page: 0, PageSize: 10}
	filter := domain.ArticleFilter{
		Language: "en",
	}
	got, total, err := repo.Filter(ctx, filter, pag)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, got, 1)
	assert.Equal(t, "filter1", got[0].ID)
}
