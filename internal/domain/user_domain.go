package domain

import (
	"context"
	"time"
)

type UserRole string
const (
    RoleUser  UserRole = "user"
    RoleAdmin UserRole = "admin"
    RoleSuperAdmin UserRole = "super_admin"
)
type User struct {
    ID             string    
    Username       string   
    Email          string    
    Password       string    
    Role           UserRole  
    Bio            string
    BookMark       []string    
    ProfileImage   string
    Verified       bool  
    IsActive       bool      
    CreatedAt      time.Time 
    UpdatedAt      time.Time 
}

type RefreshToken struct {
    ID          string     // UUID of this refresh token record
    UserID      string     // UUID of the user this token belongs to
    Token       string     // Actual refresh token (opaque string or signed JWT)
    DeviceInfo  string     // Human-readable description of the device
    IP          string     // IP address where the token was issued
    UserAgent   string     // Full user-agent string from browser/request header
    Revoked     bool       // Has the token been revoked (logout, password change, etc.)
    RevokedAt   *time.Time // Optional: when it was revoked
    ExpiresAt   time.Time  // When this refresh token expires
    CreatedAt   time.Time  // When this token was created
    UpdatedAt   time.Time  // Last time the token record was updated
}

type AuthClaims struct {
    UserID string
    Role   string
}

type IUserRepository interface {
    CreateUser(ctx context.Context, user *User) error
    GetByEmail(ctx context.Context, email string) (*User, error)
    GetByID(ctx context.Context, id string) (*User, error)
    GetByUsername(ctx context.Context, username string) (*User, error)
    UpdateUser(ctx context.Context, user *User) error
    DeleteUser(ctx context.Context, id string) error
    PromoteToAdmin(ctx context.Context, userID string) error
    DemoteToUser(ctx context.Context, userID string) error
// token related
    StoreToken(ctx context.Context, token *RefreshToken) error
    GetByToken(ctx context.Context, token string) (*RefreshToken, error)
    GetValidByUser(ctx context.Context,userID string) ([]*RefreshToken, error)
    RevokeToken(ctx context.Context,token string) error
    DeleteExpiredTokens() error
    RevokeAllByUser(ctx context.Context, userID string) error
    SaveVerificationToken(ctx context.Context,emailToken *EmailVerificationToken) error
    GetVerificationToken(ctx context.Context, emailToken string) (*EmailVerificationToken, error)
}

type IPasswordService interface {
    HashPassword(password string) (string, error)
    VerifyPassword(hash, password string) bool
    IsPasswordStrong(password string) bool 
}

type ITokenService interface {
    GenerateAccessToken(user *User) (string, error)
    GenerateRefreshToken(user *User) (string, error)
    ValidateAccessToken(token string) (*AuthClaims, error)
    ValidateRefreshToken(token string) (*AuthClaims, error)

}

type IEmailService interface {
    SendPasswordReset(email, token string) error
    SendVerificationEmail(email, code string) error
}

type IUserUsecase interface {
	Register(ctx context.Context, registerInput *RegisterInput) error
    VerifyEmail(ctx context.Context, token string) error
	Login(ctx context.Context, loginInput *LoginInput, metadata *AuthMetadata) (*LoginResult, error)
	Logout(ctx context.Context, refreshToken string) error
	RefreshToken(ctx context.Context, refreshToken string, metadata *AuthMetadata) (*LoginResult, error)

	GetProfile(ctx context.Context, userID string) (*User, error)
	UpdateProfile(ctx context.Context, updateProfileInpute *UpdateProfileInput) error
    UpdateAccount(ctx context.Context, updateAccoutInput *UpdateAccountInput)
	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error
    ForgotPassword(ctx context.Context, email string) error
    ResetPassword(ctx context.Context, resetToken, newPassword string) error
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

type AuthMetadata struct {
	IP         string
	UserAgent  string
	DeviceInfo string
}
type RegisterInput struct {
    Username string
    Email    string
    Password string
}

type LoginInput struct {
    EmailOrUsername string
    Password        string
}

type UpdateProfileInput struct {
    UserID       string
    Bio          string
    ProfileImage string
}

type UpdateAccountInput struct {
    UserID   string
    Username string
    Email    string
}

type EmailVerificationToken struct {
    ID        string
    UserID    string
    Token     string
    ExpiresAt time.Time
    CreatedAt time.Time
}