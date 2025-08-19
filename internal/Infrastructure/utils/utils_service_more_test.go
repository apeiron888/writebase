package utils

import (
	"testing"
	"write_base/internal/domain"
)

func TestValidateContent_EdgeBlocks(t *testing.T) {
	u := &Utils{}
	good := []domain.ContentBlock{
		{Type: domain.BlockImage, Content: domain.BlockContent{Image: &domain.ImageContent{URL: "u", Alt: "a"}}},
		{Type: domain.BlockCode, Content: domain.BlockContent{Code: &domain.CodeContent{Code: "x", Language: "go"}}},
		{Type: domain.BlockVideoEmbed, Content: domain.BlockContent{VideoEmbed: &domain.VideoEmbedContent{Provider: "yt", URL: "http://x"}}},
		{Type: domain.BlockList, Content: domain.BlockContent{List: &domain.ListContent{Items: []string{"a"}}}},
		{Type: domain.BlockDivider, Content: domain.BlockContent{Divider: &domain.DividerContent{Style: "solid"}}},
	}
	if !u.ValidateContent(good) {
		t.Fatal("expected valid")
	}

	bad := []domain.ContentBlock{
		{Type: "unknown"},
	}
	if u.ValidateContent(bad) {
		t.Fatal("expected invalid")
	}
}
