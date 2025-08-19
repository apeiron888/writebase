package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	dtodlv "write_base/internal/delivery/http/controller/dto"
	"write_base/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// fake usecase implementing domain.IReportUsecase
type fakeReportUC struct {
	err  error
	list []*domain.Report
}

func (f *fakeReportUC) CreateReport(_ context.Context, _ *domain.Report) error { return f.err }
func (f *fakeReportUC) GetReports(_ context.Context, _ map[string]interface{}) ([]*domain.Report, error) {
	return f.list, f.err
}
func (f *fakeReportUC) UpdateReportStatus(_ context.Context, _ string, _ domain.ReportStatus) error {
	return f.err
}

func setupReportRouter(uc domain.IReportUsecase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewReportController(uc)
	r.POST("/reports", h.CreateReport)
	r.GET("/reports", h.GetReports)
	r.PUT("/reports/:id/status", h.UpdateReportStatus)
	return r
}

func TestReportController_HappyPaths(t *testing.T) {
	uc := &fakeReportUC{list: []*domain.Report{{ID: "r1"}}}
	r := setupReportRouter(uc)

	// Create
	body, _ := json.Marshal(dtodlv.ReportRequest{ReporterID: "u1", TargetID: "p1", TargetType: "post", Reason: "spam"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/reports", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	// Get all
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/reports", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	// Update status
	up, _ := json.Marshal(dtodlv.UpdateReportStatusRequest{Status: "resolved"})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPut, "/reports/r1/status", bytes.NewReader(up))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
