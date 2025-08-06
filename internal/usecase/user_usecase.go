package usecase

import (
	"context"
	"strings"
	"time"
	"write_base/internal/domain"

	"github.com/google/uuid"
)

type UserUsercase struct {
	userRepo        domain.IUserRepository
	passwordService domain.IPasswordService
	tokenService    domain.ITokenService
	emailService    domain.IEmailService
}

// func NewUserUsecase(repo domain.IUserRepository,pass domain.IPasswordService, tk domain.ITokenService , em domain.IEmailService) domain.IUserUsecase{
// 	return &UserUsercase{userRepo: repo, passwordService: pass, tokenService: tk, emailService: em}
// }

func (uu *UserUsercase) Register(ctx context.Context, req *domain.RegisterInput) error {
	if _, err := uu.userRepo.GetByEmail(ctx, req.Email); err == nil {
		return domain.ErrEmailAlreadyExists
	}
	if _, err := uu.userRepo.GetByUsername(ctx, req.Username); err == nil {
		return domain.ErrUsernameAlreadyExists
	}

	if !uu.passwordService.IsPasswordStrong(req.Password) {
		return domain.ErrWeakPassword

	}
	hashedPassword, err := uu.passwordService.HashPassword(req.Password)
	if err != nil {
		return err
	}
	user := &domain.User{
		ID:        uuid.New().String(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		Role:      domain.RoleUser,
		IsActive:  true,
		Verified:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = uu.userRepo.CreateUser(ctx, user)
	if err != nil {
		return err
	}
	// 2. Generate email verification token
	token := uuid.New().String()
	verificationToken := &domain.EmailVerificationToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(5 * time.Minute), 
		CreatedAt: time.Now(),
	}
	err = uu.userRepo.SaveVerificationToken(ctx, verificationToken)
	if err != nil {
		return err
	}
	err = uu.emailService.SendVerificationEmail(user.Email, token)
	if err != nil {
		return err
	}

	return nil
}

func (uu *UserUsercase) VerifyEmail(ctx context.Context, emailToken string) error {
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

func (uu *UserUsercase) Login(ctx context.Context, req *domain.LoginInput, metadata *domain.AuthMetadata) (*domain.LoginResult, error) {
	var existingUser *domain.User
	var err error
	if strings.Contains(req.EmailOrUsername, "@") {
		existingUser, err = uu.userRepo.GetByEmail(ctx, req.EmailOrUsername)
	} else {
		existingUser, err = uu.userRepo.GetByUsername(ctx, req.EmailOrUsername)
	}
	if err != nil {
		return nil, domain.ErrUserNotFound
	}
	if !existingUser.Verified {
		return nil, domain.ErrUserNotVerified
	}
	if !existingUser.IsActive {
		return nil, domain.ErrUserDeactivated
	}

	if !uu.passwordService.VerifyPassword(existingUser.Password, req.Password) {
		return nil, domain.ErrInvalidCredentials
	}
	accessToken, err := uu.tokenService.GenerateAccessToken(existingUser)
	if err != nil {
		return nil, err
	}

	refreshTokenString, err := uu.tokenService.GenerateRefreshToken(existingUser)
	if err != nil {
		return nil, err
	}
	refreshToken := &domain.RefreshToken{
		ID:         uuid.New().String(),
		UserID:     existingUser.ID,
		Token:      refreshTokenString,
		DeviceInfo: metadata.DeviceInfo,
		IP:         metadata.IP,
		UserAgent:  metadata.UserAgent,
		Revoked:    false,
		ExpiresAt:  time.Now().Add(7 * 24 * time.Hour), // Same as token expiry
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	
	if err := uu.userRepo.StoreToken(ctx, refreshToken); err != nil {
		return nil, err
	}
	loginResult := &domain.LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
		ExpiresAt:    time.Now().Add(15 * time.Minute),
	}

	return loginResult, nil
}

func (uu *UserUsercase) LoginOrRegisterOAuthUser(ctx context.Context, registerInput *domain.RegisterInput, metadata *domain.AuthMetadata) (*domain.LoginResult, error) {
	existingUser, err := uu.userRepo.GetByEmail(ctx, registerInput.Email)
	var user domain.User
	if err != nil {
		// User does not exist — register them
		user.Username = registerInput.Username
		user.Email = registerInput.Email
		user.ID = uuid.New().String()
		user.Role = domain.RoleUser
		user.IsActive = true
		user.Verified = true
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()

		if err := uu.userRepo.CreateUser(ctx, &user); err != nil {
			return nil, err
		}
		existingUser = &user
	}
	accessToken, err := uu.tokenService.GenerateAccessToken(existingUser)
	if err != nil {
		return nil, err
	}

	refreshTokenString, err := uu.tokenService.GenerateRefreshToken(existingUser)
	if err != nil {
		return nil, err
	}
	refreshToken := &domain.RefreshToken{
		ID:         uuid.New().String(),
		UserID:     existingUser.ID,
		Token:      refreshTokenString,
		DeviceInfo: metadata.DeviceInfo,
		IP:         metadata.IP,
		UserAgent:  metadata.UserAgent,
		Revoked:    false,
		ExpiresAt:  time.Now().Add(7 * 24 * time.Hour), // Same as token expiry
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	
	if err := uu.userRepo.StoreToken(ctx, refreshToken); err != nil {
		return nil, err
	}
	loginResult := &domain.LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
		ExpiresAt:    time.Now().Add(15 * time.Minute),
	}

	return loginResult, nil

}

func(uu *UserUsercase) RefreshToken(ctx context.Context, refreshTokenString string) (*domain.LoginResult, error){
	refreshToken, err := uu.userRepo.GetByToken(ctx, refreshTokenString)
	if err != nil{
		return nil , err
	}
	if time.Now().After(refreshToken.ExpiresAt){
		return nil, domain.ErrRefreshTokenExpired
	}
	if refreshToken.Revoked{
		return nil, domain.ErrRefreshTokenRevoked
	}
	user, err:= uu.userRepo.GetByID(ctx,refreshToken.UserID)
	if err != nil{
		return nil, err
	}
	accesToken, err:= uu.tokenService.GenerateAccessToken(user)
	if err != nil{
		return nil, err
	}
	loginResult := &domain.LoginResult{
		AccessToken: accesToken,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}
	return loginResult, nil

}

func(uu *UserUsercase) ForgotPassword(ctx context.Context, email string) error{
	user, err:= uu.userRepo.GetByEmail(ctx, email)
	if err != nil{
		return domain.ErrEmailNotRegistered
	}
	tokenStirng := uuid.New().String()
	resetToken := &domain.EmailVerificationToken{
		ID: uuid.New().String(),
		UserID: user.ID,
		Token: tokenStirng,
		ExpiresAt: time.Now().Add(5 *time.Minute),
		CreatedAt: time.Now(),
	}
	err = uu.userRepo.SaveVerificationToken(ctx, resetToken)
	if err != nil{
		return err
	}
	err = uu.emailService.SendPasswordReset(email, tokenStirng)
	if err != nil{
		return err
	}
	return nil
	
}

