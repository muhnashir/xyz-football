package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xyz-corp/xyz-football-api/internal/service"
	"github.com/xyz-corp/xyz-football-api/pkg/response"
)

type ReportHandler struct {
	svc service.ReportService
}

func NewReportHandler(svc service.ReportService) *ReportHandler {
	return &ReportHandler{svc: svc}
}

// ListMatchReports godoc
// @Summary      Report seluruh pertandingan
// @Tags         Reports
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} response.Envelope{data=[]dto.MatchReportResponse}
// @Router       /reports/matches [get]
func (h *ReportHandler) List(c *gin.Context) {
	reports, meta, err := h.svc.MatchReportList(c)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.Paginated(c, "Report pertandingan berhasil diambil", reports, meta)
}

// MatchReport godoc
// @Summary      Report satu pertandingan
// @Tags         Reports
// @Security     BearerAuth
// @Produce      json
// @Param        uuid path string true "UUID Pertandingan"
// @Success      200 {object} response.Envelope{data=dto.MatchReportResponse}
// @Router       /reports/matches/{uuid} [get]
func (h *ReportHandler) Detail(c *gin.Context) {
	id, ok := paramUUID(c, "uuid")
	if !ok {
		return
	}
	report, err := h.svc.MatchReport(id)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, "Report pertandingan berhasil diambil", report)
}
