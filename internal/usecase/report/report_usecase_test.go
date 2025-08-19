package usecase

import (
	"context"
	"testing"
	"write_base/internal/domain"
)

type reportRepoMock struct {
	CreateReportFn       func(ctx context.Context, r *domain.Report) error
	GetReportsFn         func(ctx context.Context, f map[string]interface{}) ([]*domain.Report, error)
	UpdateReportStatusFn func(ctx context.Context, id string, s domain.ReportStatus) error
}

func (m *reportRepoMock) CreateReport(ctx context.Context, r *domain.Report) error {
	return m.CreateReportFn(ctx, r)
}
func (m *reportRepoMock) GetReports(ctx context.Context, f map[string]interface{}) ([]*domain.Report, error) {
	return m.GetReportsFn(ctx, f)
}
func (m *reportRepoMock) UpdateReportStatus(ctx context.Context, id string, s domain.ReportStatus) error {
	return m.UpdateReportStatusFn(ctx, id, s)
}

func TestReportService_Basic(t *testing.T) {
	repo := &reportRepoMock{
		CreateReportFn: func(ctx context.Context, r *domain.Report) error { return nil },
		GetReportsFn: func(ctx context.Context, f map[string]interface{}) ([]*domain.Report, error) {
			return []*domain.Report{{ID: "rep1"}}, nil
		},
		UpdateReportStatusFn: func(ctx context.Context, id string, s domain.ReportStatus) error { return nil },
	}
	s := NewReportService(repo)
	if err := s.CreateReport(context.Background(), &domain.Report{ID: "rep1"}); err != nil {
		t.Fatal(err)
	}
	if list, err := s.GetReports(context.Background(), map[string]interface{}{}); err != nil || len(list) != 1 {
		t.Fatalf("get bad")
	}
	if err := s.UpdateReportStatus(context.Background(), "rep1", domain.ReportResolved); err != nil {
		t.Fatal(err)
	}
}
