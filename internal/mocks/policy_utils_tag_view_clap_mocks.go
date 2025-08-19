package mocks

import (
	"context"
	"write_base/internal/domain"
)

// Policy mock
type PolicyMock struct {
	UserExistsFn         func(userID string) bool
	ArticleCreateValidFn func(input *domain.Article) bool
	UserOwnsArticleFn    func(userID string, input *domain.Article) bool
	CheckChangesValidFn  func(oldArticle *domain.Article, newArticle *domain.Article) bool
	IsAdminFn            func(userID, userRole string) bool
}

func (m *PolicyMock) UserExists(userID string) bool {
	if m.UserExistsFn != nil {
		return m.UserExistsFn(userID)
	}
	return true
}
func (m *PolicyMock) ArticleCreateValid(input *domain.Article) bool {
	if m.ArticleCreateValidFn != nil {
		return m.ArticleCreateValidFn(input)
	}
	return true
}
func (m *PolicyMock) UserOwnsArticle(userID string, input *domain.Article) bool {
	if m.UserOwnsArticleFn != nil {
		return m.UserOwnsArticleFn(userID, input)
	}
	return true
}
func (m *PolicyMock) CheckArticleChangesAndValid(oldArticle *domain.Article, newArticle *domain.Article) bool {
	if m.CheckChangesValidFn != nil {
		return m.CheckChangesValidFn(oldArticle, newArticle)
	}
	return true
}
func (m *PolicyMock) IsAdmin(userID string, userRole string) bool {
	if m.IsAdminFn != nil {
		return m.IsAdminFn(userID, userRole)
	}
	return false
}

// Utils mock
type UtilsMock struct {
	GenerateUUIDFn      func() string
	GenerateSlugFn      func(string) string
	GenerateShortUUIDFn func() string
	ValidateContentFn   func([]domain.ContentBlock) bool
}

func (u *UtilsMock) GenerateUUID() string {
	if u.GenerateUUIDFn != nil {
		return u.GenerateUUIDFn()
	}
	return "uuid-1"
}
func (u *UtilsMock) GenerateSlug(title string) string {
	if u.GenerateSlugFn != nil {
		return u.GenerateSlugFn(title)
	}
	return "slug-1"
}
func (u *UtilsMock) GenerateShortUUID() string {
	if u.GenerateShortUUIDFn != nil {
		return u.GenerateShortUUIDFn()
	}
	return "x1"
}
func (u *UtilsMock) ValidateContent(blocks []domain.ContentBlock) bool {
	if u.ValidateContentFn != nil {
		return u.ValidateContentFn(blocks)
	}
	return true
}

// Tag usecase mock
type TagUsecaseMock struct {
	ValidateTagsFn  func([]string) error
	IsTagApprovedFn func(string) bool
}

func (t *TagUsecaseMock) CreateTag(ctx context.Context, userID string, name string) (*domain.Tag, error) {
	return nil, nil
}
func (t *TagUsecaseMock) ApproveTag(ctx context.Context, tagID string) (*domain.Tag, error) {
	return nil, nil
}
func (t *TagUsecaseMock) RejectTag(ctx context.Context, tagID string) (*domain.Tag, error) {
	return nil, nil
}
func (t *TagUsecaseMock) ListTags(ctx context.Context, status domain.TagStatus) ([]domain.Tag, error) {
	return nil, nil
}
func (t *TagUsecaseMock) DeleteTag(ctx context.Context, tagID string) error { return nil }
func (t *TagUsecaseMock) IsTagApproved(name string) bool {
	if t.IsTagApprovedFn != nil {
		return t.IsTagApprovedFn(name)
	}
	return true
}
func (t *TagUsecaseMock) ValidateTags(tags []string) error {
	if t.ValidateTagsFn != nil {
		return t.ValidateTagsFn(tags)
	}
	return nil
}

// View usecase mock
type ViewUsecaseMock struct {
	RecordViewFn func(ctx context.Context, userID, articleID, clientIP string) error
}

func (v *ViewUsecaseMock) RecordView(ctx context.Context, userID, articleID, clientIP string) error {
	if v.RecordViewFn != nil {
		return v.RecordViewFn(ctx, userID, articleID, clientIP)
	}
	return nil
}

// Clap usecase mock
type ClapUsecaseMock struct {
	AddClapFn func(ctx context.Context, userID, articleID string) (int, error)
}

func (c *ClapUsecaseMock) AddClap(ctx context.Context, userID, articleID string) (int, error) {
	if c.AddClapFn != nil {
		return c.AddClapFn(ctx, userID, articleID)
	}
	return 0, nil
}
