package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"
	"write_base/internal/domain"
	"write_base/internal/infrastructure/ai"
)
type ArticleUsecase struct {
    Repo        domain.IArticleRepository
    Policy      domain.IPolicy 
    Utils       domain.IUtils
    AIClient    domain.IAI  // <-- Use interface instead of *ai.GeminiClient
    TagUsecase  domain.TagUsecase
    ViewUsecase domain.ViewUsecase
    ClapUsecase domain.ClapUsecase

}

func NewArticleUsecase(repo domain.IArticleRepository, policy domain.IPolicy, util domain.IUtils,tagusecase domain.TagUsecase, vuc domain.ViewUsecase, clap domain.ClapUsecase, aiClient *ai.GeminiClient) domain.IArticleUsecase{
	return &ArticleUsecase{Repo: repo, Policy: policy, Utils: util, TagUsecase: tagusecase, ViewUsecase: vuc, ClapUsecase: clap, AIClient: aiClient,}
}
//===============================================================================//
//                                CRUD                                           //
//===============================================================================//
// =============================== Article Create ================================
func (au *ArticleUsecase) CreateArticle(ctx context.Context, userID string, input *domain.Article) (string, error) {
    c, cancel := context.WithTimeout(ctx, domain.DefaultTimeout)
    defer cancel()

    // if !au.Policy.UserExists(userID) {
    //     return "", domain.ErrUnauthorized
    // }
    if !au.Policy.ArticleCreateValid(input) {
        return "", domain.ErrInvalidArticlePayload
    }
    // Initialize article fields
    input.ID = au.Utils.GenerateUUID()
    input.AuthorID = userID
    now := time.Now()
    input.Timestamps = domain.ArticleTimes{
        CreatedAt: now,
        UpdatedAt: now,
    }
    input.Status = domain.StatusDraft
    input.Stats = domain.ArticleStats{
        ViewCount: 0,
        ClapCount: 0,
    }
    if input.Slug == "" {
        input.Slug = au.Utils.GenerateSlug(input.Title)
        if _,err:=au.Repo.GetBySlug(c,input.Slug);err!=nil{
            suffix := au.Utils.GenerateShortUUID()
            input.Slug += suffix
        }
    }
    if input.Excerpt == "" {
        input.Excerpt = input.Title
        if len(input.Excerpt) > domain.MaxExcerptLength {
            input.Excerpt = input.Excerpt[:domain.MaxExcerptLength]
        }
    }
	if err := au.TagUsecase.ValidateTags(input.Tags); err != nil {
		return "", domain.ErrInvalidTagName
	}
    if err := au.Repo.Create(c, input); err != nil {
        return "", fmt.Errorf("repository error: %w", err)
    }
    return input.ID, nil
}
// =============================== Article Update ================================
func (au *ArticleUsecase) UpdateArticle(ctx context.Context, userID string, input *domain.Article) error {
    c, cancel := context.WithTimeout(ctx, domain.DefaultTimeout)
    defer cancel()

    // if !au.Policy.UserExists(userID) {
    //     return domain.ErrUnauthorized
    // }

    if !au.Policy.UserOwnsArticle(userID,input){
        return domain.ErrUnauthorized
    }
    if !au.Policy.ArticleCreateValid(input) {
        return domain.ErrInvalidArticlePayload
    }
    if input.Slug=="" || strings.Contains(input.Slug," ") {
        input.Slug = au.Utils.GenerateSlug(input.Title)
        if _,err:=au.Repo.GetBySlug(c,input.Slug);err!=nil{
            suffix := au.Utils.GenerateShortUUID()
            input.Slug += suffix
        }
    }
    if err := au.TagUsecase.ValidateTags(input.Tags); err != nil {
		return domain.ErrInvalidTagName
	}
    old,err:= au.GetArticleByID(c, input.ID, userID)
    if err!=nil || old.ID!=input.ID {
        return domain.ErrArticleNotFound
    }

    if err:=au.Repo.Update(c,input); err!=nil{
        return domain.ErrInternalServer
    }
    return nil
}
// =============================== Article Delete ================================
func (au *ArticleUsecase) DeleteArticle(ctx context.Context, articleID,userID string) error {
	c, close := context.WithTimeout(ctx, domain.DefaultTimeout)
	defer close()

    // if !au.Policy.UserExists(userID) {
    //     return domain.ErrUnauthorized
    // }

    res,err:= au.GetArticleByID(c,articleID, userID)
    if err!=nil {
        return domain.ErrArticleNotFound
    }

    if !au.Policy.UserOwnsArticle(userID,res){
        return domain.ErrUnauthorized
    }

	if err = au.Repo.Delete(c, articleID); err != nil {
		if err == domain.ErrArticleNotFound {
			return domain.ErrArticleNotFound
		}
		return domain.ErrInternalServer
	}
	return nil
}
// =============================== Article Restore ================================
func (au *ArticleUsecase) RestoreArticle(ctx context.Context,  articleID, userID string)  error {
	c, close := context.WithTimeout(ctx, domain.DefaultTimeout)
	defer close()

    // if !au.Policy.UserExists(userID) {
    //     return domain.ErrUnauthorized
    // }

    res,err:= au.GetArticleByID(c,articleID, userID)
    if err!=nil {
        return domain.ErrArticleNotFound
    }

    if !au.Policy.UserOwnsArticle(userID,res){
        return domain.ErrUnauthorized
    }

	if err = au.Repo.Restore(c, articleID); err != nil {
		if err == domain.ErrArticleNotFound {
			return domain.ErrArticleNotFound
		}
		return domain.ErrInternalServer
	}
	return nil
}
//===============================================================================//
//                   Article State Management                                    //
//===============================================================================//
// ======================== Article Publish =======================================
func (au *ArticleUsecase) PublishArticle(ctx context.Context, articleID, userID string) (*domain.Article, error) {
    c, cancel := context.WithTimeout(ctx, domain.DefaultTimeout)
    defer cancel()

    // if !au.Policy.UserExists(userID) {
    //     return domain.ErrUnauthorized
    // }
    article,err := au.GetArticleByID(c, articleID,userID)
    if err!=nil{
        if err == domain.ErrArticleNotFound{
            return nil,err
        }
        return nil,domain.ErrInternalServer
    }
    if !au.Policy.UserOwnsArticle(userID,article){
        return nil,domain.ErrUnauthorized
    }
    if article.Status==domain.ArticleStatus(domain.StatusPublished){
        return nil, domain.ErrArticlePublished
    }
	// Check all tags are approved
	for _, tag := range article.Tags {
		if !au.TagUsecase.IsTagApproved(tag) {
			return nil, domain.ErrUnapprovedTags
		}
	}
	article.Status = domain.StatusPublished
	now := time.Now()
	article.Timestamps.PublishedAt = &now
    if err := au.Repo.Publish(c,articleID,now);err!=nil {
        return nil,domain.ErrInternalServer
    }
	return article, nil
}
// ======================== Article Unpublish =====================================
func (au *ArticleUsecase) UnpublishArticle(ctx context.Context, articleID, userID string) (*domain.Article, error) {
    c, cancel := context.WithTimeout(ctx, domain.DefaultTimeout)
    defer cancel()

    // if !au.Policy.UserExists(userID) {
    //     return domain.ErrUnauthorized
    // }

    article,err := au.GetArticleByID(c, articleID,userID)
    if err!=nil{
        if err == domain.ErrArticleNotFound{
            return nil,err
        }
        return nil,domain.ErrInternalServer
    }
    if !au.Policy.UserOwnsArticle(userID,article){
        return nil,domain.ErrUnauthorized
    }
    if article.Status!=domain.ArticleStatus(domain.StatusPublished){
        return nil, domain.ErrArticleNotPublished
    }
	article.Status = domain.StatusDraft
    if err := au.Repo.Unpublish(c,articleID);err!=nil {
        return nil,domain.ErrInternalServer
    }
	return article, nil
}
// ======================== Article Archive =======================================
func (au *ArticleUsecase) ArchiveArticle(ctx context.Context, articleID, userID string) (*domain.Article, error) {
    c, cancel := context.WithTimeout(ctx, domain.DefaultTimeout)
    defer cancel()

    // if !au.Policy.UserExists(userID) {
    //     return domain.ErrUnauthorized
    // }
    article,err := au.GetArticleByID(c, articleID,userID)
    if err!=nil{
        if err == domain.ErrArticleNotFound{
            return nil,err
        }
        return nil,domain.ErrInternalServer
    }
    if !au.Policy.UserOwnsArticle(userID,article){
        return nil,domain.ErrUnauthorized
    }
    if article.Status==domain.StatusArchived {
        return nil, domain.ErrArticleArchived
    }
	article.Status = domain.StatusArchived
	now := time.Now()
	article.Timestamps.ArchivedAt = &now
    if err := au.Repo.Archive(c,articleID,now);err!=nil {
        return nil,domain.ErrInternalServer
    }
	return article, nil
}
// ======================== Article Unarchive =====================================
func (au *ArticleUsecase) UnarchiveArticle(ctx context.Context, articleID, userID string) (*domain.Article, error) {
    c, cancel := context.WithTimeout(ctx, domain.DefaultTimeout)
    defer cancel()

    // if !au.Policy.UserExists(userID) {
    //     return domain.ErrUnauthorized
    // }
    article,err := au.GetArticleByID(c, articleID,userID)
    if err!=nil{
        if err == domain.ErrArticleNotFound{
            return nil,err
        }
        return nil,domain.ErrInternalServer
    }
    if !au.Policy.UserOwnsArticle(userID,article){
        return nil,domain.ErrUnauthorized
    }
    if article.Status!=domain.ArticleStatus(domain.StatusArchived){
        return nil, domain.ErrArticleNotArchived
    }
	article.Status = domain.StatusDraft
    if err := au.Repo.Unarchive(c,articleID);err!=nil {
        return nil,domain.ErrInternalServer
    }
	return article, nil
}
//===============================================================================//
//                               Retrieve                                        //
//===============================================================================//
// =============================== Article GetByID ================================
func (au *ArticleUsecase) GetArticleByID(ctx context.Context, articleID, userID string) (*domain.Article, error) {
    c, cancel := context.WithTimeout(ctx, domain.DefaultTimeout)
    defer cancel()

    // if !au.Policy.UserExists(userID) {
    //     return domain.ErrUnauthorized
    // }
    if articleID==""{
        return nil,domain.ErrInvalidArticlePayload
    }
    article,err := au.Repo.GetByID(c, articleID)
    if err!=nil{
        if err == domain.ErrArticleNotFound{
            return nil,err
        }
        return nil,domain.ErrInternalServer
    }

    if article.AuthorID == userID {
		return article, nil
	}
	if article.Status == domain.StatusPublished {
		// Record view for authenticated users
		if userID != "" {
			au.ViewUsecase.RecordView(c, userID, articleID, "")
		}
		
		// Increment article view count
		au.Repo.IncrementView(c, article.ID)
		return article, nil
	}
    return nil,domain.ErrUnauthorized
}

