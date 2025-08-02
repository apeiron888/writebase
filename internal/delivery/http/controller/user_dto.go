package controller

import (
	"time"
	"write_base/internal/domain"
)

// --- Auth DTOs ---

// RegisterRequest is used for user registration
type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

//rigister user to domain user
func (r *RegisterRequest) ToRegisterInput() *domain.RegisterInput {
    return &domain.RegisterInput{
        Username: r.Username,
        Email:    r.Email,
        Password: r.Password, 
    }
}

// LoginRequest is used for user login
type LoginRequest struct {
    EmailOrUsername string `json:"email_or_username" binding:"required"`
    Password        string `json:"password" binding:"required"`
}
// login uer to domain user 
func (l *LoginRequest) ToLoginInput() *domain.LoginInput{
    return &domain.LoginInput{
        EmailOrUsername: l.EmailOrUsername,
        Password: l.Password,
    }
}

// LoginResponse returns tokens after login
type LoginResponse struct {
    AccessToken  string    `json:"access_token"`
    RefreshToken string    `json:"refresh_token"`
    ExpiresAt    time.Time `json:"expires_at"` // access token expiration
}

// logout request need refresh token
type LogoutRequest struct {
    RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenRequest is used to request new access token
type RefreshTokenRequest struct {
    RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenResponse returns new access token after refresh
type RefreshTokenResponse struct {
    AccessToken string    `json:"access_token"`
    ExpiresAt   time.Time `json:"expires_at"` // new access token expiration
}

// ForgotPasswordRequest for sending password reset email
type ForgotPasswordRequest struct {
    Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest for resetting password using a token
type ResetPasswordRequest struct {
    Token       string `json:"token" binding:"required"`
    NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ChangePasswordRequest for changing password after login
type ChangePasswordRequest struct {
    OldPassword string `json:"old_password" binding:"required"`
    NewPassword string `json:"new_password" binding:"required,min=6"`
}

// --- User Profile DTOs ---

// UserResponse returns user profile data
type UserResponse struct {
    ID           string    `json:"id"`
    Username     string    `json:"username"`
    Email        string    `json:"email"`
    Bio          string    `json:"bio,omitempty"`
    ProfileImage string    `json:"profile_image,omitempty"`
    Role         string    `json:"role"`
    CreatedAt    time.Time `json:"created_at"`
}

// UpdateProfileRequest for updating user profile
type UpdateProfileRequest struct {
    Bio          *string `json:"bio,omitempty"`
    ProfileImage *string `json:"profile_image,omitempty"`
}

// --- Optional Admin DTOs ---

// PromoteDemoteUserRequest - for promote/demote if needed
type PromoteDemoteUserRequest struct {
    // optional: reason string `json:"reason,omitempty"`
}


//--- convert functions