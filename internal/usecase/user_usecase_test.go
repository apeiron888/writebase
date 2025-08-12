package usecase_test

import (
	"context"
	"testing"
	"time"

	"write_base/internal/domain"
	"write_base/internal/mocks"
	"write_base/internal/usecase"

	"github.com/stretchr/testify/mock"
)

func TestUserUsecase_Login_Success(t *testing.T) {
    ctx := context.Background()
    userRepo := mocks.NewMockIUserRepository(t)
    passSvc := mocks.NewMockIPasswordService(t)
    tokenSvc := mocks.NewMockITokenService(t)
    emailSvc := mocks.NewMockIEmailService(t)

    u := &domain.User{ID: "u1", Email: "a@b.com", Username: "alice", Password: "HASH", Verified: true, IsActive: true}
    userRepo.EXPECT().GetByEmail(ctx, "a@b.com").Return(u, nil)
    passSvc.EXPECT().VerifyPassword("HASH", "pw").Return(true)
    tokenSvc.EXPECT().GenerateAccessToken(u).Return("AT", nil)
    tokenSvc.EXPECT().GenerateRefreshToken(u).Return("RT", nil)
    userRepo.EXPECT().StoreToken(ctx, mock.AnythingOfType("*domain.RefreshToken")).Return(nil)

    uc := usecase.NewUserUsecase(userRepo, passSvc, tokenSvc, emailSvc)

    res, err := uc.Login(ctx, &domain.LoginInput{EmailOrUsername: "a@b.com", Password: "pw"}, &domain.AuthMetadata{IP: "1.1.1.1", UserAgent: "ua", DeviceInfo: "dev"})
    if err != nil {
        t.Fatalf("Login error: %v", err)
    }
    if res.AccessToken != "AT" || res.RefreshToken != "RT" {
        t.Fatalf("unexpected tokens: %+v", res)
    }
    if res.ExpiresAt.IsZero() {
        t.Fatalf("expected non-zero ExpiresAt")
    }
}

func TestUserUsecase_Register_Success(t *testing.T) {
    ctx := context.Background()
    userRepo := mocks.NewMockIUserRepository(t)
    passSvc := mocks.NewMockIPasswordService(t)
    tokenSvc := mocks.NewMockITokenService(t)
    emailSvc := mocks.NewMockIEmailService(t)

    // Simulate non-existing user/email by returning errors from GetByEmail/GetByUsername
    userRepo.EXPECT().GetByEmail(ctx, "a@b.com").Return(nil, domain.ErrUserNotFound)
    userRepo.EXPECT().GetByUsername(ctx, "alice").Return(nil, domain.ErrUserNotFound)
    passSvc.EXPECT().IsPasswordStrong("Pw0!abcd").Return(true)
    passSvc.EXPECT().HashPassword("Pw0!abcd").Return("HASH", nil)
    userRepo.EXPECT().SaveVerificationToken(ctx, mock.AnythingOfType("*domain.EmailVerificationToken")).Return(nil)
    emailSvc.EXPECT().SendVerificationEmail("a@b.com", mock.AnythingOfType("string")).Return(nil)
    userRepo.EXPECT().CreateUser(ctx, mock.AnythingOfType("*domain.User")).Return(nil)

    uc := usecase.NewUserUsecase(userRepo, passSvc, tokenSvc, emailSvc)
    err := uc.Register(ctx, &domain.RegisterInput{Username: "alice", Email: "a@b.com", Password: "Pw0!abcd"})
    if err != nil {
        t.Fatalf("Register error: %v", err)
    }
}

func TestUserUsecase_ChangePassword_Success(t *testing.T) {
    ctx := context.Background()
    userRepo := mocks.NewMockIUserRepository(t)
    passSvc := mocks.NewMockIPasswordService(t)
    tokenSvc := mocks.NewMockITokenService(t)
    emailSvc := mocks.NewMockIEmailService(t)

    u := &domain.User{ID: "u1", Password: "OLDHASH", UpdatedAt: time.Time{}}
    userRepo.EXPECT().GetByID(ctx, "u1").Return(u, nil)
    passSvc.EXPECT().VerifyPassword("OLDHASH", "oldpw").Return(true)
    passSvc.EXPECT().IsPasswordStrong("New0!pass").Return(true)
    passSvc.EXPECT().HashPassword("New0!pass").Return("NEWHASH", nil)
    userRepo.EXPECT().UpdateUser(ctx, mock.AnythingOfType("*domain.User")).Return(nil)

    uc := usecase.NewUserUsecase(userRepo, passSvc, tokenSvc, emailSvc)
    if err := uc.ChangePassword(ctx, "u1", "oldpw", "New0!pass"); err != nil {
        t.Fatalf("ChangePassword error: %v", err)
    }
}
