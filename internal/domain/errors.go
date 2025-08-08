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

	ErrPasswordHashingFailed        = Error{Code: "USER_017", Message: "Failed to hash password"}
	ErrUserCreationFailed           = Error{Code: "USER_018", Message: "Failed to create user"}
	ErrUserUpdateFailed             = Error{Code: "USER_021", Message: "Failed to update user"}
	ErrAccessTokenGenerationFailed  = Error{Code: "USER_022", Message: "Failed to generate access token"}
	ErrRefreshTokenGenerationFailed = Error{Code: "USER_023", Message: "Failed to generate refresh token"}
	ErrVerificationTokenSaveFailed  = Error{Code: "USER_019", Message: "Failed to save verification token"}
	ErrSendVerificationEmailFailed  = Error{Code: "USER_020", Message: "Failed to send verification email"}

	ErrEmailAlreadyExists          = Error{Code: "USER_009", Message: "Email already exists"}
	ErrUsernameAlreadyExists       = Error{Code: "USER_010", Message: "Username already exists"}
	ErrEmailNotRegistered          = Error{Code: "USER_016", Message: "This email is not registered"}
	ErrMissingVerifyCode           = Error{Code: "USER_015", Message: "Missing verification code"}
	ErrUserNotVerified             = Error{Code: "USER_013", Message: "User email is not verified"}
	ErrUserDeactivated             = Error{Code: "USER_014", Message: "User account is deactivated"}
	ErrInvalidToken                = Error{Code: "USER_011", Message: "Invalid or malformed token"}
	ErrExpiredToken                = Error{Code: "USER_012", Message: "Token has expired"}
    ErrUserIDNotFound = Error{Code: "USER_024", Message: "User ID not found in context"}
	ErrJWTExpired                  = Error{Code: "TOKEN_002", Message: "JWT has expired"}
	ErrUnexpectedSigningMethod     = Error{Code: "TOKEN_001", Message: "Unexpected signing method"}
	ErrAuthorizationHeaderRequired = Error{Code: "AUTH_001", Message: "Authorization header is required"}
	ErrMissingOrExpiredStateCookie = Error{Code: "AUTH_002", Message: "Missing or expired state cookie"}
	ErrMissingState                = Error{Code: "AUTH_003", Message: "Missing state parameter in OAuth callback"}
	ErrMissingOAuthStateToken      = Error{Code: "AUTH_004", Message: "Missing oauthStateToken parameter in OAuth callback"}
	ErrOAuthLoginFailed            = Error{Code: "AUTH_005", Message: "OAuth login failed"}
	ErrFailedToFetchUserInfo       = Error{Code: "AUTH_006", Message: "Failed to fetch user info from OAuth provider"}
	ErrTokenExchangeFailed         = Error{Code: "AUTH_007", Message: "Token exchange failed during OAuth process"}
	ErrMissingOAuthCode            = Error{Code: "AUTH_008", Message: "Missing code parameter in OAuth callback"}
	ErrRefreshTokenExpired         = Error{Code: "AUTH_009", Message: "Refresh token has expired"}
	ErrRefreshTokenRevoked         = Error{Code: "AUTH_010", Message: "Refresh token has been revoked"}
	ErrRefreshTokenNotFound        = Error{Code: "AUTH_011", Message: "Refresh token not found"}
	ErrSuperAdminCannotBeDemoted = Error{Code:"AUTH_012", Message:"Super admin cannot be demoted"}
	ErrSuperAdminCannotBePromoted = Error{Code:"AUTH_014", Message:"Super admin cannot be promoted"}
	ErrEmailUpdateFailed = Error{Code:"AUTH_013", Message:"failed to update email"}
	ErrSuperAdminCannotBeDisable = Error{Code:"AUTH_015", Message:"super admin cannot be disable"}
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
