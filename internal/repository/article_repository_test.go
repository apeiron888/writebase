package repository_test

import (
	"context"
	"testing"
	"time"
	"write_base/internal/domain"
	"write_base/internal/repository"
	"write_base/tests/test_data"
	"write_base/tests/test_utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ArticleRepoTestSuite struct {
	db         *mongo.Database
	repo       domain.IArticleRepository
	ctx        context.Context
	collection *mongo.Collection
	cleanup    func()
}

func (s *ArticleRepoTestSuite) SetupSuite(t *testing.T) {
	s.ctx = context.Background()
	db, cleanup := test_utils.SetupTestDatabase(t)
	s.db = db
	s.collection = db.Collection("articles")
	s.cleanup = cleanup
	s.repo = repository.NewArticleRepository(db, "articles")

	// Ensure compound index for search tests
	_, err := s.collection.Indexes().CreateOne(s.ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "status", Value: 1}, {Key: "title", Value: "text"}, {Key: "excerpt", Value: "text"}},
		Options: options.Index().SetDefaultLanguage("none").SetLanguageOverride("none").SetName("compound_status_title_excerpt"),
	})
	require.NoError(t, err)
}

func (s *ArticleRepoTestSuite) TearDownSuite(t *testing.T) { s.cleanup() }

// helpers
func (s *ArticleRepoTestSuite) resetCollection(t *testing.T) {
	t.Helper()
	_, _ = s.collection.DeleteMany(s.ctx, bson.M{})
}

func (s *ArticleRepoTestSuite) mustCreateArticle(t *testing.T, a *domain.Article) *domain.Article {
	t.Helper()
	require.NoError(t, s.repo.Create(s.ctx, a))
	got, err := s.repo.GetByID(s.ctx, a.ID)
	require.NoError(t, err)
	return got
}

func (s *ArticleRepoTestSuite) mustPublish(t *testing.T, id string, when time.Time) {
	t.Helper()
	require.NoError(t, s.repo.Publish(s.ctx, id, when))
}

// ========================= Create =========================
func TestArticleRepo_Create(t *testing.T) {
	suite := &ArticleRepoTestSuite{}
	suite.SetupSuite(t)
	defer suite.TearDownSuite(t)
	suite.resetCollection(t)

	for _, tc := range test_data.CreateTests {
		t.Run(tc.Name, func(t *testing.T) {
			err := suite.repo.Create(suite.ctx, tc.Article)
			if tc.ExpectError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.ExpectError)
				return
			}
			assert.NoError(t, err)
			var dto repository.ArticleDTO
			err = suite.collection.FindOne(suite.ctx, bson.M{"_id": tc.Article.ID}).Decode(&dto)
			assert.NoError(t, err)
			result := dto.ToDomain()
			assert.Equal(t, tc.Article.ID, result.ID)
			assert.Equal(t, tc.Article.Title, result.Title)
			assert.Equal(t, tc.Article.AuthorID, result.AuthorID)
		})
	}
}

// ========================= GetByID =========================
func TestArticleRepo_GetByID(t *testing.T) {
	s := &ArticleRepoTestSuite{}
	s.SetupSuite(t)
	defer s.TearDownSuite(t)
	s.resetCollection(t)

	a := &domain.Article{ID: "a1", Title: "Hello", Slug: "hello", AuthorID: "u1", Status: domain.StatusDraft}
	s.mustCreateArticle(t, a)

	got, err := s.repo.GetByID(s.ctx, "a1")
	require.NoError(t, err)
	assert.Equal(t, "Hello", got.Title)
}

// ========================= Update =========================
func TestArticleRepo_Update(t *testing.T) {
	s := &ArticleRepoTestSuite{}
	s.SetupSuite(t)
	defer s.TearDownSuite(t)
	s.resetCollection(t)

	a := &domain.Article{ID: "a2", Title: "Old", Slug: "old", AuthorID: "u1", Status: domain.StatusDraft}
	s.mustCreateArticle(t, a)
	a.Title = "New"
	require.NoError(t, s.repo.Update(s.ctx, a))
	got, err := s.repo.GetByID(s.ctx, a.ID)
	require.NoError(t, err)
	assert.Equal(t, "New", got.Title)
}

