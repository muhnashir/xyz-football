package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xyz-corp/xyz-football-api/internal/dto"
	"github.com/xyz-corp/xyz-football-api/internal/service"
	"github.com/xyz-corp/xyz-football-api/pkg/apperror"
	"github.com/xyz-corp/xyz-football-api/pkg/response"
)

type TeamHandler struct {
	svc       service.TeamService
	playerSvc service.PlayerService
}

func NewTeamHandler(svc service.TeamService, playerSvc service.PlayerService) *TeamHandler {
	return &TeamHandler{svc: svc, playerSvc: playerSvc}
}

// CreateTeam godoc
// @Summary      Tambah tim baru
// @Tags         Teams
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateTeamRequest true "Payload tim"
// @Success      201 {object} response.Envelope{data=dto.TeamResponse}
// @Router       /teams [post]
func (h *TeamHandler) Create(c *gin.Context) {
	var req dto.CreateTeamRequest
	if !bindJSON(c, &req) {
		return
	}
	team, err := h.svc.Create(req)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.Created(c, "Tim berhasil dibuat", team)
}

// ListTeams godoc
// @Summary      Daftar tim
// @Tags         Teams
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} response.Envelope{data=[]dto.TeamResponse}
// @Router       /teams [get]
func (h *TeamHandler) List(c *gin.Context) {
	teams, meta, err := h.svc.List(c)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.Paginated(c, "Data tim berhasil diambil", teams, meta)
}

// DetailTeam godoc
// @Summary      Detail tim
// @Tags         Teams
// @Security     BearerAuth
// @Produce      json
// @Param        uuid path string true "UUID Tim"
// @Success      200 {object} response.Envelope{data=dto.TeamDetailResponse}
// @Router       /teams/{uuid} [get]
func (h *TeamHandler) Detail(c *gin.Context) {
	id, ok := paramUUID(c, "uuid")
	if !ok {
		return
	}
	team, err := h.svc.Detail(id)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, "Detail tim berhasil diambil", team)
}

// UpdateTeam godoc
// @Summary      Ubah tim
// @Tags         Teams
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        uuid path string true "UUID Tim"
// @Param        request body dto.UpdateTeamRequest true "Payload tim"
// @Success      200 {object} response.Envelope{data=dto.TeamResponse}
// @Router       /teams/{uuid} [put]
func (h *TeamHandler) Update(c *gin.Context) {
	id, ok := paramUUID(c, "uuid")
	if !ok {
		return
	}
	var req dto.UpdateTeamRequest
	if !bindJSON(c, &req) {
		return
	}
	team, err := h.svc.Update(id, req)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, "Tim berhasil diubah", team)
}

// DeleteTeam godoc
// @Summary      Soft delete tim
// @Tags         Teams
// @Security     BearerAuth
// @Produce      json
// @Param        uuid path string true "UUID Tim"
// @Success      200 {object} response.Envelope
// @Router       /teams/{uuid} [delete]
func (h *TeamHandler) Delete(c *gin.Context) {
	id, ok := paramUUID(c, "uuid")
	if !ok {
		return
	}
	if err := h.svc.Delete(id); err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, "Tim berhasil dihapus", nil)
}

// UploadLogo godoc
// @Summary      Upload logo tim
// @Tags         Teams
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        uuid path string true "UUID Tim"
// @Param        logo formData file true "File logo"
// @Success      200 {object} response.Envelope{data=dto.TeamResponse}
// @Router       /teams/{uuid}/logo [post]
func (h *TeamHandler) UploadLogo(c *gin.Context) {
	id, ok := paramUUID(c, "uuid")
	if !ok {
		return
	}
	file, err := c.FormFile("logo")
	if err != nil {
		handleErr(c, apperror.Validation("File logo wajib diisi"))
		return
	}
	team, err := h.svc.UploadLogo(id, file)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, "Logo berhasil diunggah", team)
}

// TeamPlayers godoc
// @Summary      Daftar pemain satu tim
// @Tags         Teams
// @Security     BearerAuth
// @Produce      json
// @Param        uuid path string true "UUID Tim"
// @Success      200 {object} response.Envelope{data=[]dto.PlayerResponse}
// @Router       /teams/{uuid}/players [get]
func (h *TeamHandler) Players(c *gin.Context) {
	id, ok := paramUUID(c, "uuid")
	if !ok {
		return
	}
	players, err := h.playerSvc.ListByTeam(id)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, "Daftar pemain berhasil diambil", players)
}
