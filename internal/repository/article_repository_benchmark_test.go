package repository_test

import (
	"context"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"write_base/internal/domain"
	"write_base/internal/repository"
	"write_base/tests/test_utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Benchmark setup
type benchEnv struct {
	db         *mongo.Database
	coll       *mongo.Collection
	repo       domain.IArticleRepository
	totalDocs  int
	authorPick []string
}

func setupBench(b *testing.B) *benchEnv {
	b.Helper()
	db := test_utils.SetupBenchmarkDatabase(b)
	coll := db.Collection("bench_articles")
	repo := repository.NewArticleRepository(db, "bench_articles")

	// Optional: create indexes controlled via env
	// Back-compat: BENCH_NO_INDEX=1 -> no indexes
	// New: BENCH_INDEX_MODE in {none, text, full}. Defaults to full when unset.
	indexMode := os.Getenv("BENCH_INDEX_MODE")
	if os.Getenv("BENCH_NO_INDEX") == "1" {
		indexMode = "none"
	}
	switch indexMode {
	case "none":
		// no indexes
	case "text":
		// only the compound text index required for Search
		_, _ = coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
			Keys:    bson.D{{Key: "status", Value: 1}, {Key: "title", Value: "text"}, {Key: "excerpt", Value: "text"}},
			Options: options.Index().SetDefaultLanguage("none").SetLanguageOverride("none").SetName("text_title_excerpt"),
		})
	default: // "full"
		idxes := []mongo.IndexModel{
			// author + status lists sorted by created_at
			{Keys: bson.D{{Key: "author_id", Value: 1}, {Key: "status", Value: 1}, {Key: "timestamps.created_at", Value: -1}}},
			// new articles by published_at
			{Keys: bson.D{{Key: "status", Value: 1}, {Key: "timestamps.published_at", Value: -1}}},
			// popular by view_count
			{Keys: bson.D{{Key: "status", Value: 1}, {Key: "stats.view_count", Value: -1}}},
			// trending ESR: equality(status), sort(view_count), range(published_at)
			{Keys: bson.D{{Key: "status", Value: 1}, {Key: "stats.view_count", Value: -1}, {Key: "timestamps.published_at", Value: -1}}},
			// tags filter + created_at sort under status equality
			{Keys: bson.D{{Key: "status", Value: 1}, {Key: "tags", Value: 1}, {Key: "timestamps.created_at", Value: -1}}},
			// slug lookup
			{Keys: bson.D{{Key: "slug", Value: 1}}},
		}
		_, _ = coll.Indexes().CreateMany(context.Background(), idxes)
		// text index (compound with status)
		_, _ = coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
			Keys:    bson.D{{Key: "status", Value: 1}, {Key: "title", Value: "text"}, {Key: "excerpt", Value: "text"}},
			Options: options.Index().SetDefaultLanguage("none").SetLanguageOverride("none").SetName("text_title_excerpt"),
		})
	}

	// Seed size via env (default 5000)
	seedN := 5000
	if v := os.Getenv("BENCH_SEED_COUNT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			seedN = n
		}
	}

	// Seed data (only once per run if empty)
	cnt, _ := coll.CountDocuments(context.Background(), bson.M{})
	if cnt == 0 {
		seedBenchmarkArticles(b, coll, seedN)
	}

	return &benchEnv{
		db:        db,
		coll:      coll,
		repo:      repo,
		totalDocs: seedN,
		authorPick: func() []string {
			ids := make([]string, 100)
			for i := range ids {
				ids[i] = "user_" + strconv.Itoa(i)
			}
			return ids
		}(),
	}
}

func teardownBench(b *testing.B, env *benchEnv) {
	b.Helper()
	test_utils.CleanupBenchmarkDatabase(b, env.db)
}

func seedBenchmarkArticles(b *testing.B, coll *mongo.Collection, n int) {
	b.Helper()
	rng := rand.New(rand.NewSource(42))
	tagsList := []string{"go", "db", "web", "ai", "cli"}
	authors := make([]string, 100)
	for i := range authors {
		authors[i] = "user_" + strconv.Itoa(i)
	}

	// batch insert for speed
	batch := make([]interface{}, 0, 1000)
	now := time.Now()
	for i := 0; i < n; i++ {
		id := "art_" + strconv.Itoa(i)
		author := authors[rng.Intn(len(authors))]
		// 70% published, 25% draft, 5% archived
		stPick := rng.Intn(100)
		status := string(domain.StatusDraft)
		var publishedAt *time.Time
		if stPick < 70 {
			status = string(domain.StatusPublished)
			daysAgo := rng.Intn(30)
			pa := now.AddDate(0, 0, -daysAgo)
			publishedAt = &pa
		} else if stPick >= 95 {
			status = string(domain.StatusArchived)
		}
		// tags 1-3
		tgCount := 1 + rng.Intn(3)
		tg := make([]string, tgCount)
		for j := 0; j < tgCount; j++ {
			tg[j] = tagsList[rng.Intn(len(tagsList))]
		}

		viewCount := rng.Intn(1000)
		clapCount := rng.Intn(300)

		dto := repository.ArticleDTO{
			ID:            id,
			Title:         "Golang DB Tricks " + strconv.Itoa(i),
			Slug:          "golang-db-tricks-" + strconv.Itoa(i),
			AuthorID:      author,
			ContentBlocks: []repository.ContentBlockDTO{},
			Excerpt:       "Excerpt for benchmark article",
			Language:      "en",
			Tags:          tg,
			Status:        status,
			Stats:         repository.ArticleStatsDTO{ViewsCount: viewCount, ClapCount: clapCount},
			Timestamps: repository.ArticleTimesDTO{
				CreatedAt:   now.Add(-time.Duration(rng.Intn(60)) * time.Hour),
				UpdatedAt:   now,
				PublishedAt: publishedAt,
			},
		}
		batch = append(batch, dto)
		if len(batch) == cap(batch) {
			_, _ = coll.InsertMany(context.Background(), batch)
			batch = batch[:0]
		}
	}
	if len(batch) > 0 {
		_, _ = coll.InsertMany(context.Background(), batch)
	}
}