// ========================= Delete & Restore =========================
func TestArticleRepo_DeleteRestore(t *testing.T) {
	s := &ArticleRepoTestSuite{}
	s.SetupSuite(t)
	defer s.TearDownSuite(t)
	s.resetCollection(t)

	a := &domain.Article{ID: "a3", Title: "T", Slug: "t", AuthorID: "u1", Status: domain.StatusDraft}
	s.mustCreateArticle(t, a)
	require.NoError(t, s.repo.Delete(s.ctx, a.ID))
	got, err := s.repo.GetByID(s.ctx, a.ID)
	require.NoError(t, err)
	assert.Equal(t, domain.StatusDeleted, got.Status)
	require.NoError(t, s.repo.Restore(s.ctx, a.ID))
	got, err = s.repo.GetByID(s.ctx, a.ID)
	require.NoError(t, err)
	assert.Equal(t, domain.StatusDraft, got.Status)
}

// ========================= Publish/Unpublish/Archive/Unarchive =========================
func TestArticleRepo_StateTransitions(t *testing.T) {
	s := &ArticleRepoTestSuite{}
	s.SetupSuite(t)
	defer s.TearDownSuite(t)
	s.resetCollection(t)

	a := &domain.Article{ID: "a4", Title: "S", Slug: "s", AuthorID: "u1", Status: domain.StatusDraft}
	s.mustCreateArticle(t, a)
	now := time.Now()

	require.NoError(t, s.repo.Publish(s.ctx, a.ID, now))
	got, _ := s.repo.GetByID(s.ctx, a.ID)
	assert.Equal(t, domain.StatusPublished, got.Status)

	require.NoError(t, s.repo.Unpublish(s.ctx, a.ID))
	got, _ = s.repo.GetByID(s.ctx, a.ID)
	assert.Equal(t, domain.StatusDraft, got.Status)

	require.NoError(t, s.repo.Archive(s.ctx, a.ID, now))
	got, _ = s.repo.GetByID(s.ctx, a.ID)
	assert.Equal(t, domain.StatusArchived, got.Status)

	require.NoError(t, s.repo.Unarchive(s.ctx, a.ID))
	got, _ = s.repo.GetByID(s.ctx, a.ID)
	assert.Equal(t, domain.StatusDraft, got.Status)
}

// ========================= ListByAuthor =========================
func TestArticleRepo_ListByAuthor(t *testing.T) {
	s := &ArticleRepoTestSuite{}
	s.SetupSuite(t)
	defer s.TearDownSuite(t)
	s.resetCollection(t)

	// create and publish two articles for u2
	for i := 0; i < 2; i++ {
		id := "la" + string(rune('1'+i))
		a := &domain.Article{ID: id, Title: "T", Slug: id, AuthorID: "u2", Status: domain.StatusDraft}
		s.mustCreateArticle(t, a)
		s.mustPublish(t, id, time.Now())
	}
	arts, total, err := s.repo.ListByAuthor(s.ctx, "u2", domain.Pagination{Page: 1, PageSize: 10})
	require.NoError(t, err)
	assert.Len(t, arts, 2)
	assert.Equal(t, 2, total)
}

// ========================= FindTrending =========================
func TestArticleRepo_FindTrending(t *testing.T) {
	s := &ArticleRepoTestSuite{}
	s.SetupSuite(t)
	defer s.TearDownSuite(t)
	s.resetCollection(t)

	a := &domain.Article{ID: "tr1", Title: "Trend", Slug: "tr1", AuthorID: "u1", Status: domain.StatusDraft}
	s.mustCreateArticle(t, a)
	s.mustPublish(t, a.ID, time.Now())

	arts, total, err := s.repo.FindTrending(s.ctx, 7, domain.Pagination{Page: 1, PageSize: 10})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, 1)
	assert.GreaterOrEqual(t, len(arts), 1)
}

// ========================= FindNewArticles =========================
func TestArticleRepo_FindNewArticles(t *testing.T) {
	s := &ArticleRepoTestSuite{}
	s.SetupSuite(t)
	defer s.TearDownSuite(t)
	s.resetCollection(t)

	a := &domain.Article{ID: "na1", Title: "New", Slug: "na1", AuthorID: "u1", Status: domain.StatusDraft}
	s.mustCreateArticle(t, a)
	s.mustPublish(t, a.ID, time.Now())

	arts, total, err := s.repo.FindNewArticles(s.ctx, domain.Pagination{Page: 1, PageSize: 10})
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, arts, 1)
}

