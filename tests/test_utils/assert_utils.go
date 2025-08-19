package test_utils

import (
	"testing"
	"write_base/internal/domain"

	"github.com/stretchr/testify/assert"
)

// AssertArticleEqual compares two articles while ignoring dynamic fields
func AssertArticleEqual(t *testing.T, expected, actual *domain.Article, msg string) {
	assert.Equal(t, expected.ID, actual.ID, msg+" (ID mismatch)")
	assert.Equal(t, expected.Title, actual.Title, msg+" (Title mismatch)")
	assert.Equal(t, expected.Slug, actual.Slug, msg+" (Slug mismatch)")
	assert.Equal(t, expected.AuthorID, actual.AuthorID, msg+" (AuthorID mismatch)")
	assert.Equal(t, expected.Status, actual.Status, msg+" (Status mismatch)")
	assert.Equal(t, expected.Language, actual.Language, msg+" (Language mismatch)")
	assert.ElementsMatch(t, expected.Tags, actual.Tags, msg+" (Tags mismatch)")

	// Compare content blocks
	assert.Len(t, actual.ContentBlocks, len(expected.ContentBlocks), msg+" (ContentBlocks length mismatch)")
	for i, expBlock := range expected.ContentBlocks {
		actBlock := actual.ContentBlocks[i]
		assert.Equal(t, expBlock.Type, actBlock.Type, msg+" (Block type mismatch at index %d)", i)
		assert.Equal(t, expBlock.Order, actBlock.Order, msg+" (Block order mismatch at index %d)", i)
		
		switch expBlock.Type {
		case domain.BlockParagraph:
			assert.Equal(t, expBlock.Content.Paragraph.Text, actBlock.Content.Paragraph.Text)
		case domain.BlockHeading:
			assert.Equal(t, expBlock.Content.Heading.Text, actBlock.Content.Heading.Text)
			assert.Equal(t, expBlock.Content.Heading.Level, actBlock.Content.Heading.Level)
		case domain.BlockImage:
			assert.Equal(t, expBlock.Content.Image.URL, actBlock.Content.Image.URL)
			assert.Equal(t, expBlock.Content.Image.Alt, actBlock.Content.Image.Alt)
		case domain.BlockCode:
			assert.Equal(t, expBlock.Content.Code.Code, actBlock.Content.Code.Code)
			assert.Equal(t, expBlock.Content.Code.Language, actBlock.Content.Code.Language)
		}
	}
}

// AssertTimestampsValid verifies timestamps are set and logical
func AssertTimestampsValid(t *testing.T, timestamps domain.ArticleTimes, msg string) {
	assert.False(t, timestamps.CreatedAt.IsZero(), msg+" (CreatedAt should be set)")
	assert.False(t, timestamps.UpdatedAt.IsZero(), msg+" (UpdatedAt should be set)")
	assert.True(t, timestamps.UpdatedAt.Equal(timestamps.CreatedAt) || 
		timestamps.UpdatedAt.After(timestamps.CreatedAt), 
		msg+" (UpdatedAt should be after or equal to CreatedAt)")
}