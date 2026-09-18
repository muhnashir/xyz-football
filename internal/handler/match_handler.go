package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xyz-corp/xyz-football-api/internal/dto"
	"github.com/xyz-corp/xyz-football-api/internal/service"
	"github.com/xyz-corp/xyz-football-api/pkg/response"
)

type MatchHandler struct {
	svc      service.MatchService
	eventSvc service.MatchEventService
}

func NewMatchHandler(svc service.MatchService, eventSvc service.MatchEventService) *MatchHandler {
	return &MatchHandler{svc: svc, eventSvc: eventSvc}
}

// CreateMatch godoc
// @Summary      Tambah jadwal pertandingan
// @Tags         Matches
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateMatchRequest true "Payload jadwal"
// @Success      201 {object} response.Envelope{data=dto.MatchResponse}
// @Router       /matches [post]
func (h *MatchHandler) Create(c *gin.Context) {
	var req dto.CreateMatchRequest
	if !bindJSON(c, &req) {
		return
	}
	match, err := h.svc.Create(req)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.Created(c, "Jadwal pertandingan berhasil dibuat", match)
}

// ListMatches godoc
// @Summary      Daftar jadwal pertandingan
// @Tags         Matches
// @Security     BearerAuth
// @Produce      json
// @Param        page       query    int    false  "Nomor halaman (default 1)"
// @Param        limit      query    int    false  "Jumlah data per halaman (default 10, maksimum 100)"
// @Param        status     query    string false  "Filter status: scheduled, in_progress, finished, postpone"
// @Param        team_uuid  query    string false  "Filter berdasarkan UUID tim (home atau away)"
// @Param        date_from  query    string false  "Filter tanggal mulai (YYYY-MM-DD)"
// @Param        date_to    query    string false  "Filter tanggal akhir (YYYY-MM-DD)"
// @Param        sort_by    query    string false  "Kolom sorting: created_at, match_date, status"
// @Param        order      query    string false  "Arah sorting: asc atau desc (default desc)"
// @Success      200 {object} response.Envelope{data=[]dto.MatchResponse,meta=dto.PaginationMeta}
// @Router       /matches [get]
func (h *MatchHandler) List(c *gin.Context) {
	matches, meta, err := h.svc.List(c)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.Paginated(c, "Data jadwal berhasil diambil", matches, meta)
}

// DetailMatch godoc
// @Summary      Detail jadwal pertandingan
// @Tags         Matches
// @Security     BearerAuth
// @Produce      json
// @Param        uuid path string true "UUID Pertandingan"
// @Success      200 {object} response.Envelope{data=dto.MatchResponse}
// @Router       /matches/{uuid} [get]
func (h *MatchHandler) Detail(c *gin.Context) {
	id, ok := paramUUID(c, "uuid")
	if !ok {
		return
	}
	match, err := h.svc.Detail(id)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, "Detail pertandingan berhasil diambil", match)
}

// UpdateMatch godoc
// @Summary      Ubah jadwal pertandingan
// @Tags         Matches
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        uuid path string true "UUID Pertandingan"
// @Param        request body dto.UpdateMatchRequest true "Payload jadwal"
// @Success      200 {object} response.Envelope{data=dto.MatchResponse}
// @Router       /matches/{uuid} [put]
func (h *MatchHandler) Update(c *gin.Context) {
	id, ok := paramUUID(c, "uuid")
	if !ok {
		return
	}
	var req dto.UpdateMatchRequest
	if !bindJSON(c, &req) {
		return
	}
	match, err := h.svc.Update(id, req)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, "Jadwal pertandingan berhasil diubah", match)
}

// UpdateMatchStatus godoc
// @Summary      Ubah status pertandingan
// @Tags         Matches
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        uuid path string true "UUID Pertandingan"
// @Param        request body dto.UpdateMatchStatusRequest true "Status baru"
// @Success      200 {object} response.Envelope{data=dto.MatchResponse}
// @Router       /matches/{uuid}/status [patch]
func (h *MatchHandler) UpdateStatus(c *gin.Context) {
	id, ok := paramUUID(c, "uuid")
	if !ok {
		return
	}
	var req dto.UpdateMatchStatusRequest
	if !bindJSON(c, &req) {
		return
	}
	match, err := h.svc.UpdateStatus(id, req)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, "Status pertandingan berhasil diubah", match)
}