// Benchmarks

func BenchmarkRepo_ListByAuthor(b *testing.B) {
	env := setupBench(b)
	defer teardownBench(b, env)

	author := env.authorPick[0]
	pag := domain.Pagination{Page: 1, PageSize: 20}
	b.ReportAllocs()
	b.ResetTimer()
	var totalDocs int
	for i := 0; i < b.N; i++ {
		res, _, err := env.repo.ListByAuthor(context.Background(), author, pag)
		if err != nil {
			b.Fatalf("ListByAuthor error: %v", err)
		}
		totalDocs += len(res)
	}
	b.StopTimer()
	if b.N > 0 {
		b.ReportMetric(float64(totalDocs)/float64(b.N), "docs/op")
	}
}

func BenchmarkRepo_Filter_TagsPublished(b *testing.B) {
	env := setupBench(b)
	defer teardownBench(b, env)

	filter := domain.ArticleFilter{Tags: []string{"go"}, Statuses: []domain.ArticleStatus{domain.StatusPublished}}
	pag := domain.Pagination{Page: 1, PageSize: 20}
	b.ReportAllocs()
	b.ResetTimer()
	var totalDocs int
	for i := 0; i < b.N; i++ {
		res, _, err := env.repo.Filter(context.Background(), filter, pag)
		if err != nil {
			b.Fatalf("Filter error: %v", err)
		}
		totalDocs += len(res)
	}
	b.StopTimer()
	if b.N > 0 {
		b.ReportMetric(float64(totalDocs)/float64(b.N), "docs/op")
	}
}

func BenchmarkRepo_Search_TitleRegex(b *testing.B) {
	env := setupBench(b)
	defer teardownBench(b, env)

	pag := domain.Pagination{Page: 1, PageSize: 20}
	b.ReportAllocs()
	b.ResetTimer()
	var totalDocs int
	for i := 0; i < b.N; i++ {
		res, _, err := env.repo.Search(context.Background(), "Golang", pag)
		if err != nil {
			b.Fatalf("Search error: %v", err)
		}
		totalDocs += len(res)
	}
	b.StopTimer()
	if b.N > 0 {
		b.ReportMetric(float64(totalDocs)/float64(b.N), "docs/op")
	}
}

func BenchmarkRepo_FindTrending(b *testing.B) {
	env := setupBench(b)
	defer teardownBench(b, env)

	pag := domain.Pagination{Page: 1, PageSize: 20}
	b.ReportAllocs()
	b.ResetTimer()
	var totalDocs int
	for i := 0; i < b.N; i++ {
		res, _, err := env.repo.FindTrending(context.Background(), 7, pag)
		if err != nil {
			b.Fatalf("FindTrending error: %v", err)
		}
		totalDocs += len(res)
	}
	b.StopTimer()
	if b.N > 0 {
		b.ReportMetric(float64(totalDocs)/float64(b.N), "docs/op")
	}
}

func BenchmarkRepo_FindPopular(b *testing.B) {
	env := setupBench(b)
	defer teardownBench(b, env)

	pag := domain.Pagination{Page: 1, PageSize: 20}
	b.ReportAllocs()
	b.ResetTimer()
	var totalDocs int
	for i := 0; i < b.N; i++ {
		res, _, err := env.repo.FindPopularArticles(context.Background(), pag)
		if err != nil {
			b.Fatalf("FindPopularArticles error: %v", err)
		}
		totalDocs += len(res)
	}
	b.StopTimer()
	if b.N > 0 {
		b.ReportMetric(float64(totalDocs)/float64(b.N), "docs/op")
	}
}

func BenchmarkRepo_Create(b *testing.B) {
	env := setupBench(b)
	defer teardownBench(b, env)

	// Creating new docs each iteration (keep simple unique IDs)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := "bench_new_" + strconv.Itoa(i)
		a := &domain.Article{
			ID:       id,
			Title:    "Bench Create " + id,
			Slug:     id,
			AuthorID: "user_0",
			Status:   domain.StatusDraft,
		}
		if err := env.repo.Create(context.Background(), a); err != nil {
			b.Fatalf("Create error: %v", err)
		}
	}
}