// ========================= FindPopularArticles =========================
func TestArticleRepo_FindPopularArticles(t *testing.T) {
	s := &ArticleRepoTestSuite{}
	s.SetupSuite(t)
	defer s.TearDownSuite(t)
	s.resetCollection(t)

	a := &domain.Article{ID: "pa1", Title: "Popular", Slug: "pa1", AuthorID: "u1", Status: domain.StatusDraft}
	s.mustCreateArticle(t, a)
	s.mustPublish(t, a.ID, time.Now())
	// add some views
	require.NoError(t, s.repo.IncrementView(s.ctx, a.ID))
	require.NoError(t, s.repo.IncrementView(s.ctx, a.ID))

	arts, total, err := s.repo.FindPopularArticles(s.ctx, domain.Pagination{Page: 1, PageSize: 10})
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, arts, 1)
}

// ========================= FilterAuthorArticles =========================
func TestArticleRepo_FilterAuthorArticles(t *testing.T) {
	s := &ArticleRepoTestSuite{}
	s.SetupSuite(t)
	defer s.TearDownSuite(t)
	s.resetCollection(t)

	// create two published with tag "go" and one draft
	a1 := &domain.Article{ID: "fa1", Title: "A1", Slug: "fa1", AuthorID: "ux", Tags: []string{"go"}, Status: domain.StatusDraft}
	a2 := &domain.Article{ID: "fa2", Title: "A2", Slug: "fa2", AuthorID: "ux", Tags: []string{"go", "db"}, Status: domain.StatusDraft}
	a3 := &domain.Article{ID: "fa3", Title: "A3", Slug: "fa3", AuthorID: "ux", Tags: []string{"db"}, Status: domain.StatusDraft}
	s.mustCreateArticle(t, a1)
	s.mustCreateArticle(t, a2)
	s.mustCreateArticle(t, a3)
	s.mustPublish(t, a1.ID, time.Now())
	s.mustPublish(t, a2.ID, time.Now())

	filter := domain.ArticleFilter{Tags: []string{"go"}, Statuses: []domain.ArticleStatus{domain.StatusPublished}}
	arts, total, err := s.repo.FilterAuthorArticles(s.ctx, "ux", filter, domain.Pagination{Page: 1, PageSize: 10})
	require.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, arts, 2)
}

// ========================= Filter =========================
func TestArticleRepo_Filter(t *testing.T) {
	s := &ArticleRepoTestSuite{}
	s.SetupSuite(t)
	defer s.TearDownSuite(t)
	s.resetCollection(t)

	a1 := &domain.Article{ID: "f1", Title: "A1", Slug: "f1", AuthorID: "u1", Tags: []string{"go"}, Status: domain.StatusDraft}
	a2 := &domain.Article{ID: "f2", Title: "A2", Slug: "f2", AuthorID: "u2", Tags: []string{"go", "db"}, Status: domain.StatusDraft}
	a3 := &domain.Article{ID: "f3", Title: "A3", Slug: "f3", AuthorID: "u3", Tags: []string{"db"}, Status: domain.StatusDraft}
	s.mustCreateArticle(t, a1)
	s.mustCreateArticle(t, a2)
	s.mustCreateArticle(t, a3)
	s.mustPublish(t, a1.ID, time.Now())
	s.mustPublish(t, a2.ID, time.Now())

	filter := domain.ArticleFilter{Tags: []string{"go"}, Statuses: []domain.ArticleStatus{domain.StatusPublished}}
	arts, total, err := s.repo.Filter(s.ctx, filter, domain.Pagination{Page: 1, PageSize: 10})
	require.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, arts, 2)
}

// ========================= Search =========================
func TestArticleRepo_Search(t *testing.T) {
	s := &ArticleRepoTestSuite{}
	s.SetupSuite(t)
	defer s.TearDownSuite(t)
	s.resetCollection(t)

	a := &domain.Article{ID: "s1", Title: "Golang Tips", Slug: "s1", AuthorID: "u1", Status: domain.StatusDraft}
	s.mustCreateArticle(t, a)
	s.mustPublish(t, a.ID, time.Now())

	arts, total, err := s.repo.Search(s.ctx, "Golang", domain.Pagination{Page: 1, PageSize: 10})
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, arts, 1)
}

