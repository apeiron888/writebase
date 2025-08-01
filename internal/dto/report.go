package dto

type ReportRequest struct {
	ReporterID string `json:"reporter_id"`
	TargetID   string `json:"target_id"`
	TargetType string `json:"target_type"`
	Reason     string `json:"reason"`
}

type ReportResponse struct {
	ID         string  `json:"id"`
	ReporterID string  `json:"reporter_id"`
	TargetID   string  `json:"target_id"`
	TargetType string  `json:"target_type"`
	Reason     string  `json:"reason"`
	Status     string  `json:"status"`
	CreatedAt  int64   `json:"created_at"`
	ResolvedAt *int64  `json:"resolved_at,omitempty"`
}
