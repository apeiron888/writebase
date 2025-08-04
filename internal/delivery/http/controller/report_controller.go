package controller

import (
	"net/http"
	dtodlv "write_base/internal/delivery/http/controller/dto"
	"write_base/internal/domain"

	"github.com/gin-gonic/gin"
)

type ReportController struct {
	usecase domain.IReportUsecase
}

func NewReportController(usecase domain.IReportUsecase) *ReportController {
	return &ReportController{usecase: usecase}
}

func (rc *ReportController) CreateReport(c *gin.Context) {
	   var req dtodlv.ReportRequest
	   if err := c.ShouldBindJSON(&req); err != nil {
			   c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			   return
	   }
	   report := &domain.Report{
			   ReporterID: req.ReporterID,
			   TargetID:   req.TargetID,
			   TargetType: req.TargetType,
			   Reason:     req.Reason,
	   }
	   if err := rc.usecase.CreateReport(c, report); err != nil {
			   c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			   return
	   }
	   c.JSON(http.StatusCreated, gin.H{"message": "Report created"})
}

func (rc *ReportController) GetReports(c *gin.Context) {
	   // For simplicity, no filter from query params
	   reports, err := rc.usecase.GetReports(c, map[string]interface{}{})
	   if err != nil {
			   c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			   return
	   }
	   c.JSON(http.StatusOK, reports)
}

func (rc *ReportController) UpdateReportStatus(c *gin.Context) {
	   id := c.Param("id")
	   var req dtodlv.UpdateReportStatusRequest
	   if err := c.ShouldBindJSON(&req); err != nil {
			   c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			   return
	   }
	   if err := rc.usecase.UpdateReportStatus(c, id, domain.ReportStatus(req.Status)); err != nil {
			   c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			   return
	   }
	   c.JSON(http.StatusOK, gin.H{"message": "Report status updated"})
}