// =============================== Article Stats ================================
func (au *ArticleUsecase) GetArticleStats(ctx context.Context, articleID, userID string) (*domain.ArticleStats, error) {
    c, cancel := context.WithTimeout(ctx, domain.DefaultTimeout)
    defer cancel()

    // if !au.Policy.UserExists(userID) {
    //     return domain.ErrUnauthorized
    // }
    if articleID==""{
        return nil,domain.ErrInvalidArticlePayload
    }
    articleStats,err := au.Repo.GetStats(c, articleID)
    if err!=nil{
        if err == domain.ErrArticleNotFound{
            return nil,err
        }
        return nil,domain.ErrInternalServer
    }
    return articleStats, nil
}
// =================== All Article Stats of Author ==============================
func (au *ArticleUsecase) GetAllArticleStats(ctx context.Context, userID string) ([]domain.ArticleStats,int, error) {
    c, cancel := context.WithTimeout(ctx, domain.DefaultTimeout)
    defer cancel()

    // if !au.Policy.UserExists(userID) {
    //     return domain.ErrUnauthorized
    // }
    articlesStats,total,err := au.Repo.GetAllArticleStats(c,userID)
    if err!=nil {
        if err == domain.ErrArticleNotFound{
            return nil,0,err
        }
        return nil,0,domain.ErrInternalServer
    }
    return articlesStats, total, nil
}
// =================== Article GetBySlug ==============================
func (au *ArticleUsecase) GetArticleBySlug(ctx context.Context, slug string, clientIP string) (*domain.Article, error) {
    c, cancel := context.WithTimeout(ctx, domain.DefaultTimeout)
    defer cancel()

    article,err := au.Repo.GetBySlug(c, slug)
    if err!=nil{
        if err == domain.ErrArticleNotFound{
            return nil,err
        }
        return nil,domain.ErrInternalServer
    }
	// Record view with client IP
	if clientIP != "" {
		au.ViewUsecase.RecordView(c, "", article.ID, clientIP)
	}
	
	// Increment article view count
	au.Repo.IncrementView(c, article.ID)
    return article,nil
}

