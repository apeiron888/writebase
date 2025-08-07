package usecase

import (
	"context"
	"strings"
	"time"
	"write_base/internal/domain"
	"github.com/google/uuid"
)

type ArticleUsecase struct {
	ArticleRepo domain.IArticleRepository
	Policy      domain.IArticlePolicy
}

func NewArticleUsecase(repo domain.IArticleRepository, policy domain.IArticlePolicy) domain.IArticleUsecase {
	return &ArticleUsecase{ArticleRepo: repo, Policy: policy}
}

//===========================================================================//
//                       General Functions                                   //
//===========================================================================//
func GenerateSlug(title string) string {
	parts := strings.Fields(title)
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, "-")
}

func GenerateID() string {
	return uuid.NewString()
}

func ValidateContent(blocks []domain.ContentBlock) bool {

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

func CheckArticleChangesAndValid(oldArticle *domain.Article, newArticle *domain.Article) bool {
	if newArticle.Title == "" || len(newArticle.Tags) == 0 || len(newArticle.ContentBlocks) == 0 || newArticle.Excerpt == "" || len(newArticle.Title) > domain.MaxTitleLength || len(newArticle.Excerpt) > domain.MaxContentLength {
		return false
	}
	if newArticle.AuthorID != oldArticle.AuthorID {
		return false
	}
	if newArticle.Language == "" {
		newArticle.Language = "en"
		// AI can be used to detect language
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
	// Check if the content blocks have changed
	if !ValidateContent(newArticle.ContentBlocks) {
		return true
	}
	return false
}

//===========================================================================//
//                 Article Lifecycle (Author only)                           //
//===========================================================================//
func (u *ArticleUsecase) CreateArticle(ctx context.Context, UserID string, input *domain.Article) (*domain.Article, error) {
	c, close := context.WithTimeout(ctx, domain.DefaultTimeout)
	defer close()
	// Validation
	if !u.Policy.UserExists(UserID) {
		return nil, domain.ErrUnauthorized
	}
	if input.Title == "" || len(input.Tags) == 0 || len(input.ContentBlocks) == 0 || input.Excerpt == "" || len(input.Title) > domain.MaxTitleLength || len(input.Excerpt) > domain.MaxContentLength {
		return nil, domain.ErrInvalidRequest
	}
	if input.AuthorID != UserID {
		return nil, domain.ErrUnauthorizedArticleEdit
	}
	if input.Language == "" {
		input.Language = "en" 
		// AI can be used to detect language
	}
	if len(input.Tags) > 5 {
		return nil, domain.ErrArticleTagLimitExceeded
	}
	if input.Status != domain.StatusDraft {
		input.Status = domain.StatusDraft
	}
	if !ValidateContent(input.ContentBlocks) {
		return nil, domain.ErrInvalidRequest
	}

	// Generation and Initialization
	if input.Slug == "" {
		input.Slug = GenerateSlug(input.Title)
	}
	input.ID = GenerateID()
	input.Timestamps.CreatedAt = time.Now()
	input.Timestamps.UpdatedAt = time.Now()
	input.Stats.ClapCount = 0
	input.Stats.ViewCount = 0
	article,err := u.ArticleRepo.Insert(c, input)
	if err != nil {
		if err != domain.ErrArticleAlreadyExists{
			return nil, domain.ErrInternalServer
		}
		return nil, err
	}
	return article,nil
}

func (u *ArticleUsecase) UpdateArticle(ctx context.Context, userID string, input *domain.Article) (*domain.Article, error) {
	/*
		Add => Apply optimistic locking (version check)
	*/
	c,close := context.WithTimeout(ctx, domain.DefaultTimeout)
	defer close()
	if !u.Policy.UserExists(userID) {
		return nil, domain.ErrUnauthorized
	}
	if input.AuthorID != userID {
		return nil, domain.ErrUnauthorizedArticleEdit
	}
	article, err := u.ArticleRepo.GetByID(c, input.ID)
	if err != nil {
		if err == domain.ErrArticleNotFound {
			return nil, domain.ErrArticleNotFound
		}
		return nil, domain.ErrInternalServer
	}
	if !u.Policy.UserOwnsArticle(userID, *article) {
		return nil, domain.ErrUnauthorizedArticleEdit
	}
	if !CheckArticleChangesAndValid(article, input) {
		return nil, domain.ErrNoChangesDetected
	}
	if input.Slug == "" {
		input.Slug = GenerateSlug(input.Title)
	}	
	input.Timestamps.UpdatedAt = time.Now()
	res,err := u.ArticleRepo.Update(c, input)
	if err != nil {
		if err == domain.ErrArticleAlreadyExists {
			return nil, domain.ErrDuplicateArticleSlug
		}
		return nil, domain.ErrInternalServer
	}
	return res, nil
}

func (u *ArticleUsecase) DeleteArticle(ctx context.Context, userID, articleID string) error {
	c, close := context.WithTimeout(ctx, domain.DefaultTimeout)
	defer close()
	if !u.Policy.UserExists(userID) {
		return domain.ErrUnauthorized
	}
	if !u.Policy.UserOwnsArticle(userID, domain.Article{ID: articleID}) {
		return domain.ErrUnauthorizedArticleDelete
	}
	_, err := u.ArticleRepo.GetByID(c, articleID)
	if err != nil {
		if err == domain.ErrArticleNotFound {
			return domain.ErrArticleNotFound
		}
		return domain.ErrInternalServer
	}
	if err = u.ArticleRepo.Delete(c, articleID); err != nil {
		if err == domain.ErrArticleNotFound {
			return domain.ErrArticleNotFound
		}
		return domain.ErrInternalServer
	}
	return nil
}

func (u *ArticleUsecase) RestoreArticle(ctx context.Context, userID, articleID string) (*domain.Article, error) {
	c, close := context.WithTimeout(ctx, domain.DefaultTimeout)
	defer close()
	if !u.Policy.UserExists(userID) {
		return nil, domain.ErrUnauthorized
	}
	if !u.Policy.UserOwnsArticle(userID, domain.Article{ID: articleID}) {
		return nil, domain.ErrUnauthorizedArticleEdit
	}
	_, err := u.ArticleRepo.GetByID(c, articleID)
	if err != nil {
		if err == domain.ErrArticleNotFound {
			return nil, domain.ErrArticleNotFound
		}
		return nil, domain.ErrInternalServer
	}
	res, err := u.ArticleRepo.Restore(c, articleID)
	if err != nil {
		return nil, domain.ErrInternalServer
	}
	return res, nil
}

//===========================================================================//
//                 Article Statistics (Author only)                           //
//===========================================================================//
func (u *ArticleUsecase) GetArticleStats(ctx context.Context, articleID string) (domain.ArticleStats, error) {
	c, close := context.WithTimeout(ctx, domain.DefaultTimeout)
	defer close()
	article, err := u.ArticleRepo.GetByID(c, articleID)
	if err != nil {
		if err == domain.ErrArticleNotFound {
			return domain.ArticleStats{}, domain.ErrArticleNotFound
		}
		return domain.ArticleStats{}, domain.ErrInternalServer
	}
	return article.Stats, nil
}

//===========================================================================//
//             Article State Management (Author only)                        //
//===========================================================================//
func (u *ArticleUsecase) PublishArticle(ctx context.Context, userID, articleID string) (*domain.Article, error) {
	c, close := context.WithTimeout(ctx, domain.DefaultTimeout)
	defer close()
	if !u.Policy.UserExists(userID) {
		return nil, domain.ErrUnauthorized
	}
	if !u.Policy.UserOwnsArticle(userID, domain.Article{ID: articleID}) {
		return nil, domain.ErrUnauthorizedArticleEdit
	}
	article, err := u.ArticleRepo.GetByID(c, articleID)
	if err != nil {
		if err == domain.ErrArticleNotFound {
			return nil, domain.ErrArticleNotFound
		}
		return nil, domain.ErrInternalServer
	}
	article.Status = domain.StatusPublished
	now := time.Now()
	article.Timestamps.PublishedAt = &now
	if err = u.ArticleRepo.Publish(c, articleID, now); err != nil {
		if err == domain.ErrArticleAlreadyPublished {
			return nil, domain.ErrArticleAlreadyPublished
		}
		return nil, domain.ErrInternalServer
	}
	return article, nil
}

func (u *ArticleUsecase) UnpublishArticle(ctx context.Context, userID, articleID string, force bool) (*domain.Article, error) {
	c, close := context.WithTimeout(ctx, domain.DefaultTimeout)
	defer close()
	if !u.Policy.UserExists(userID) {
		return nil, domain.ErrUnauthorized
	}
	if !u.Policy.UserOwnsArticle(userID, domain.Article{ID: articleID}) {
		return nil, domain.ErrUnauthorizedArticleEdit
	}
	article, err := u.ArticleRepo.GetByID(c, articleID)
	if err != nil {
		if err == domain.ErrArticleNotFound {
			return nil, domain.ErrArticleNotFound
		}
		return nil, domain.ErrInternalServer
	}
	article.Status = domain.StatusDraft
	if err = u.ArticleRepo.Unpublish(c, articleID); err != nil {
		if err == domain.ErrArticleNotPublished {
			return nil, domain.ErrArticleNotPublished
		}
		return nil, domain.ErrInternalServer
	}
	return article, nil
}

func (u *ArticleUsecase) ArchiveArticle(ctx context.Context, userID, articleID string) (*domain.Article, error) {
	c, close := context.WithTimeout(ctx, domain.DefaultTimeout)
	defer close()
	if !u.Policy.UserExists(userID) {
		return nil, domain.ErrUnauthorized
	}
	if !u.Policy.UserOwnsArticle(userID, domain.Article{ID: articleID}) {
		return nil, domain.ErrUnauthorizedArticleEdit
	}
	article, err := u.ArticleRepo.GetByID(c, articleID)
	if err != nil {
		if err == domain.ErrArticleNotFound {
			return nil, domain.ErrArticleNotFound
		}
		return nil, domain.ErrInternalServer
	}
	article.Status = domain.StatusArchived
	now := time.Now()
	article.Timestamps.ArchivedAt = &now
	if err = u.ArticleRepo.Archive(c, articleID, now); err != nil {
		if err == domain.ErrArticleNotPublished {
			return nil, domain.ErrArticleNotPublished
		}
		return nil, domain.ErrInternalServer
	}
	return article, nil
}

func (u *ArticleUsecase) UnarchiveArticle(ctx context.Context, userID, articleID string) (*domain.Article, error) {
	c, close := context.WithTimeout(ctx, domain.DefaultTimeout)
	defer close()
	if !u.Policy.UserExists(userID) {
		return nil, domain.ErrUnauthorized
	}
	if !u.Policy.UserOwnsArticle(userID, domain.Article{ID: articleID}) {
		return nil, domain.ErrUnauthorizedArticleEdit
	}
	article, err := u.ArticleRepo.GetByID(c, articleID)
	if err != nil {
		if err == domain.ErrArticleNotFound {
			return nil, domain.ErrArticleNotFound
		}
		return nil, domain.ErrInternalServer
	}
	article.Status = domain.StatusDraft
	if err = u.ArticleRepo.Unarchive(c, articleID); err != nil {
		if err == domain.ErrArticleNotArchived {
			return nil, domain.ErrArticleNotArchived
		}
		return nil, domain.ErrInternalServer
	}
	return article, nil
}

//===========================================================================//
//                         Article Retrieval                                 //
//===========================================================================//
func (u *ArticleUsecase) GetArticleByID(ctx context.Context, viewerID, articleID, userRole string) (*domain.Article, error) {
	// Increase view count
	// connect with the view collection

	c, close := context.WithTimeout(ctx, domain.DefaultTimeout)
	defer close()
	if !u.Policy.UserExists(viewerID) {
		return nil, domain.ErrUnauthorized
	}
	
	article, err := u.ArticleRepo.GetByID(c, articleID)
	if err != nil {
		if err == domain.ErrArticleNotFound {
			return nil, domain.ErrArticleNotFound
		}
		return nil, domain.ErrInternalServer
	}
	if u.Policy.IsAdmin(viewerID, userRole) {
		// Admins can view all articles
		return article, nil
	}
	if article.AuthorID == viewerID {
		// Authors can view their own articles
		return article, nil
	}
	if article.Status == domain.StatusPublished {
		return article, nil
	}
	return nil, domain.ErrForbidden 
}

func (u *ArticleUsecase) ViewArticleBySlug(ctx context.Context, slug, clientIP string) (*domain.Article, error) {
	// Increase view count
	// connect with the view collection
	c, close := context.WithTimeout(ctx, domain.DefaultTimeout)
	defer close()
	if slug == "" {
		return nil, domain.ErrInvalidRequest
	}
	article, err := u.ArticleRepo.GetBySlug(c, slug)
	if err != nil {
		if err == domain.ErrArticleNotFound {
			return nil, domain.ErrArticleNotFound
		}
		return nil, domain.ErrInternalServer
	}
	return article, nil
}

//===========================================================================//
//                           Article Lists                                   //
//===========================================================================//
func (u *ArticleUsecase) ListUserArticles(ctx context.Context, userID, authorID string, pag domain.Pagination) ([]domain.Article, int, error) {
	c,close := context.WithTimeout(ctx, domain.DefaultTimeout)
	defer close()
	if !u.Policy.UserExists(userID) {
		return nil, 0, domain.ErrUnauthorized
	}
	articles, length, err := u.ArticleRepo.ListByAuthor(c, authorID, pag)
	if err != nil {
		if err == domain.ErrArticleNotFound {
			return nil, 0, domain.ErrArticleNotFound
		}
		return nil, 0, domain.ErrInternalServer
	}
	return articles, length, nil
}


func (u *ArticleUsecase) ListTrendingArticles(ctx context.Context, pag domain.Pagination, windowDays int) ([]domain.Article, int, error) {
	c, close := context.WithTimeout(ctx, domain.DefaultTimeout)
	defer close()
	// Think of checking if user exists
	if windowDays <= 0 {
		windowDays = 7
	}
	articles, length, err := u.ArticleRepo.ListTrending(c, pag, windowDays)
	if err != nil {
		if err == domain.ErrArticleNotFound {
			return nil, 0, domain.ErrArticleNotFound
		}
		return nil, 0, domain.ErrInternalServer
	}
	return articles, length, nil
}

func (u *ArticleUsecase) ListArticlesByTag(ctx context.Context, userID, tag string, pag domain.Pagination) ([]domain.Article, int, error) {
	c, close := context.WithTimeout(ctx, domain.DefaultTimeout)
	defer close()
	if !u.Policy.UserExists(userID) {
		return nil, 0, domain.ErrUnauthorized
	}
	articles, length, err := u.ArticleRepo.ListByTag(c, tag, pag)
	if err != nil {
		if err == domain.ErrArticleNotFound {
			return nil, 0, domain.ErrArticleNotFound
		}
		return nil, 0, domain.ErrInternalServer
	}
	return articles, length, nil
}

func (u *ArticleUsecase) SearchArticles(ctx context.Context, userID, query string, pag domain.Pagination) ([]domain.Article, int, error) {
	c, close := context.WithTimeout(ctx, domain.DefaultTimeout)
	defer close()
	if !u.Policy.UserExists(userID) {
		return nil, 0, domain.ErrUnauthorized
	}
	articles, length, err := u.ArticleRepo.Search(c, query, pag)
	if err != nil {
		if err == domain.ErrArticleNotFound {
			return nil, 0, domain.ErrArticleNotFound
		}
		return nil, 0, domain.ErrInternalServer
	}
	return articles, length, nil
}

func (u *ArticleUsecase) FilterArticles(ctx context.Context, userID string, filter domain.ArticleFilter, pag domain.Pagination) ([]domain.Article, int, error) {
	c, close := context.WithTimeout(ctx, domain.DefaultTimeout)
	defer close()
	if !u.Policy.UserExists(userID) {
		return nil, 0, domain.ErrUnauthorized
	}
	articles, length, err := u.ArticleRepo.Filter(c, filter, pag)
	if err != nil {
		if err == domain.ErrArticleNotFound {
			return nil, 0, domain.ErrArticleNotFound
		}
		return nil, 0, domain.ErrInternalServer
	}
	return articles, length, nil
}

//===========================================================================//
//                             Engagement                                    //
//===========================================================================//
func (u *ArticleUsecase) ClapArticle(ctx context.Context, articleID, userID string) (domain.ArticleStats, error) {
	c, close := context.WithTimeout(ctx, domain.DefaultTimeout)
	defer close()
	if !u.Policy.UserExists(userID) {
		return domain.ArticleStats{}, domain.ErrUnauthorized
	}
	// Increment clap count
	stats, err := u.ArticleRepo.GetByID(c, articleID)
	if err != nil {
		if err == domain.ErrArticleNotFound {
			return domain.ArticleStats{}, domain.ErrArticleNotFound
		}
		return domain.ArticleStats{}, domain.ErrInternalServer
	}
	if stats.Stats.ClapCount >= 50 {
		return domain.ArticleStats{}, domain.ErrClapLimitExceeded
	}
	if err := u.ArticleRepo.IncrementClap(c, articleID); err != nil {
		if err == domain.ErrArticleNotFound {
			return domain.ArticleStats{}, domain.ErrArticleNotFound
		}
		return domain.ArticleStats{}, domain.ErrInternalServer
	}
	stats.Stats.ClapCount++
	return stats.Stats, nil
}
//===========================================================================//
//                 Trash Management (Author only)                            //
//===========================================================================//
func (u *ArticleUsecase) EmptyTrash(ctx context.Context, userID string) error {
	c, close := context.WithTimeout(ctx, domain.DefaultTimeout)
	defer close()
	if !u.Policy.UserExists(userID) {
		return domain.ErrUnauthorized
	}
	if err:= u.ArticleRepo.EmptyTrash(c,userID); err != nil {
		if err == domain.ErrArticleNotFound {
			return domain.ErrArticleNotFound
		}
		return domain.ErrInternalServer
	}
	return nil
}

//===========================================================================//
// 					       Admin Operations                                  //
//===========================================================================//
func (u *ArticleUsecase) AdminListAllArticles(ctx context.Context, userID, userRole string, pag domain.Pagination) ([]domain.Article, int, error) {
	c, close := context.WithTimeout(ctx, domain.DefaultTimeout)
	defer close()
	if !u.Policy.UserExists(userID) {
		return nil, 0, domain.ErrUnauthorized
	}
	if !u.Policy.IsAdmin(userID, userRole) {
		return nil, 0, domain.ErrUnauthorized
	}
	articles, length, err := u.ArticleRepo.AdminListAllArticles(c,pag)
	if err != nil {
		if err == domain.ErrArticleNotFound {
			return nil, 0, domain.ErrArticleNotFound
		}
		return nil, 0, domain.ErrInternalServer
	}
	return articles, length, nil
}

func (u *ArticleUsecase) AdminHardDeleteArticle(ctx context.Context, userID, userRole, articleID string) error {
	c, close := context.WithTimeout(ctx, domain.DefaultTimeout)
	defer close()
	if !u.Policy.IsAdmin(userID, userRole) {
		return domain.ErrUnauthorized
	}
	_, err := u.ArticleRepo.GetByID(c, articleID)
	if err != nil {
		if err == domain.ErrArticleNotFound {
			return domain.ErrArticleNotFound
		}
		return domain.ErrInternalServer
	}
	if err:= u.ArticleRepo.HardDelete(c, articleID); err != nil {
		return domain.ErrInternalServer
	}
	return nil
}

func (u *ArticleUsecase) AdminUnpublishArticle(ctx context.Context, userID, userRole, articleID string) (*domain.Article, error) {
	c, close := context.WithTimeout(ctx, domain.DefaultTimeout)
	defer close()
	if !u.Policy.IsAdmin(userID, userRole) {
		return nil, domain.ErrUnauthorized
	}
	article, err := u.ArticleRepo.GetByID(c, articleID)
	if err != nil {
		if err == domain.ErrArticleNotFound {
			return nil, domain.ErrArticleNotFound
		}
		return nil, domain.ErrInternalServer
	}
	if err := u.ArticleRepo.Unpublish(c, articleID); err != nil {
		if err == domain.ErrArticleNotPublished {
			return nil, domain.ErrArticleNotPublished
		}
		return nil, domain.ErrInternalServer
	}
	return article, nil
}