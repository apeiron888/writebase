package utils

import (
	"testing"
	"write_base/internal/domain"
)

func TestValidateContent_ValidBlocks(t *testing.T) {
	u := &Utils{}
	blocks := []domain.ContentBlock{
		{
			Type:    domain.BlockParagraph,
			Order:   0,
			Content: domain.BlockContent{Paragraph: &domain.ParagraphContent{Text: "hello"}},
		},
		{
			Type:    domain.BlockHeading,
			Order:   1,
			Content: domain.BlockContent{Heading: &domain.HeadingContent{Text: "Head", Level: 2}},
		},
	}
	if !u.ValidateContent(blocks) {
		t.Fatalf("expected valid content")
	}
}

func TestValidateContent_InvalidBlock(t *testing.T) {
	u := &Utils{}
	blocks := []domain.ContentBlock{
		{Type: domain.BlockParagraph, Order: 0, Content: domain.BlockContent{Paragraph: &domain.ParagraphContent{Text: ""}}},
	}
	if u.ValidateContent(blocks) {
		t.Fatalf("expected invalid content")
	}
}
