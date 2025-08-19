package test_utils

import (
	"time"
	"write_base/internal/domain"
)

// CreateTestArticle creates a valid article with optional overrides
func CreateTestArticle(overrides ...func(*domain.Article)) *domain.Article {
	article := &domain.Article{
		ID:       "article_" + time.Now().Format("20060102150405"),
		Title:    "Test Article",
		Slug:     "test-article",
		AuthorID: "author_123",
		ContentBlocks: []domain.ContentBlock{
			{
				Type:  domain.BlockParagraph,
				Order: 1,
				Content: domain.BlockContent{
					Paragraph: &domain.ParagraphContent{
						Text: "Test content paragraph",
					},
				},
			},
		},
		Excerpt:  "Test excerpt",
		Language: "en",
		Tags:     []string{"test", "golang"},
		Status:   domain.StatusDraft,
		Timestamps: domain.ArticleTimes{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, override := range overrides {
		override(article)
	}
	return article
}

// CreateTestArticleWithStatus creates article with specific status
func CreateTestArticleWithStatus(status domain.ArticleStatus) *domain.Article {
	return CreateTestArticle(func(a *domain.Article) {
		a.Status = status
	})
}

// CreateTestArticleWithBlocks creates article with specific content blocks
func CreateTestArticleWithBlocks(blocks []domain.ContentBlock) *domain.Article {
	return CreateTestArticle(func(a *domain.Article) {
		a.ContentBlocks = blocks
	})
}

// CreateTestArticles creates multiple test articles
func CreateTestArticles(count int) []*domain.Article {
	articles := make([]*domain.Article, count)
	for i := 0; i < count; i++ {
		articles[i] = CreateTestArticle(func(a *domain.Article) {
			a.ID = a.ID + "_" + string(rune(i))
			a.Slug = a.Slug + "_" + string(rune(i))
		})
	}
	return articles
}
// In test_utils/article_utils.go
func CreateTestParagraphBlock(text string) domain.ContentBlock {
	return domain.ContentBlock{
		Type: domain.BlockParagraph,
		Order: 1,
		Content: domain.BlockContent{
			Paragraph: &domain.ParagraphContent{
				Text: text,
			},
		},
	}
}

func CreateTestHeadingBlock(text string, level int) domain.ContentBlock {
	return domain.ContentBlock{
		Type: domain.BlockHeading,
		Order: 1,
		Content: domain.BlockContent{
			Heading: &domain.HeadingContent{
				Text: text,
				Level: level,
			},
		},
	}
}