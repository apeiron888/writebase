package usecase_test

import (
	"context"
	"testing"
	"write_base/internal/domain"
	"write_base/internal/mocks"
	usecase "write_base/internal/usecase/report"
)

func TestReportService_Basic(t *testing.T) {
	ctx := context.Background()
	repo := mocks.NewMockIReportRepository(t)
	svc := usecase.NewReportService(repo)

	report := &domain.Report{ID: "rep1", ReporterID: "u1", TargetID: "p1", TargetType: "post", Reason: "spam"}
	repo.EXPECT().CreateReport(ctx, report).Return(nil)
	if err := svc.CreateReport(ctx, report); err != nil { t.Fatalf("CreateReport: %v", err) }

	repo.EXPECT().GetReports(ctx, map[string]interface{}{"status": domain.ReportPending}).Return([]*domain.Report{report}, nil)
	list, err := svc.GetReports(ctx, map[string]interface{}{"status": domain.ReportPending})
	if err != nil || len(list) != 1 { t.Fatalf("GetReports: %v %v", list, err) }

	repo.EXPECT().UpdateReportStatus(ctx, "rep1", domain.ReportResolved).Return(nil)
	if err := svc.UpdateReportStatus(ctx, "rep1", domain.ReportResolved); err != nil { t.Fatalf("UpdateReportStatus: %v", err) }
}