//===========================================================================//
//                           Article Lists                                   //
//===========================================================================//
//======================= List user articles ==================================
func (au *ArticleUsecase) ListArticlesByAuthor(ctx context.Context, userID, authorID string, pag domain.Pagination) ([]domain.Article, int, error) {
	c,close := context.WithTimeout(ctx, domain.DefaultTimeout)
	defer close()

	if !au.Policy.UserExists(userID) {
		return nil, 0, domain.ErrUnauthorized
	}
	articles, length, err := au.Repo.ListByAuthor(c, authorID, pag)
	if err != nil {
		if err == domain.ErrArticleNotFound {
			return nil, 0, domain.ErrArticleNotFound
		}
		return nil, 0, domain.ErrInternalServer
	}
	return articles, length, nil
}
//======================= List Trending articles ==================================
func (au *ArticleUsecase) GetTrendingArticles(ctx context.Context, userID string, pag domain.Pagination) ([]domain.Article, int, error) {
    c, cancel := context.WithTimeout(ctx, domain.DefaultTimeout)
    defer cancel()

    if !au.Policy.UserExists(userID) {
        return nil, 0, domain.ErrUnauthorized
    }

    // Pass time window for trending calculation => 7 days
    trendingWindow := 7
    articles, total, err := au.Repo.FindTrending(c, trendingWindow, pag)
    if err != nil {
        if err == domain.ErrArticleNotFound{
            return nil,0,domain.ErrArticleNotFound
        }
        return nil, 0, domain.ErrInternalServer
    }
    return articles, total, nil
}

