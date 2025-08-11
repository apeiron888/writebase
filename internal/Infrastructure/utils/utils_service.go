package utils

import "write_base/internal/domain"

type Utils struct{}

func NewUtils() domain.IUtils { return &Utils{} }

func (u *Utils) ValidateContent(blocks []domain.ContentBlock) bool {

	for _, block := range blocks {
		switch block.Type {
		case "heading":
			if block.Content.Heading.Text == "" || len(block.Content.Heading.Text) > domain.MaxTitleLength || block.Content.Heading.Level < 1 || block.Content.Heading.Level > 6 {
				return false
			}
		case "paragraph":
			if block.Content.Paragraph.Text == "" || len(block.Content.Paragraph.Text) > domain.MaxContentLength {
				return false
			}
		case "image":
			if block.Content.Image.URL == "" || block.Content.Image.Alt == "" {
				return false
			}
		case "code":
			if block.Content.Code.Code == "" || block.Content.Code.Language == "" {
				return false
			}
		case "video_embed":
			if block.Content.VideoEmbed.Provider == "" || block.Content.VideoEmbed.URL == "" {
				return false
			}
		case "list":
			if len(block.Content.List.Items) == 0 {
				return false
			}
		case "divider":
			if block.Content.Divider.Style == "" {
				return false
			}
		default:
			return false
		}
	}
	return true
}