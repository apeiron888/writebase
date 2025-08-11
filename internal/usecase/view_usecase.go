package usecase

import (
	"context"
	"write_base/internal/domain"
)

type ViewUsecaseImpl struct {
	viewRepo domain.ViewRepository
	utils domain.IUtils
}

func NewViewUsecase(viewRepo domain.ViewRepository, u domain.IUtils) domain.ViewUsecase {
	return &ViewUsecaseImpl{viewRepo: viewRepo, utils: u}
}

func (uc *ViewUsecaseImpl) RecordView(ctx context.Context, userID, articleID, clientIP string) error {
	view := &domain.View{
		ID: uc.utils.GenerateUUID(),
		UserID:    userID,
		ArticleID: articleID,
		ClientIP:  clientIP,
	}
	return uc.viewRepo.Create(ctx, view)
}