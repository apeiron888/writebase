package domain

// Error represents a domain error with a code and message.
type Error struct {
	Code    string
	Message string
}

func (e Error) Error() string {
	return e.Message
}

// User-related errors
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

// General errors
var (
	ErrInternalServer    = Error{Code: "GEN_001", Message: "Internal server error"}
	ErrInvalidRequest    = Error{Code: "GEN_002", Message: "Invalid request data"}
	ErrRateLimitExceeded = Error{Code: "GEN_003", Message: "Rate limit exceeded"}
	ErrNotFound          = Error{Code: "GEN_004", Message: "Requested resource not found"}
)
