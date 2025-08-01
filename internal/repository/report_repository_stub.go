package repository

import (
	"context"
	"starter/internal/domain"
)

type ReportRepositoryStub struct {
	reports map[string]*domain.Report
}

func NewReportRepositoryStub() *ReportRepositoryStub {
	return &ReportRepositoryStub{reports: make(map[string]*domain.Report)}
}

func (r *ReportRepositoryStub) CreateReport(ctx context.Context, report *domain.Report) error {
	r.reports[report.ID] = report
	return nil
}

func (r *ReportRepositoryStub) GetReports(ctx context.Context, filter map[string]interface{}) ([]*domain.Report, error) {
	var res []*domain.Report
	for _, rep := range r.reports {
		res = append(res, rep)
	}
	return res, nil
}

func (r *ReportRepositoryStub) UpdateReportStatus(ctx context.Context, reportID string, status domain.ReportStatus) error {
	rep, ok := r.reports[reportID]
	if !ok {
		return nil
	}
	rep.Status = status
	return nil
}
