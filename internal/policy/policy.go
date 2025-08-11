// UserExists:  Implement User check by ID
// ArticleCreateValid: Implement Tag existence check

package policy

import (
	"write_base/internal/domain"
)

type Policy struct {
	Utils domain.IUtils
}

func NewArticlePolicy(utils domain.IUtils)domain.IPolicy{
	return &Policy{Utils: utils}
}

func (p *Policy) UserExists(userID string) bool {
	return true
}

func (p *Policy) ArticleCreateValid(input *domain.Article) bool {
	if input.Title == "" || len(input.Title) > domain.MaxTitleLength {
		return false
    }
    if len(input.ContentBlocks) == 0 || len(input.ContentBlocks) > domain.MaxContentBlocks {
		return false
    }
    if !p.Utils.ValidateContent(input.ContentBlocks) {
		return false
    }
    if len(input.Tags) == 0 || len(input.Tags) > domain.MaxTagsPerArticle {
        return false
    }
    return true
}

func (p *Policy) UserOwnsArticle(userID string, input *domain.Article) bool {
  return input.AuthorID == userID
}

func (p *Policy) CheckArticleChangesAndValid(oldArticle *domain.Article, newArticle *domain.Article) bool {
	if newArticle.AuthorID != oldArticle.AuthorID {
		return false
	}
	if newArticle.Language == "" {
		newArticle.Language = "en"
	}
	if len(newArticle.Tags) > 5 {
		return false
	}
	if newArticle.Status != domain.StatusDraft {
		newArticle.Status = domain.StatusDraft
	}
	if oldArticle.Title != newArticle.Title || oldArticle.Excerpt != newArticle.Excerpt || oldArticle.Language != newArticle.Language || oldArticle.Slug != newArticle.Slug {
		return true
	}
	return false
}

func (p *Policy) IsAdmin(userID string, userRole string) bool {
	// Implement policy logic for checking if a user is an admin
	return userRole == "admin" && p.UserExists(userID)
}
