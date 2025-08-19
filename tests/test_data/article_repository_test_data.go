package test_data

import (
	"write_base/internal/domain"
)

// ArticleCreateTest represents a test case for Create method
type ArticleCreateTest struct {
	Name        string
	Article     *domain.Article
	ExpectError error
	Description string
}

var CreateTests = []ArticleCreateTest{
	{
		Name: "Valid Article",
		Article: &domain.Article{
			ID:       "article1",
			Title:    "Test Article",
			Slug:     "test-article",
			AuthorID: "user1",
			ContentBlocks: []domain.ContentBlock{
				{
					Type:  domain.BlockParagraph,
					Order: 1,
					Content: domain.BlockContent{
						Paragraph: &domain.ParagraphContent{
							Text: "First paragraph",
						},
					},
				},
			},
			Status: domain.StatusDraft,
		},
		ExpectError: nil,
	},
	{
		Name: "Duplicate ID",
		Article: &domain.Article{
			ID:       "article1", // Same as first test
			Title:    "Duplicate Article",
			Slug:     "duplicate-article",
			AuthorID: "user1",
			ContentBlocks: []domain.ContentBlock{
				{
					Type:  domain.BlockParagraph,
					Order: 1,
					Content: domain.BlockContent{
						Paragraph: &domain.ParagraphContent{
							Text: "This should fail",
						},
					},
				},
			},
		},
		ExpectError: domain.ErrInternalServer,
	},
	{
		Name: "Valid Content Block Types",
		Article: &domain.Article{
			ID:       "article6",
			Title:    "Valid Block Types",
			Slug:     "valid-blocks",
			AuthorID: "user1",
			ContentBlocks: []domain.ContentBlock{
				{
					Type: domain.BlockHeading,
					Order: 1,
					Content: domain.BlockContent{
						Heading: &domain.HeadingContent{
							Text: "Heading", 
							Level: 2,
						},
					},
				},
			},
		},
		ExpectError: nil,
	},
}