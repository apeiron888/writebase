package dto

type ReportRequest struct {
	ReporterID string `json:"reporter_id"`
	TargetID   string `json:"target_id"`
	TargetType string `json:"target_type"`
	Reason     string `json:"reason"`
}

type UpdateReportStatusRequest struct {
	Status string `json:"status"`
}