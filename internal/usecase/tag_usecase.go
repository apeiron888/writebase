package usecase

import (
	"context"
	"strings"
	"write_base/internal/domain"

)

type TagUsecaseImpl struct {
	tagRepo domain.TagRepository
	utils domain.IUtils
}

func NewTagUsecase(tagRepo domain.TagRepository,utils domain.IUtils) domain.TagUsecase {
	return &TagUsecaseImpl{tagRepo: tagRepo,utils:utils}
}

func (uc *TagUsecaseImpl) CreateTag(ctx context.Context, userID string, name string) (*domain.Tag, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, domain.ErrInvalidTagName
	}

	// Check if tag already exists
	existing, _ := uc.tagRepo.GetByName(ctx, name)
	if existing != nil {
		return nil, domain.ErrTagAlreadyExists
	}

	tagID := uc.utils.GenerateUUID()

	tag := &domain.Tag{
		ID: tagID,
		Name:      name,
		Status:    domain.TagStatusPending,
		CreatedBy: userID,
	}
	return uc.tagRepo.Create(ctx, tag)
}

func (uc *TagUsecaseImpl) ApproveTag(ctx context.Context, tagID string) (*domain.Tag, error) {
	tag, err := uc.tagRepo.GetByID(ctx, tagID)
	if err != nil {
		return nil, err
	}
	tag.Status = domain.TagStatusApproved
	return uc.tagRepo.Update(ctx, tag)
}

func (uc *TagUsecaseImpl) RejectTag(ctx context.Context, tagID string) (*domain.Tag, error) {
	tag, err := uc.tagRepo.GetByID(ctx, tagID)
	if err != nil {
		return nil, err
	}
	tag.Status = domain.TagStatusRejected
	return uc.tagRepo.Update(ctx, tag)
}

func (uc *TagUsecaseImpl) ListTags(ctx context.Context, status domain.TagStatus) ([]domain.Tag, error) {
	return uc.tagRepo.List(ctx, domain.TagFilter{Status: status})
}

func (uc *TagUsecaseImpl) DeleteTag(ctx context.Context, tagID string) error {
	return uc.tagRepo.Delete(ctx, tagID)
}

func (uc *TagUsecaseImpl) IsTagApproved(name string) bool {
	tag, err := uc.tagRepo.GetByName(context.Background(), name)
	return err == nil && tag.Status == domain.TagStatusApproved
}

func (uc *TagUsecaseImpl) ValidateTags(tags []string) error {
	for _, tag := range tags {
		tagStatus, err := uc.tagRepo.GetByName(context.Background(), tag)
		if err != nil {
			return domain.ErrTagNotFound
		}
		if tagStatus.Status == domain.TagStatusRejected {
			return domain.ErrTagRejected
		}
	}
	return nil
}