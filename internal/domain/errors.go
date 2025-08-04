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
	ErrInternalServer = Error{Code: "GEN_001", Message: "Internal server error"}
	ErrInvalidRequest = Error{Code: "GEN_002", Message: "Invalid request payload"}
	// ErrUnauthorized        = Error{Code: "GEN_003", Message: "Unauthorized access"}
	ErrForbidden           = Error{Code: "GEN_004", Message: "Forbidden"}
	ErrNotFound            = Error{Code: "GEN_005", Message: "Resource not found"}
	ErrRateLimitExceeded   = Error{Code: "GEN_006", Message: "Too many requests"}
	ErrMethodNotAllowed    = Error{Code: "GEN_007", Message: "HTTP method not allowed"}
	ErrConflict            = Error{Code: "GEN_008", Message: "Resource conflict"}
	ErrUnprocessableEntity = Error{Code: "GEN_009", Message: "Unprocessable entity"}
)

// ===========================================================================//
//
//	User Errors                                     //
//
// ===========================================================================//
var (
	ErrUserNotFound       = Error{Code: "USER_001", Message: "User not found"}
	ErrEmailAlreadyTaken  = Error{Code: "USER_002", Message: "Email already registered"}
	ErrInvalidCredentials = Error{Code: "USER_003", Message: "Invalid email or password"}
	ErrUnauthorized       = Error{Code: "USER_004", Message: "Unauthorized access"}
	ErrInvalidUserID      = Error{Code: "USER_005", Message: "Invalid user ID format"}
	ErrUserBanned         = Error{Code: "USER_006", Message: "User account is banned"}
	ErrProfileIncomplete  = Error{Code: "USER_007", Message: "User profile is incomplete"}

	ErrWeakPassword = Error{Code: "USER_008", Message: "Password must contain at least one uppercase letter, one lowercase letter, one number, and one special character"}

	ErrEmailAlreadyExists          = Error{Code: "USER_009", Message: "Email already exists"}
	ErrUsernameAlreadyExists       = Error{Code: "USER_010", Message: "Username already exists"}
	ErrUserNotVerified             = Error{Code: "USER_013", Message: "User email is not verified"}
	ErrUserDeactivated             = Error{Code: "USER_014", Message: "User account is deactivated"}
	ErrInvalidToken                = Error{Code: "USER_011", Message: "Invalid or malformed token"}
	ErrExpiredToken                = Error{Code: "USER_012", Message: "Token has expired"}
	ErrJWTExpired                  = Error{Code: "TOKEN_002", Message: "JWT has expired"}
	ErrUnexpectedSigningMethod     = Error{Code: "TOKEN_001", Message: "Unexpected signing method"}
	ErrAuthorizationHeaderRequired = Error{Code: "AUTH_001", Message: "Authorization header is required"}
)

// ===========================================================================//
//
//	Article Errors                                     //
//
// ===========================================================================//
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
