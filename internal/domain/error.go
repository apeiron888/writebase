package domain

type Error struct {
	Code    string
	Message string
}

func (e Error) Error() string {
	return e.Message
}

var (
	//General
	ErrInternalServer = Error{Code: "GEN001", Message: "Internal server error"}

	// User
	ErrUnauthorized = Error{Code: "USER001", Message: "User Unauthorized"}

	// Article
	ErrInvalidArticlePayload = Error{Code: "ARTICLE001", Message: "Invalid Article Payload Request"}
	ErrArticleNotFound       = Error{Code: "ARTICLE_002", Message: "Article not found"}
	ErrArticleInvalidID      = Error{Code: "ARTICLE_003", Message: "Article ID is required"}
	ErrArticlePublished      = Error{Code: "ARTICLE_004", Message: "Article already published"}
	ErrArticleNotPublished   = Error{Code: "ARTICLE_005", Message: "Article is not published"}
	ErrArticleArchived       = Error{Code: "ARTICLE_006", Message: "Article is already archived"}
	ErrArticleNotArchived    = Error{Code: "ARTICLE_007", Message: "Article is not archived"}
	ErrArticleInvalidSlug    = Error{Code: "ARTICLE_008", Message: "Invalid article slug"}
	ErrAuthorNotFound        = Error{Code: "ARTICLE_009", Message: "Author not found"}
	ErrArticleContentEmpty   = Error{Code: "ARTICLE_010", Message: "empty content"}
	// Tag
	ErrTagNotFound      = Error{Code: "TAG001", Message: "Tag not found"}
	ErrInvalidTagName   = Error{Code: "TAG002", Message: "Invalid tag name"}
	ErrTagAlreadyExists = Error{Code: "TAG003", Message: "Tag already exists"}
	ErrTagRejected      = Error{Code: "TAG004", Message: "Tag rejected"}
	ErrUnapprovedTags   = Error{Code: "TAG005", Message: "article contains unapproved tags"}

	// CLAP
	ErrClapLimitExceeded = Error{Code: "CLAP001", Message: "clap limit exceeded"}

	ErrContentPolicyViolation = Error{Code: "GEN002",Message: "violation"}
)
