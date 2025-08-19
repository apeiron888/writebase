package di

import (
	"context"
	"testing"
	"time"
	"write_base/internal/domain"
)

type fakeUserRepo struct {
	created *domain.User
	exists  *domain.User
}

func (f *fakeUserRepo) CreateUser(ctx context.Context, u *domain.User) error {
	f.created = u
	return nil
}
func (f *fakeUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return f.exists, nil
}

// unused methods
func (f *fakeUserRepo) GetByID(context.Context, string) (*domain.User, error)       { return nil, nil }
func (f *fakeUserRepo) GetByUsername(context.Context, string) (*domain.User, error) { return nil, nil }
func (f *fakeUserRepo) UpdateUser(context.Context, *domain.User) error              { return nil }
func (f *fakeUserRepo) DeleteUser(context.Context, string) error                    { return nil }
func (f *fakeUserRepo) PromoteToAdmin(context.Context, string) error                { return nil }
func (f *fakeUserRepo) DemoteToUser(context.Context, string) error                  { return nil }
func (f *fakeUserRepo) StoreToken(context.Context, *domain.RefreshToken) error      { return nil }
func (f *fakeUserRepo) GetByToken(context.Context, string) (*domain.RefreshToken, error) {
	return nil, nil
}
func (f *fakeUserRepo) GetValidByUser(context.Context, string) ([]*domain.RefreshToken, error) {
	return nil, nil
}
func (f *fakeUserRepo) RevokeToken(context.Context, string) error     { return nil }
func (f *fakeUserRepo) DeleteExpiredTokens() error                    { return nil }
func (f *fakeUserRepo) RevokeAllByUser(context.Context, string) error { return nil }
func (f *fakeUserRepo) SaveVerificationToken(context.Context, *domain.EmailVerificationToken) error {
	return nil
}
func (f *fakeUserRepo) GetVerificationToken(context.Context, string) (*domain.EmailVerificationToken, error) {
	return nil, nil
}
func (f *fakeUserRepo) DeleteUnverifiedExpiredUsers(context.Context, time.Duration) error { return nil }
func (f *fakeUserRepo) DeleteOldRevokedTokens(context.Context, time.Duration) error       { return nil }
func (f *fakeUserRepo) DeleteVerificationToken(context.Context, string) error             { return nil }

type fakePassword struct{}

func (fakePassword) HashPassword(p string) (string, error)     { return "hash", nil }
func (fakePassword) VerifyPassword(hash, password string) bool { return true }
func (fakePassword) IsPasswordStrong(password string) bool     { return true }

func TestSeedSuperAdmin_CreatesWhenMissing(t *testing.T) {
	repo := &fakeUserRepo{}
	if err := SeedSuperAdmin(context.Background(), repo, fakePassword{}); err != nil {
		t.Fatal(err)
	}
	if repo.created == nil || repo.created.Role != "super_admin" || !repo.created.Verified {
		t.Fatal("user not created properly")
	}
}

func TestSeedSuperAdmin_NoCreateWhenExists(t *testing.T) {
	repo := &fakeUserRepo{exists: &domain.User{ID: "1"}}
	if err := SeedSuperAdmin(context.Background(), repo, fakePassword{}); err != nil {
		t.Fatal(err)
	}
	if repo.created != nil {
		t.Fatal("should not create when exists")
	}
}