//======================= List New articles =======================================


func (au *ArticleUsecase) GetNewArticles(ctx context.Context, userID string, pag domain.Pagination) ([]domain.Article, int, error) {
    c, cancel := context.WithTimeout(ctx, domain.DefaultTimeout)
    defer cancel()

    if !au.Policy.UserExists(userID) {
        return nil, 0, domain.ErrUnauthorized
    }

    articles, total, err := au.Repo.FindNewArticles(c,pag)
    if err != nil {
        if err == domain.ErrArticleNotFound{
            return nil,0,domain.ErrArticleNotFound
        }
        return nil, 0, domain.ErrInternalServer
    }
    return articles, total, nil
}
//======================= List Popular articles ===================================
func (au *ArticleUsecase) GetPopularArticles(ctx context.Context, userID string, pag domain.Pagination) ([]domain.Article, int, error) {
    c, cancel := context.WithTimeout(ctx, domain.DefaultTimeout)
    defer cancel()

    if !au.Policy.UserExists(userID) {
        return nil, 0, domain.ErrUnauthorized
    }

    articles, total, err := au.Repo.FindPopularArticles(c, pag)
    if err != nil {
        if err == domain.ErrArticleNotFound{
            return nil,0,domain.ErrArticleNotFound
        }
        return nil, 0, domain.ErrInternalServer
    }
    return articles, total, nil
}

