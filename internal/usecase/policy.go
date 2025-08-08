package usecase

import "write_base/internal/domain"

type ArticlePolicy struct{}

func NewArticlePolicy() domain.IArticlePolicy {
	return &ArticlePolicy{}
}

func (p *ArticlePolicy) UserExists(userID string) bool {
	// Implement policy logic for checking user existence
	return true // Placeholder, implement actual logic
}

func (p *ArticlePolicy) UserOwnsArticle(userID string, article domain.Article) bool {
	// Implement policy logic for checking article ownership
	return article.AuthorID == userID && p.UserExists(userID)
}

func (p *ArticlePolicy) CanViewByID(userID string, userRole string, article domain.Article) bool {
	// Implement policy logic for viewing an article by ID
	if  p.UserExists(userID) && article.Status == domain.StatusPublished {
		return true
	}
	if  p.UserExists(userID) && (userRole == "admin" || userRole == "author") {
		return true
	}
	return false
}

func (p *ArticlePolicy) IsAdmin(userID string, userRole string) bool {
	// Implement policy logic for checking if a user is an admin
	return userRole == "admin" && p.UserExists(userID)
}