// ========================= ListByTags =========================
func TestArticleRepo_ListByTags(t *testing.T) {
	s := &ArticleRepoTestSuite{}
	s.SetupSuite(t)
	defer s.TearDownSuite(t)
	s.resetCollection(t)

	a1 := &domain.Article{ID: "t1", Title: "T1", Slug: "t1", AuthorID: "u1", Tags: []string{"go"}, Status: domain.StatusDraft}
	a2 := &domain.Article{ID: "t2", Title: "T2", Slug: "t2", AuthorID: "u1", Tags: []string{"go"}, Status: domain.StatusDraft}
	s.mustCreateArticle(t, a1)
	s.mustCreateArticle(t, a2)
	s.mustPublish(t, a1.ID, time.Now())
	s.mustPublish(t, a2.ID, time.Now())

	arts, total, err := s.repo.ListByTags(s.ctx, []string{"go"}, domain.Pagination{Page: 1, PageSize: 10})
	require.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, arts, 2)
}

// ========================= EmptyTrash =========================
func TestArticleRepo_TrashOperations(t *testing.T) {
	s := &ArticleRepoTestSuite{}
	s.SetupSuite(t)
	defer s.TearDownSuite(t)
	s.resetCollection(t)

	// create two, mark deleted
	a1 := &domain.Article{ID: "d1", Title: "D1", Slug: "d1", AuthorID: "ux", Status: domain.StatusDraft}
	a2 := &domain.Article{ID: "d2", Title: "D2", Slug: "d2", AuthorID: "ux", Status: domain.StatusDraft}
	s.mustCreateArticle(t, a1)
	s.mustCreateArticle(t, a2)
	require.NoError(t, s.repo.Delete(s.ctx, a1.ID))
	require.NoError(t, s.repo.Delete(s.ctx, a2.ID))

	// delete single from trash
	require.NoError(t, s.repo.DeleteFromTrash(s.ctx, a1.ID, "ux"))
	// empty remaining trash
	require.NoError(t, s.repo.EmptyTrash(s.ctx, "ux"))
}

// ========================= AdminListAll & HardDelete =========================
func TestArticleRepo_AdminAndHardDelete(t *testing.T) {
	s := &ArticleRepoTestSuite{}
	s.SetupSuite(t)
	defer s.TearDownSuite(t)
	s.resetCollection(t)

	a1 := &domain.Article{ID: "ad1", Title: "A1", Slug: "ad1", AuthorID: "u1", Status: domain.StatusDraft}
	a2 := &domain.Article{ID: "ad2", Title: "A2", Slug: "ad2", AuthorID: "u2", Status: domain.StatusDraft}
	s.mustCreateArticle(t, a1)
	s.mustCreateArticle(t, a2)

	arts, total, err := s.repo.AdminListAllArticles(s.ctx, domain.Pagination{Page: 1, PageSize: 10})
	require.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, arts, 2)

	// hard delete one
	require.NoError(t, s.repo.HardDelete(s.ctx, a1.ID))
	_, err = s.repo.GetByID(s.ctx, a1.ID)
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrArticleNotFound)
}

// ========================= Stats & Counters =========================
func TestArticleRepo_StatsAndCounters(t *testing.T) {
	s := &ArticleRepoTestSuite{}
	s.SetupSuite(t)
	defer s.TearDownSuite(t)
	s.resetCollection(t)

	a := &domain.Article{ID: "sc1", Title: "SC", Slug: "sc1", AuthorID: "u1", Status: domain.StatusDraft}
	s.mustCreateArticle(t, a)
	s.mustPublish(t, a.ID, time.Now())

	// increment view and update clap
	require.NoError(t, s.repo.IncrementView(s.ctx, a.ID))
	require.NoError(t, s.repo.UpdateClapCount(s.ctx, a.ID, 5))

	// verify via GetByID (since GetStats projects differently)
	got, err := s.repo.GetByID(s.ctx, a.ID)
	require.NoError(t, err)
	assert.Equal(t, 1, got.Stats.ViewCount)
	assert.Equal(t, 5, got.Stats.ClapCount)
}

// ========================= GetBySlug =========================
func TestArticleRepo_GetBySlug(t *testing.T) {
	s := &ArticleRepoTestSuite{}
	s.SetupSuite(t)
	defer s.TearDownSuite(t)
	s.resetCollection(t)

	a := &domain.Article{ID: "gs1", Title: "By Slug", Slug: "by-slug", AuthorID: "u1", Status: domain.StatusDraft}
	s.mustCreateArticle(t, a)
	got, err := s.repo.GetBySlug(s.ctx, "by-slug")
	require.NoError(t, err)
	assert.Equal(t, a.ID, got.ID)
}
