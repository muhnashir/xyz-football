package service

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xyz-corp/xyz-football-api/internal/dto"
	"github.com/xyz-corp/xyz-football-api/internal/entity"
	"github.com/xyz-corp/xyz-football-api/internal/repository"
	"github.com/xyz-corp/xyz-football-api/pkg/apperror"
	"github.com/xyz-corp/xyz-football-api/pkg/pagination"
	"github.com/xyz-corp/xyz-football-api/pkg/uploader"
)

var teamAllowedSort = []string{"created_at", "name", "founded_year", "city"}

type TeamService interface {
	Create(req dto.CreateTeamRequest) (*dto.TeamResponse, error)
	Update(id uuid.UUID, req dto.UpdateTeamRequest) (*dto.TeamResponse, error)
	Delete(id uuid.UUID) error
	Detail(id uuid.UUID) (*dto.TeamDetailResponse, error)
	List(c *gin.Context) ([]dto.TeamResponse, pagination.Meta, error)
	UploadLogo(id uuid.UUID, file *multipart.FileHeader) (*dto.TeamResponse, error)
}

type teamService struct {
	teamRepo repository.TeamRepository
	up       *uploader.Uploader
}

func NewTeamService(teamRepo repository.TeamRepository, up *uploader.Uploader) TeamService {
	return &teamService{teamRepo: teamRepo, up: up}
}

func (s *teamService) Create(req dto.CreateTeamRequest) (*dto.TeamResponse, error) {
	exists, err := s.teamRepo.ExistsByNameActive(req.Name, 0)
	if err != nil {
		return nil, apperror.Internal(err.Error())
	}
	if exists {
		return nil, apperror.Conflict("Nama tim sudah digunakan", apperrorDetail("name", "nama tim harus unik"))
	}

	team := &entity.Team{
		Name:        req.Name,
		FoundedYear: req.FoundedYear,
		LogoURL:     req.LogoURL,
		Address:     req.Address,
		City:        req.City,
	}
	if err := s.teamRepo.Create(team); err != nil {
		return nil, apperror.Internal(err.Error())
	}

	return toTeamResponse(team, 0), nil
}

func (s *teamService) Update(id uuid.UUID, req dto.UpdateTeamRequest) (*dto.TeamResponse, error) {
	team, err := s.teamRepo.FindByUUID(id)
	if err != nil {
		return nil, wrapNotFound(err, "Tim tidak ditemukan")
	}

	exists, err := s.teamRepo.ExistsByNameActive(req.Name, team.ID)
	if err != nil {
		return nil, apperror.Internal(err.Error())
	}
	if exists {
		return nil, apperror.Conflict("Nama tim sudah digunakan", apperrorDetail("name", "nama tim harus unik"))
	}

	team.Name = req.Name
	team.FoundedYear = req.FoundedYear
	team.LogoURL = req.LogoURL
	team.Address = req.Address
	team.City = req.City

	if err := s.teamRepo.Update(team); err != nil {
		return nil, apperror.Internal(err.Error())
	}

	total, _ := s.teamRepo.CountPlayers(team.ID)
	return toTeamResponse(team, total), nil
}

func (s *teamService) Delete(id uuid.UUID) error {
	team, err := s.teamRepo.FindByUUID(id)
	if err != nil {
		return wrapNotFound(err, "Tim tidak ditemukan")
	}

	playerCount, err := s.teamRepo.CountPlayers(team.ID)
	if err != nil {
		return apperror.Internal(err.Error())
	}
	matchCount, err := s.teamRepo.CountActiveMatches(team.ID)
	if err != nil {
		return apperror.Internal(err.Error())
	}
	if playerCount > 0 || matchCount > 0 {
		return apperror.BusinessRule("Tim tidak dapat dihapus karena masih memiliki pemain aktif atau jadwal berstatus scheduled/in_progress")
	}

	return s.teamRepo.Delete(team)
}

func (s *teamService) Detail(id uuid.UUID) (*dto.TeamDetailResponse, error) {
	team, err := s.teamRepo.FindByUUIDWithPlayers(id)
	if err != nil {
		return nil, wrapNotFound(err, "Tim tidak ditemukan")
	}

	players := make([]dto.PlayerResponse, 0, len(team.Players))
	for _, p := range team.Players {
		players = append(players, dto.PlayerResponse{
			UUID: p.UUID, Name: p.Name, HeightCM: p.HeightCM, WeightKG: p.WeightKG,
			Position: p.Position, JerseyNumber: p.JerseyNumber,
		})
	}

	return &dto.TeamDetailResponse{
		TeamResponse: *toTeamResponse(team, int64(len(team.Players))),
		Players:      players,
	}, nil
}

func (s *teamService) List(c *gin.Context) ([]dto.TeamResponse, pagination.Meta, error) {
	p := pagination.Parse(c, teamAllowedSort, "created_at")
	city := c.Query("city")

	teams, total, err := s.teamRepo.List(p, city)
	if err != nil {
		return nil, pagination.Meta{}, apperror.Internal(err.Error())
	}

	result := make([]dto.TeamResponse, 0, len(teams))
	for _, t := range teams {
		count, _ := s.teamRepo.CountPlayers(t.ID)
		tt := t
		result = append(result, *toTeamResponse(&tt, count))
	}

	return result, pagination.BuildMeta(p.Page, p.Limit, total), nil
}

func (s *teamService) UploadLogo(id uuid.UUID, file *multipart.FileHeader) (*dto.TeamResponse, error) {
	team, err := s.teamRepo.FindByUUID(id)
	if err != nil {
		return nil, wrapNotFound(err, "Tim tidak ditemukan")
	}

	url, err := s.up.SaveLogo(file)
	if err != nil {
		return nil, apperror.Validation(err.Error())
	}

	team.LogoURL = &url
	if err := s.teamRepo.Update(team); err != nil {
		return nil, apperror.Internal(err.Error())
	}

	total, _ := s.teamRepo.CountPlayers(team.ID)
	return toTeamResponse(team, total), nil
}

func toTeamResponse(t *entity.Team, totalPlayers int64) *dto.TeamResponse {
	return &dto.TeamResponse{
		UUID:         t.UUID,
		Name:         t.Name,
		FoundedYear:  t.FoundedYear,
		LogoURL:      t.LogoURL,
		Address:      t.Address,
		City:         t.City,
		TotalPlayers: totalPlayers,
		CreatedAt:    t.CreatedAt,
	}
}
