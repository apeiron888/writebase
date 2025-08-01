package dto

type ReportResponse struct {
 ID         string  `json:"id" bson:"_id,omitempty"`
 ReporterID string  `json:"reporter_id" bson:"reporter_id"`
 TargetID   string  `json:"target_id" bson:"target_id"`
 TargetType string  `json:"target_type" bson:"target_type"`
 Reason     string  `json:"reason" bson:"reason"`
 Status     string  `json:"status" bson:"status"`
 CreatedAt  int64   `json:"created_at" bson:"created_at"`
 ResolvedAt *int64  `json:"resolved_at,omitempty" bson:"resolved_at,omitempty"`
}
