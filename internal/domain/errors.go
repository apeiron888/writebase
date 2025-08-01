package domain

type Error struct {
	Code    string
	Message string
}

func (e Error) Error() string {
	return e.Message
}

//===========================================================================//
//                        General Errors                                     //
//===========================================================================//
var (
	ErrInternalServer      = Error{Code: "GEN_001", Message: "Internal server error"}
	ErrInvalidRequest      = Error{Code: "GEN_002", Message: "Invalid request payload"}
	ErrUnauthorized        = Error{Code: "GEN_003", Message: "Unauthorized access"}
	ErrForbidden           = Error{Code: "GEN_004", Message: "Forbidden"}
	ErrNotFound            = Error{Code: "GEN_005", Message: "Resource not found"}
	ErrRateLimitExceeded   = Error{Code: "GEN_006", Message: "Too many requests"}
	ErrMethodNotAllowed    = Error{Code: "GEN_007", Message: "HTTP method not allowed"}
	ErrConflict            = Error{Code: "GEN_008", Message: "Resource conflict"}
	ErrUnprocessableEntity = Error{Code: "GEN_009", Message: "Unprocessable entity"}
)
//===========================================================================//
//                        Article Errors                                     //
//===========================================================================//
var (
	ErrArticleNotFound           = Error{Code: "ARTICLE_001", Message: "Article not found"}
	ErrInvalidArticleID          = Error{Code: "ARTICLE_002", Message: "Invalid article ID format"}
	ErrArticleTitleEmpty         = Error{Code: "ARTICLE_003", Message: "Article title cannot be empty"}
	ErrArticleContentEmpty       = Error{Code: "ARTICLE_004", Message: "Article content cannot be empty"}
	ErrArticleTooLong            = Error{Code: "ARTICLE_005", Message: "Article exceeds maximum length"}
	ErrArticleTagLimitExceeded   = Error{Code: "ARTICLE_006", Message: "Maximum number of tags exceeded"}
	ErrDuplicateArticleSlug      = Error{Code: "ARTICLE_007", Message: "Article slug already exists"}
	ErrUnauthorizedArticleEdit   = Error{Code: "ARTICLE_008", Message: "You are not allowed to edit this article"}
	ErrUnauthorizedArticleDelete = Error{Code: "ARTICLE_009", Message: "You are not allowed to delete this article"}
	ErrArticleAlreadyPublished   = Error{Code: "ARTICLE_010", Message: "Article is already published"}
	ErrMaxArticlesPerUser        = Error{Code: "ARTICLE_011", Message: "User has reached the article creation limit"}
	ErrArticleSearchFailed       = Error{Code: "ARTICLE_012", Message: "Failed to search articles"}

	
	ErrMaxClaps = Error{Code: "CLAP_001", Message: "Maximum claps reached"}
)