//======================= List Author articles ===============================
func (au *ArticleUsecase) FilterAuthorArticles(ctx context.Context, callerID, authorID string, filter domain.ArticleFilter, pag domain.Pagination) ([]domain.Article, int, error) {
    c, cancel := context.WithTimeout(ctx, domain.DefaultTimeout)
    defer cancel()

    // Caller must exist (auth middleware normally ensures this, but be defensive)
    if !au.Policy.UserExists(callerID) {
        return nil, 0, domain.ErrUnauthorized
    }

    // Ensure the target author exists
    if !au.Policy.UserExists(authorID) {
        return nil, 0, domain.ErrAuthorNotFound
    }

    // Authorization: only the author themselves (or admin if you add it)
    if callerID != authorID {
        return nil, 0, domain.ErrUnauthorized
    }

    // Defensive pagination defaults (repo also defends)
    if pag.Page < 1 {
        pag.Page = 1
    }
    if pag.PageSize <= 0 {
        pag.PageSize = 20
    }

    articles, total, err := au.Repo.FilterAuthorArticles(c, authorID, filter, pag)
    if err != nil {
        if err == domain.ErrArticleNotFound {
            return nil, 0, domain.ErrArticleNotFound
        }
        return nil, 0, domain.ErrInternalServer
    }
    return articles, total, nil
}
//======================= List Author articles (all users) =============================
func (au *ArticleUsecase) FilterArticles(ctx context.Context, filter domain.ArticleFilter, pag domain.Pagination) ([]domain.Article, int, error) {
    c, cancel := context.WithTimeout(ctx, domain.DefaultTimeout)
    defer cancel()

    // If statuses not provided, default to published (public view)
    if len(filter.Statuses) == 0 {
        filter.Statuses = []domain.ArticleStatus{domain.StatusPublished}
    }

    // Defensive pagination defaults
    if pag.Page < 1 {
        pag.Page = 1
    }
    if pag.PageSize <= 0 {
        pag.PageSize = 20
    }

    articles, total, err := au.Repo.Filter(c, filter, pag)
    if err != nil {
        if err == domain.ErrArticleNotFound {
            return nil, 0, domain.ErrArticleNotFound
        }
        return nil, 0, domain.ErrInternalServer
    }
    return articles, total, nil
}

//=========================== Search ==============================================
func (u *ArticleUsecase) SearchArticles(ctx context.Context, userID, query string, pag domain.Pagination) ([]domain.Article, int, error) {
	c, close := context.WithTimeout(ctx, domain.DefaultTimeout)
	defer close()
	if !u.Policy.UserExists(userID) {
		return nil, 0, domain.ErrUnauthorized
	}
	articles, length, err := u.Repo.Search(c, query, pag)
	if err != nil {
		if err == domain.ErrArticleNotFound {
			return nil, 0, domain.ErrArticleNotFound
		}
		return nil, 0, domain.ErrInternalServer
	}
	return articles, length, nil
}
//======================== List By Tags =======================================
func (u *ArticleUsecase) ListArticlesByTags(ctx context.Context, userID string, tags []string, pag domain.Pagination) ([]domain.Article, int, error) {
	c, close := context.WithTimeout(ctx, domain.DefaultTimeout)
	defer close()
	
	if !u.Policy.UserExists(userID) {
		return nil, 0, domain.ErrUnauthorized
	}
	
	// Validate tags
	for _, tag := range tags {
		if !u.TagUsecase.IsTagApproved(tag) {
			return nil, 0, domain.ErrUnapprovedTags
		}
	}
	
	return u.Repo.ListByTags(c, tags, pag)
}
//===========================================================================//
//                 Trash Management (Author only)                            //
//===========================================================================//
func (au *ArticleUsecase) EmptyTrash(ctx context.Context, userID string) error {
	c, close := context.WithTimeout(ctx, domain.DefaultTimeout)
	defer close()
	if !au.Policy.UserExists(userID) {
		return domain.ErrUnauthorized
	}
	if err:= au.Repo.EmptyTrash(c,userID); err != nil {
		if err == domain.ErrArticleNotFound {
			return domain.ErrArticleNotFound
		}
		return domain.ErrInternalServer
	}
	return nil
}
//========================== Delete Article From Trash ==========================
func (au *ArticleUsecase) DeleteArticleFromTrash(ctx context.Context,articleID, userID string) error {
	c, close := context.WithTimeout(ctx, domain.DefaultTimeout)
	defer close()
	if !au.Policy.UserExists(userID) {
		return domain.ErrUnauthorized
	}
    input,err := au.GetArticleByID(c,articleID,userID)
    if err!=nil {
        return domain.ErrArticleNotFound
    }
    if input.Status!=domain.StatusDeleted {
        return domain.ErrArticleNotFound
    }
    if !au.Policy.UserOwnsArticle(userID,input){
        return domain.ErrUnauthorized
    }

	if err:= au.Repo.DeleteFromTrash(c,articleID,userID); err != nil {
		if err == domain.ErrArticleNotFound {
			return domain.ErrArticleNotFound
		}
		return domain.ErrInternalServer
	}
	return nil
}

