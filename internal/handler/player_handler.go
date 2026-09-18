package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xyz-corp/xyz-football-api/internal/dto"
	"github.com/xyz-corp/xyz-football-api/internal/service"
	"github.com/xyz-corp/xyz-football-api/pkg/response"
)

type PlayerHandler struct {
	svc service.PlayerService
}

func NewPlayerHandler(svc service.PlayerService) *PlayerHandler {
	return &PlayerHandler{svc: svc}
}

// CreatePlayer godoc
// @Summary      Tambah pemain
// @Tags         Players
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body dto.CreatePlayerRequest true "Payload pemain"
// @Success      201 {object} response.Envelope{data=dto.PlayerResponse}
// @Router       /players [post]
func (h *PlayerHandler) Create(c *gin.Context) {
	var req dto.CreatePlayerRequest
	if !bindJSON(c, &req) {
		return
	}
	player, err := h.svc.Create(req)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.Created(c, "Pemain berhasil dibuat", player)
}

// ListPlayers godoc
// @Summary      Daftar pemain
// @Tags         Players
// @Security     BearerAuth
// @Produce      json
// @Param        page       query    int    false  "Nomor halaman (default 1)"
// @Param        limit      query    int    false  "Jumlah data per halaman (default 10, maksimum 100)"
// @Param        search     query    string false  "Cari berdasarkan nama pemain"
// @Param        team_uuid  query    string false  "Filter berdasarkan UUID tim"
// @Param        position   query    string false  "Filter posisi: PENYERANG, GELANDANG, BERTAHAN, PENJAGA_GAWANG"
// @Param        sort_by    query    string false  "Kolom sorting: created_at, name, jersey_number"
// @Param        order      query    string false  "Arah sorting: asc atau desc (default desc)"
// @Success      200 {object} response.Envelope{data=[]dto.PlayerResponse,meta=dto.PaginationMeta}
// @Router       /players [get]
func (h *PlayerHandler) List(c *gin.Context) {
	players, meta, err := h.svc.List(c)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.Paginated(c, "Data pemain berhasil diambil", players, meta)
}

// DetailPlayer godoc
// @Summary      Detail pemain
// @Tags         Players
// @Security     BearerAuth
// @Produce      json
// @Param        uuid path string true "UUID Pemain"
// @Success      200 {object} response.Envelope{data=dto.PlayerResponse}
// @Router       /players/{uuid} [get]
func (h *PlayerHandler) Detail(c *gin.Context) {
	id, ok := paramUUID(c, "uuid")
	if !ok {
		return
	}
	player, err := h.svc.Detail(id)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, "Detail pemain berhasil diambil", player)
}

// UpdatePlayer godoc
// @Summary      Ubah pemain
// @Tags         Players
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        uuid path string true "UUID Pemain"
// @Param        request body dto.UpdatePlayerRequest true "Payload pemain"
// @Success      200 {object} response.Envelope{data=dto.PlayerResponse}
// @Router       /players/{uuid} [put]
func (h *PlayerHandler) Update(c *gin.Context) {
	id, ok := paramUUID(c, "uuid")
	if !ok {
		return
	}
	var req dto.UpdatePlayerRequest
	if !bindJSON(c, &req) {
		return
	}
	player, err := h.svc.Update(id, req)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, "Pemain berhasil diubah", player)
}

// DeletePlayer godoc
// @Summary      Soft delete pemain
// @Tags         Players
// @Security     BearerAuth
// @Produce      json
// @Param        uuid path string true "UUID Pemain"
// @Success      200 {object} response.Envelope
// @Router       /players/{uuid} [delete]
func (h *PlayerHandler) Delete(c *gin.Context) {
	id, ok := paramUUID(c, "uuid")
	if !ok {
		return
	}
	if err := h.svc.Delete(id); err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, "Pemain berhasil dihapus", nil)
}
