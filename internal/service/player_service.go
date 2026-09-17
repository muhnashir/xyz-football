package service

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xyz-corp/xyz-football-api/internal/dto"
	"github.com/xyz-corp/xyz-football-api/internal/entity"
	"github.com/xyz-corp/xyz-football-api/internal/repository"
	"github.com/xyz-corp/xyz-football-api/pkg/apperror"
	"github.com/xyz-corp/xyz-football-api/pkg/pagination"
)

var playerAllowedSort = []string{"created_at", "name", "jersey_number"}

type PlayerService interface {
	Create(req dto.CreatePlayerRequest) (*dto.PlayerResponse, error)
	Update(id uuid.UUID, req dto.UpdatePlayerRequest) (*dto.PlayerResponse, error)
	Delete(id uuid.UUID) error
	Detail(id uuid.UUID) (*dto.PlayerResponse, error)
	List(c *gin.Context) ([]dto.PlayerResponse, pagination.Meta, error)
	ListByTeam(teamUUID uuid.UUID) ([]dto.PlayerResponse, error)
}

type playerService struct {
	playerRepo repository.PlayerRepository
	teamRepo   repository.TeamRepository
}

func NewPlayerService(playerRepo repository.PlayerRepository, teamRepo repository.TeamRepository) PlayerService {
	return &playerService{playerRepo: playerRepo, teamRepo: teamRepo}
}

func (s *playerService) Create(req dto.CreatePlayerRequest) (*dto.PlayerResponse, error) {
	teamID, err := s.teamRepo.FindIDByUUID(req.TeamUUID)
	if err != nil {
		return nil, wrapNotFound(err, "Tim tidak ditemukan")
	}

	exists, err := s.playerRepo.ExistsByTeamJerseyActive(teamID, req.JerseyNumber, 0)
	if err != nil {
		return nil, apperror.Internal(err.Error())
	}
	if exists {
		return nil, apperror.Conflict("Nomor punggung sudah digunakan di tim ini",
			apperrorDetail("jersey_number", "nomor punggung harus unik dalam satu tim"))
	}

	player := &entity.Player{
		TeamID:       teamID,
		Name:         req.Name,
		HeightCM:     req.HeightCM,
		WeightKG:     req.WeightKG,
		Position:     req.Position,
		JerseyNumber: req.JerseyNumber,
	}
	if err := s.playerRepo.Create(player); err != nil {
		return nil, apperror.Internal(err.Error())
	}

	created, err := s.playerRepo.FindByID(player.ID)
	if err != nil {
		return nil, apperror.Internal(err.Error())
	}
	return toPlayerResponse(created), nil
}

func (s *playerService) Update(id uuid.UUID, req dto.UpdatePlayerRequest) (*dto.PlayerResponse, error) {
	player, err := s.playerRepo.FindByUUID(id)
	if err != nil {
		return nil, wrapNotFound(err, "Pemain tidak ditemukan")
	}

	teamID, err := s.teamRepo.FindIDByUUID(req.TeamUUID)
	if err != nil {
		return nil, wrapNotFound(err, "Tim tidak ditemukan")
	}

	// FR: saat team_uuid berubah, nomor punggung divalidasi ulang terhadap tim baru.
	exists, err := s.playerRepo.ExistsByTeamJerseyActive(teamID, req.JerseyNumber, player.ID)
	if err != nil {
		return nil, apperror.Internal(err.Error())
	}
	if exists {
		return nil, apperror.Conflict("Nomor punggung sudah digunakan di tim ini",
			apperrorDetail("jersey_number", "nomor punggung harus unik dalam satu tim"))
	}

	player.TeamID = teamID
	player.Name = req.Name
	player.HeightCM = req.HeightCM
	player.WeightKG = req.WeightKG
	player.Position = req.Position
	player.JerseyNumber = req.JerseyNumber

	if err := s.playerRepo.Update(player); err != nil {
		return nil, apperror.Internal(err.Error())
	}

	updated, err := s.playerRepo.FindByID(player.ID)
	if err != nil {
		return nil, apperror.Internal(err.Error())
	}
	return toPlayerResponse(updated), nil
}

func (s *playerService) Delete(id uuid.UUID) error {
	player, err := s.playerRepo.FindByUUID(id)
	if err != nil {
		return wrapNotFound(err, "Pemain tidak ditemukan")
	}
	return s.playerRepo.Delete(player)
}

func (s *playerService) Detail(id uuid.UUID) (*dto.PlayerResponse, error) {
	player, err := s.playerRepo.FindByUUID(id)
	if err != nil {
		return nil, wrapNotFound(err, "Pemain tidak ditemukan")
	}
	return toPlayerResponse(player), nil
}

func (s *playerService) List(c *gin.Context) ([]dto.PlayerResponse, pagination.Meta, error) {
	p := pagination.Parse(c, playerAllowedSort, "created_at")

	var teamID uint64
	if teamUUIDStr := c.Query("team_uuid"); teamUUIDStr != "" {
		teamUUID, err := parseUUID(teamUUIDStr)
		if err != nil {
			return nil, pagination.Meta{}, apperror.Validation("team_uuid tidak valid")
		}
		teamID, err = s.teamRepo.FindIDByUUID(teamUUID)
		if err != nil {
			return nil, pagination.Meta{}, wrapNotFound(err, "Tim tidak ditemukan")
		}
	}

	position := c.Query("position")

	players, total, err := s.playerRepo.List(p, teamID, position)
	if err != nil {
		return nil, pagination.Meta{}, apperror.Internal(err.Error())
	}

	result := make([]dto.PlayerResponse, 0, len(players))
	for i := range players {
		result = append(result, *toPlayerResponse(&players[i]))
	}
	return result, pagination.BuildMeta(p.Page, p.Limit, total), nil
}

func (s *playerService) ListByTeam(teamUUID uuid.UUID) ([]dto.PlayerResponse, error) {
	teamID, err := s.teamRepo.FindIDByUUID(teamUUID)
	if err != nil {
		return nil, wrapNotFound(err, "Tim tidak ditemukan")
	}

	players, err := s.playerRepo.ListByTeam(teamID)
	if err != nil {
		return nil, apperror.Internal(err.Error())
	}

	result := make([]dto.PlayerResponse, 0, len(players))
	for i := range players {
		result = append(result, *toPlayerResponse(&players[i]))
	}
	return result, nil
}

func toPlayerResponse(p *entity.Player) *dto.PlayerResponse {
	resp := &dto.PlayerResponse{
		UUID: p.UUID, Name: p.Name, HeightCM: p.HeightCM, WeightKG: p.WeightKG,
		Position: p.Position, JerseyNumber: p.JerseyNumber,
	}
	if p.Team.UUID != uuid.Nil {
		resp.Team = &dto.TeamRef{UUID: p.Team.UUID, Name: p.Team.Name}
	}
	return resp
}
