package domain

type Error struct {
	Code    string
	Message string
}

func (e Error) Error() string {
	return e.Message
}

// ===========================================================================//
//
//	General Errors                                     //
//
// ===========================================================================//
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

// ===========================================================================//
//
//	Article Errors                                     //
//
// ===========================================================================//
var (
	ErrArticleNotFound           = Error{Code: "ARTICLE_001", Message: "Article not found"}
	ErrArticleAlreadyExists      = Error{Code: "ARTICLE_002", Message: "Article already exists"}
	ErrInvalidArticleID          = Error{Code: "ARTICLE_003", Message: "Invalid article ID format"}
	ErrArticleTitleEmpty         = Error{Code: "ARTICLE_004", Message: "Article title cannot be empty"}
	ErrArticleContentEmpty       = Error{Code: "ARTICLE_005", Message: "Article content cannot be empty"}
	ErrArticleTooLong            = Error{Code: "ARTICLE_006", Message: "Article exceeds maximum length"}
	ErrArticleTagLimitExceeded   = Error{Code: "ARTICLE_007", Message: "Maximum number of tags exceeded"}
	ErrDuplicateArticleSlug      = Error{Code: "ARTICLE_008", Message: "Article slug already exists"}
	ErrUnauthorizedArticleEdit   = Error{Code: "ARTICLE_009", Message: "You are not allowed to edit this article"}
	ErrUnauthorizedArticleDelete = Error{Code: "ARTICLE_010", Message: "You are not allowed to delete this article"}
	ErrArticleAlreadyPublished   = Error{Code: "ARTICLE_011", Message: "Article is already published"}
	ErrMaxArticlesPerUser        = Error{Code: "ARTICLE_012", Message: "User has reached the article creation limit"}
	ErrArticleSearchFailed       = Error{Code: "ARTICLE_013", Message: "Failed to search articles"}
	ErrArticleNotPublished       = Error{Code: "ARTICLE_014", Message: "Article is not published"}
	ErrArticleNotArchived        = Error{Code: "ARTICLE_015", Message: "Article is not archived"}
	ErrNoChangesDetected         = Error{Code: "ARTICLE_016", Message: "No changes detected in the article"}
	ErrClapLimitExceeded         = Error{Code: "ARTICLE_017", Message: "Clap limit exceeded"}
)

// Comment-related errors
var (
	ErrCommentNotFound     = Error{Code: "COMMENT_001", Message: "Comment not found"}
	ErrInvalidCommentID    = Error{Code: "COMMENT_002", Message: "Invalid comment ID format"}
	ErrEmptyCommentContent = Error{Code: "COMMENT_003", Message: "Comment content cannot be empty"}
	ErrCommentPermission   = Error{Code: "COMMENT_004", Message: "You do not have permission to modify this comment"}
)

// Reaction-related errors
var (
	ErrReactionNotFound    = Error{Code: "REACTION_001", Message: "Reaction not found"}
	ErrInvalidReactionType = Error{Code: "REACTION_002", Message: "Invalid reaction type"}
	ErrAlreadyReacted      = Error{Code: "REACTION_003", Message: "Already reacted"}
)

// Follow-related errors
var (
	ErrFollowNotFound   = Error{Code: "FOLLOW_001", Message: "Follow relationship not found"}
	ErrAlreadyFollowing = Error{Code: "FOLLOW_002", Message: "Already following this user"}
	ErrCannotFollowSelf = Error{Code: "FOLLOW_003", Message: "Cannot follow yourself"}
)

// Report-related errors
var (
	ErrReportNotFound        = Error{Code: "REPORT_001", Message: "Report not found"}
	ErrInvalidReportTarget   = Error{Code: "REPORT_002", Message: "Invalid report target"}
	ErrReportAlreadyResolved = Error{Code: "REPORT_003", Message: "Report already resolved"}
)
