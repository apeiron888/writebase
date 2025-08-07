package repository

import (
	"write_base/internal/domain"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DTOs for MongoDB serialization (infrastructure layer only)
type userDTO struct {
	ID           string    `bson:"_id"`
	Username     string    `bson:"username"`
	Email        string    `bson:"email"`
	Password     string    `bson:"password"`
	Role         string    `bson:"role"`
	Bio          string    `bson:"bio"`
	BookMark     []string  `bson:"bookmark"`
	ProfileImage string    `bson:"profile_image"`
	Verified     bool      `bson:"verified"`
	IsActive     bool      `bson:"is_active"`
	CreatedAt    primitive.DateTime `bson:"created_at"`
	UpdatedAt    primitive.DateTime `bson:"updated_at"`
}

type refreshTokenDTO struct {
	ID         string    `bson:"_id"`
	UserID     string    `bson:"user_id"`
	Token      string    `bson:"token"`
	DeviceInfo string    `bson:"device_info"`
	IP         string    `bson:"ip"`
	UserAgent  string    `bson:"user_agent"`
	Revoked    bool      `bson:"revoked"`
	RevokedAt  *primitive.DateTime `bson:"revoked_at,omitempty"`
	ExpiresAt  primitive.DateTime  `bson:"expires_at"`
	CreatedAt  primitive.DateTime  `bson:"created_at"`
	UpdatedAt  primitive.DateTime  `bson:"updated_at"`
}

type emailVerificationTokenDTO struct {
	ID        string    `bson:"_id"`
	UserID    string    `bson:"user_id"`
	Token     string    `bson:"token"`
	ExpiresAt primitive.DateTime `bson:"expires_at"`
	CreatedAt primitive.DateTime `bson:"created_at"`
}
// Conversion functions between DTOs and domain models

func userDTOToDomain(u userDTO) *domain.User {
	return &domain.User{
		ID:           u.ID,
		Username:     u.Username,
		Email:        u.Email,
		Password:     u.Password,
		Role:         domain.UserRole(u.Role),
		Bio:          u.Bio,
		BookMark:     u.BookMark,
		ProfileImage: u.ProfileImage,
		Verified:     u.Verified,
		IsActive:     u.IsActive,
		CreatedAt:    u.CreatedAt.Time(),
		UpdatedAt:    u.UpdatedAt.Time(),
	}
}

func userDomainToDTO(u *domain.User) userDTO {
	return userDTO{
		ID:           u.ID,
		Username:     u.Username,
		Email:        u.Email,
		Password:     u.Password,
		Role:         string(u.Role),
		Bio:          u.Bio,
		BookMark:     u.BookMark,
		ProfileImage: u.ProfileImage,
		Verified:     u.Verified,
		IsActive:     u.IsActive,
		CreatedAt:    primitive.NewDateTimeFromTime(u.CreatedAt),
		UpdatedAt:    primitive.NewDateTimeFromTime(u.UpdatedAt),
	}
}

func refreshTokenDTOToDomain(r refreshTokenDTO) *domain.RefreshToken {
	var revokedAt *time.Time
	if r.RevokedAt != nil {
		t := r.RevokedAt.Time()
		revokedAt = &t
	}
	return &domain.RefreshToken{
		ID:         r.ID,
		UserID:     r.UserID,
		Token:      r.Token,
		DeviceInfo: r.DeviceInfo,
		IP:         r.IP,
		UserAgent:  r.UserAgent,
		Revoked:    r.Revoked,
		RevokedAt:  revokedAt,
		ExpiresAt:  r.ExpiresAt.Time(),
		CreatedAt:  r.CreatedAt.Time(),
		UpdatedAt:  r.UpdatedAt.Time(),
	}
}

func refreshTokenDomainToDTO(r *domain.RefreshToken) refreshTokenDTO {
	var revokedAt *primitive.DateTime
	if r.RevokedAt != nil {
		t := primitive.NewDateTimeFromTime(*r.RevokedAt)
		revokedAt = &t
	}
	return refreshTokenDTO{
		ID:         r.ID,
		UserID:     r.UserID,
		Token:      r.Token,
		DeviceInfo: r.DeviceInfo,
		IP:         r.IP,
		UserAgent:  r.UserAgent,
		Revoked:    r.Revoked,
		RevokedAt:  revokedAt,
		ExpiresAt:  primitive.NewDateTimeFromTime(r.ExpiresAt),
		CreatedAt:  primitive.NewDateTimeFromTime(r.CreatedAt),
		UpdatedAt:  primitive.NewDateTimeFromTime(r.UpdatedAt),
	}
}

func emailVerificationTokenDTOToDomain(e emailVerificationTokenDTO) *domain.EmailVerificationToken {
	return &domain.EmailVerificationToken{
		ID:        e.ID,
		UserID:    e.UserID,
		Token:     e.Token,
		ExpiresAt: e.ExpiresAt.Time(),
		CreatedAt: e.CreatedAt.Time(),
	}
}

func emailVerificationTokenDomainToDTO(e *domain.EmailVerificationToken) emailVerificationTokenDTO {
	return emailVerificationTokenDTO{
		ID:        e.ID,
		UserID:    e.UserID,
		Token:     e.Token,
		ExpiresAt: primitive.NewDateTimeFromTime(e.ExpiresAt),
		CreatedAt: primitive.NewDateTimeFromTime(e.CreatedAt),
	}
}