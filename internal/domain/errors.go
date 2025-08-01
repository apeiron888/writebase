package domain


// Error represents a domain error with a code and message.
type Error struct {
    Code    string
    Message string
}

func (e Error) Error() string {
    return e.Message
}


// Comment-related errors
var (
    ErrCommentNotFound      = Error{Code: "COMMENT_001", Message: "Comment not found"}
    ErrInvalidCommentID     = Error{Code: "COMMENT_002", Message: "Invalid comment ID format"}
    ErrEmptyCommentContent  = Error{Code: "COMMENT_003", Message: "Comment content cannot be empty"}
    ErrCommentPermission    = Error{Code: "COMMENT_004", Message: "You do not have permission to modify this comment"}
)

// Reaction-related errors
var (
    ErrReactionNotFound     = Error{Code: "REACTION_001", Message: "Reaction not found"}
    ErrInvalidReactionType  = Error{Code: "REACTION_002", Message: "Invalid reaction type"}
    ErrAlreadyReacted       = Error{Code: "REACTION_003", Message: "Already reacted"}
)

// Follow-related errors
var (
    ErrFollowNotFound       = Error{Code: "FOLLOW_001", Message: "Follow relationship not found"}
    ErrAlreadyFollowing     = Error{Code: "FOLLOW_002", Message: "Already following this user"}
    ErrCannotFollowSelf     = Error{Code: "FOLLOW_003", Message: "Cannot follow yourself"}
)

// Report-related errors
var (
    ErrReportNotFound       = Error{Code: "REPORT_001", Message: "Report not found"}
    ErrInvalidReportTarget  = Error{Code: "REPORT_002", Message: "Invalid report target"}
    ErrReportAlreadyResolved = Error{Code: "REPORT_003", Message: "Report already resolved"}
)


// General errors
var (
    ErrInternalServer   = Error{Code: "GEN_001", Message: "Internal server error"}
    ErrInvalidRequest   = Error{Code: "GEN_002", Message: "Invalid request data"}
    ErrRateLimitExceeded = Error{Code: "GEN_003", Message: "Rate limit exceeded"}
    ErrNotFound         = Error{Code: "GEN_004", Message: "Requested resource not found"}
)
