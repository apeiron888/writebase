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
    ProfileImage   string    
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
}

type ITokenRepository interface {
    StoreToken(token *RefreshToken) error
    GetByToken(token string) (*RefreshToken, error)
    GetValidByUser(userID string) ([]*RefreshToken, error)
    RevokeToken(token string) error
    DeleteExpiredTokens() error
    RevokeAllByUser(userID string) error
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

type IEmailSender interface {
    SendPasswordReset(email, token string) error
    SendVerificationEmail(email, code string) error
}

type IUserUsecase interface {
	Register(ctx context.Context, user *User) error
	Login(ctx context.Context, user *User, metadata *AuthMetadata) (*LoginResult, error)
	Logout(ctx context.Context, refreshToken string) error
	RefreshToken(ctx context.Context, refreshToken string, metadata *AuthMetadata) (*LoginResult, error)

	GetProfile(ctx context.Context, userID string) (*User, error)
	UpdateProfile(ctx context.Context, userID string, bio, profileImage string) error
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