// DeleteMatch godoc
// @Summary      Soft delete jadwal pertandingan
// @Tags         Matches
// @Security     BearerAuth
// @Produce      json
// @Param        uuid path string true "UUID Pertandingan"
// @Success      200 {object} response.Envelope
// @Router       /matches/{uuid} [delete]
func (h *MatchHandler) Delete(c *gin.Context) {
	id, ok := paramUUID(c, "uuid")
	if !ok {
		return
	}
	if err := h.svc.Delete(id); err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, "Jadwal pertandingan berhasil dihapus", nil)
}

// SubmitResult godoc
// @Summary      Laporkan hasil pertandingan
// @Tags         Matches
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        uuid path string true "UUID Pertandingan"
// @Param        request body dto.MatchResultRequest true "Skor & event"
// @Success      201 {object} response.Envelope{data=dto.MatchResponse}
// @Router       /matches/{uuid}/result [post]
func (h *MatchHandler) SubmitResult(c *gin.Context) {
	id, ok := paramUUID(c, "uuid")
	if !ok {
		return
	}
	var req dto.MatchResultRequest
	if !bindJSON(c, &req) {
		return
	}
	match, err := h.eventSvc.SubmitResult(id, req, false)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.Created(c, "Hasil pertandingan berhasil disimpan", match)
}

// CorrectResult godoc
// @Summary      Koreksi hasil pertandingan
// @Tags         Matches
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        uuid path string true "UUID Pertandingan"
// @Param        request body dto.MatchResultRequest true "Skor & event"
// @Success      200 {object} response.Envelope{data=dto.MatchResponse}
// @Router       /matches/{uuid}/result [put]
func (h *MatchHandler) CorrectResult(c *gin.Context) {
	id, ok := paramUUID(c, "uuid")
	if !ok {
		return
	}
	var req dto.MatchResultRequest
	if !bindJSON(c, &req) {
		return
	}
	match, err := h.eventSvc.SubmitResult(id, req, true)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, "Hasil pertandingan berhasil dikoreksi", match)
}

// ListEvents godoc
// @Summary      Daftar event satu pertandingan
// @Tags         Matches
// @Security     BearerAuth
// @Produce      json
// @Param        uuid path string true "UUID Pertandingan"
// @Success      200 {object} response.Envelope{data=[]dto.MatchEventResponse}
// @Router       /matches/{uuid}/events [get]
func (h *MatchHandler) ListEvents(c *gin.Context) {
	id, ok := paramUUID(c, "uuid")
	if !ok {
		return
	}
	events, err := h.eventSvc.ListEvents(id)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, "Daftar event berhasil diambil", events)
}

// AddEvent godoc
// @Summary      Tambah satu event pertandingan
// @Tags         Matches
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        uuid path string true "UUID Pertandingan"
// @Param        request body dto.MatchEventRequest true "Payload event"
// @Success      201 {object} response.Envelope{data=dto.MatchEventResponse}
// @Router       /matches/{uuid}/events [post]
func (h *MatchHandler) AddEvent(c *gin.Context) {
	id, ok := paramUUID(c, "uuid")
	if !ok {
		return
	}
	var req dto.MatchEventRequest
	if !bindJSON(c, &req) {
		return
	}
	event, err := h.eventSvc.AddEvent(id, req)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.Created(c, "Event berhasil ditambahkan", event)
}

// DeleteEvent godoc
// @Summary      Soft delete satu event pertandingan
// @Tags         Matches
// @Security     BearerAuth
// @Produce      json
// @Param        uuid path string true "UUID Pertandingan"
// @Param        event_uuid path string true "UUID Event"
// @Success      200 {object} response.Envelope
// @Router       /matches/{uuid}/events/{event_uuid} [delete]
func (h *MatchHandler) DeleteEvent(c *gin.Context) {
	matchID, ok := paramUUID(c, "uuid")
	if !ok {
		return
	}
	eventID, ok := paramUUID(c, "event_uuid")
	if !ok {
		return
	}
	if err := h.eventSvc.DeleteEvent(matchID, eventID); err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, "Event berhasil dihapus", nil)
}
