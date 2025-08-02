package usecase

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
	"write_base/internal/domain"

	"github.com/google/uuid"
)



type userUsercase struct{
	userRepo domain.IUserRepository
	passwordService domain.IPasswordService
	tokenService domain.ITokenService
	emailService domain.IEmailService
}

// func NewUserUsecase(repo domain.IUserRepository,pass domain.IPasswordService, tk domain.ITokenService ) domain.IUserUsecase{
// 	return &userUsercase{userRepo: repo, passwordService: pass, tokenService: tk}
// }


func (uu *userUsercase) Register(ctx context.Context, req *domain.RegisterInput) error{
	if _, err := uu.userRepo.GetByEmail(ctx, req.Email); err == nil {
		return domain.ErrEmailAlreadyExists
	}
	if _, err := uu.userRepo.GetByUsername(ctx, req.Username); err == nil {
		return domain.ErrUsernameAlreadyExists
	}

	if !uu.passwordService.IsPasswordStrong(req.Password){
		return domain.ErrWeakPassword

	}
	hashedPassword, err := uu.passwordService.HashPassword(req.Password)
	if err != nil{
		return err
	}
	user := &domain.User{
		ID:   uuid.New().String(),
		Username:req.Username,
		Email: req.Email,
		Password: hashedPassword,
		Role: domain.RoleUser,
		IsActive: true,
		Verified: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = uu.userRepo.CreateUser(ctx, user)
	if err != nil{
		return err
	}
	// 2. Generate email verification token
	token := uuid.New().String()
	verificationToken := &domain.EmailVerificationToken{
		ID: uuid.New().String(),
		UserID: user.ID,
		Token: token,
		ExpiresAt: time.Now().Add(24 * time.Hour), // expires in 1 day
		CreatedAt: time.Now(),
	}
	err = uu.userRepo.SaveVerificationToken(ctx, verificationToken)
	if err != nil {
		return err
	}
	frontendURL := os.Getenv("FRONTEND_BASE_URL") // e.g. "http://localhost:3000"
	verificationURL := fmt.Sprintf("%s/verify-email?token=%s", frontendURL, token)
	err = uu.emailService.SendVerificationEmail(user.Email, verificationURL)
	if err != nil {
		return err
	}

	return nil
}

func (uu *userUsercase) VerifyEmail(ctx context.Context, emailToken string) error {
    verification, err := uu.userRepo.GetVerificationToken(ctx, emailToken)
    if err != nil {
        return domain.ErrInvalidToken
    }
	if time.Now().After(verification.ExpiresAt) {
        return domain.ErrExpiredToken
    }

    user, err := uu.userRepo.GetByID(ctx, verification.UserID)
    if err != nil {
        return err
    }

    user.Verified = true
    user.UpdatedAt = time.Now()

    return uu.userRepo.UpdateUser(ctx, user)
}

func(uu *userUsercase) Login(ctx context.Context, req *domain.LoginInput) (*domain.LoginResult, error){
	var existingUser *domain.User
	var err error
    if strings.Contains(req.EmailOrUsername, "@") {
        existingUser, err = uu.userRepo.GetByEmail(ctx, req.EmailOrUsername)
	}else {
        existingUser, err = uu.userRepo.GetByUsername(ctx, req.EmailOrUsername)
	}
	if err != nil{
		return nil, domain.ErrUserNotFound
    }
	if !existingUser.Verified{
		return nil,domain.ErrUserNotVerified
	}
	if !existingUser.IsActive {
		return nil, domain.ErrUserDeactivated
	}

	if !uu.passwordService.VerifyPassword(existingUser.Password, req.Password){
		return nil, domain.ErrInvalidCredentials 
	}
	accessToken, err := uu.tokenService.GenerateAccessToken(existingUser)
	if err != nil {
		return nil, err
	}

	refreshToken, err := uu.tokenService.GenerateRefreshToken(existingUser)
	if err != nil {
		return nil, err 

	}
	loginResult := &domain.LoginResult{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}
	
	return loginResult, nil
}