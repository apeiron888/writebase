package usecase

import (
	"context"
	"write_base/internal/domain"
)

type ReportService struct {
	repo domain.IReportRepository
}

func NewReportService(repo domain.IReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) CreateReport(ctx context.Context, report *domain.Report) error {
	return s.repo.CreateReport(ctx, report)
}

func (s *ReportService) GetReports(ctx context.Context, filter map[string]interface{}) ([]*domain.Report, error) {
	return s.repo.GetReports(ctx, filter)
}

func (s *ReportService) UpdateReportStatus(ctx context.Context, reportID string, status domain.ReportStatus) error {
	return s.repo.UpdateReportStatus(ctx, reportID, status)
}
