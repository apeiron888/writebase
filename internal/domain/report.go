package domain

import "context"

type ReportStatus string

const (
	ReportPending  ReportStatus = "pending"
	ReportResolved ReportStatus = "resolved"
)

type Report struct {
	ID         string
	ReporterID string
	TargetID   string
	TargetType string
	Reason     string
	Status     ReportStatus
	CreatedAt  int64
	ResolvedAt *int64
}

type IReportRepository interface {
	CreateReport(ctx context.Context, report *Report) error
	GetReports(ctx context.Context, filter map[string]interface{}) ([]*Report, error)
	UpdateReportStatus(ctx context.Context, reportID string, status ReportStatus) error
}
