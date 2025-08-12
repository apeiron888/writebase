package usecase

import (
	"context"
	"write_base/internal/domain"
)

type ClapUsecaseImpl struct {
	clapRepo domain.ClapRepository
	utils domain.IUtils
}

func NewClapUsecase(clapRepo domain.ClapRepository,u domain.IUtils) domain.ClapUsecase {
	return &ClapUsecaseImpl{clapRepo: clapRepo, utils: u}
}

func (uc *ClapUsecaseImpl) AddClap(ctx context.Context, userID, articleID string) (int, error) {
	// Check if user has existing clap record
	clap, err := uc.clapRepo.GetByUserAndArticle(ctx, userID, articleID)
	if err != nil {
		return 0, err
	}
	
	// Create new clap if doesn't exist
	if clap == nil {
		clap = &domain.Clap{
			ID: uc.utils.GenerateUUID(),
			UserID:    userID,
			ArticleID: articleID,
			Count:     1,
		}
		if err := uc.clapRepo.Create(ctx, clap); err != nil {
			return 0, err
		}
	} else {
		// Check clap limit
		if clap.Count >= domain.MaxClapsPerUser {
			return 0, domain.ErrClapLimitExceeded
		}
		
		// Increment existing clap
		clap.Count++
		if err := uc.clapRepo.Update(ctx, clap); err != nil {
			return 0, err
		}
	}
	
	// Get total clap count for article
	total, err := uc.clapRepo.GetArticleClapCount(ctx, articleID)
	if err != nil {
		return 0, err
	}
	
	return total, nil
}