// ===========================================================================//
//                            Admin Operations                                //
// ===========================================================================//
func (u *ArticleUsecase) AdminListAllArticles(ctx context.Context, userID, userRole string, pag domain.Pagination) ([]domain.Article, int, error) {
	c, close := context.WithTimeout(ctx, domain.DefaultTimeout)
	defer close()
	if !u.Policy.UserExists(userID) {
		return nil, 0, domain.ErrUnauthorized
	}
	if !u.Policy.IsAdmin(userID, userRole) {
		return nil, 0, domain.ErrUnauthorized
	}
	articles, length, err := u.Repo.AdminListAllArticles(c,pag)
	if err != nil {
		if err == domain.ErrArticleNotFound {
			return nil, 0, domain.ErrArticleNotFound
		}
		return nil, 0, domain.ErrInternalServer
	}
	return articles, length, nil
}
//================================ Hard Delete (Admin) ====================================
func (u *ArticleUsecase) AdminHardDeleteArticle(ctx context.Context, userID, userRole, articleID string) error {
	c, close := context.WithTimeout(ctx, domain.DefaultTimeout)
	defer close()
	if !u.Policy.IsAdmin(userID, userRole) {
		return domain.ErrUnauthorized
	}
	_, err := u.Repo.GetByID(c, articleID)
	if err != nil {
		if err == domain.ErrArticleNotFound {
			return domain.ErrArticleNotFound
		}
		return domain.ErrInternalServer
	}
	if err:= u.Repo.HardDelete(c, articleID); err != nil {
		return domain.ErrInternalServer
	}
	return nil
}
//============================ Unpublish Article (Admin) ================================

func (au *ArticleUsecase) AdminUnpublishArticle(ctx context.Context, userID, userRole, articleID string) (*domain.Article, error) {
	c, close := context.WithTimeout(ctx, domain.DefaultTimeout)
	defer close()
	if !au.Policy.IsAdmin(userID, userRole) {
		return nil, domain.ErrUnauthorized
	}
	article, err := au.Repo.GetByID(c, articleID)
	if err != nil {
		if err == domain.ErrArticleNotFound {
			return nil, domain.ErrArticleNotFound
		}
		return nil, domain.ErrInternalServer
	}
	if err := au.Repo.Unpublish(c, articleID); err != nil {
		if err == domain.ErrArticleNotPublished {
			return nil, domain.ErrArticleNotPublished
		}
		return nil, domain.ErrInternalServer
	}
	return article, nil
}


//================================== CLAPPING ===========================================
// Add new method
func (u *ArticleUsecase) AddClap(ctx context.Context, userID, articleID string) (domain.ArticleStats, error) {
	c, cancel := context.WithTimeout(ctx, domain.DefaultTimeout)
	defer cancel()
	
	// Check if user exists
	if !u.Policy.UserExists(userID) {
		return domain.ArticleStats{}, domain.ErrUnauthorized
	}
	
	// Add clap
	totalClaps, err := u.ClapUsecase.AddClap(c, userID, articleID)
	if err != nil {
		return domain.ArticleStats{}, err
	}
	
	// Update article clap count
	if err := u.Repo.UpdateClapCount(c, articleID, totalClaps); err != nil {
		return domain.ArticleStats{}, err
	}
	
	// Get updated article stats
	article, err := u.Repo.GetByID(c, articleID)
	if err != nil {
		return domain.ArticleStats{}, err
	}
	
	return article.Stats, nil
}