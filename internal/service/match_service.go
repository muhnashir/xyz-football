package service

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xyz-corp/xyz-football-api/internal/dto"
	"github.com/xyz-corp/xyz-football-api/internal/entity"
	"github.com/xyz-corp/xyz-football-api/internal/repository"
	"github.com/xyz-corp/xyz-football-api/pkg/apperror"
	"github.com/xyz-corp/xyz-football-api/pkg/pagination"
)

var matchAllowedSort = []string{"created_at", "match_date", "status"}

const dateLayout = "2006-01-02"

type MatchService interface {
	Create(req dto.CreateMatchRequest) (*dto.MatchResponse, error)
	Update(id uuid.UUID, req dto.UpdateMatchRequest) (*dto.MatchResponse, error)
	UpdateStatus(id uuid.UUID, req dto.UpdateMatchStatusRequest) (*dto.MatchResponse, error)
	Delete(id uuid.UUID) error
	Detail(id uuid.UUID) (*dto.MatchResponse, error)
	List(c *gin.Context) ([]dto.MatchResponse, pagination.Meta, error)
}

type matchService struct {
	matchRepo repository.MatchRepository
	teamRepo  repository.TeamRepository
}

func NewMatchService(matchRepo repository.MatchRepository, teamRepo repository.TeamRepository) MatchService {
	return &matchService{matchRepo: matchRepo, teamRepo: teamRepo}
}

func (s *matchService) resolveTeams(homeUUID, awayUUID uuid.UUID) (uint64, uint64, error) {
	if homeUUID == awayUUID {
		return 0, 0, apperror.BusinessRule("Tim tuan rumah dan tim tamu tidak boleh sama")
	}
	homeID, err := s.teamRepo.FindIDByUUID(homeUUID)
	if err != nil {
		return 0, 0, wrapNotFound(err, "Tim tuan rumah tidak ditemukan")
	}
	awayID, err := s.teamRepo.FindIDByUUID(awayUUID)
	if err != nil {
		return 0, 0, wrapNotFound(err, "Tim tamu tidak ditemukan")
	}
	return homeID, awayID, nil
}

func (s *matchService) Create(req dto.CreateMatchRequest) (*dto.MatchResponse, error) {
	homeID, awayID, err := s.resolveTeams(req.HomeTeamUUID, req.AwayTeamUUID)
	if err != nil {
		return nil, err
	}

	matchDate, _ := time.Parse(dateLayout, req.MatchDate)

	if conflict, err := s.matchRepo.ExistsScheduleConflict(homeID, matchDate, req.MatchTime, 0); err != nil {
		return nil, apperror.Internal(err.Error())
	} else if conflict {
		return nil, apperror.Conflict("Tim tuan rumah sudah memiliki jadwal pada tanggal & jam yang sama")
	}
	if conflict, err := s.matchRepo.ExistsScheduleConflict(awayID, matchDate, req.MatchTime, 0); err != nil {
		return nil, apperror.Internal(err.Error())
	} else if conflict {
		return nil, apperror.Conflict("Tim tamu sudah memiliki jadwal pada tanggal & jam yang sama")
	}

	match := &entity.Match{
		MatchDate:  matchDate,
		MatchTime:  req.MatchTime,
		HomeTeamID: homeID,
		AwayTeamID: awayID,
		Status:     entity.MatchStatusScheduled,
	}
	if err := s.matchRepo.Create(match); err != nil {
		return nil, apperror.Internal(err.Error())
	}

	created, err := s.matchRepo.FindByUUID(match.UUID)
	if err != nil {
		return nil, apperror.Internal(err.Error())
	}
	return toMatchResponse(created), nil
}

func (s *matchService) Update(id uuid.UUID, req dto.UpdateMatchRequest) (*dto.MatchResponse, error) {
	match, err := s.matchRepo.FindByUUID(id)
	if err != nil {
		return nil, wrapNotFound(err, "Jadwal pertandingan tidak ditemukan")
	}
	if match.Status == entity.MatchStatusFinished {
		return nil, apperror.BusinessRule("Jadwal pertandingan yang sudah finished tidak dapat diubah")
	}

	homeID, awayID, err := s.resolveTeams(req.HomeTeamUUID, req.AwayTeamUUID)
	if err != nil {
		return nil, err
	}

	matchDate, _ := time.Parse(dateLayout, req.MatchDate)

	if conflict, err := s.matchRepo.ExistsScheduleConflict(homeID, matchDate, req.MatchTime, match.ID); err != nil {
		return nil, apperror.Internal(err.Error())
	} else if conflict {
		return nil, apperror.Conflict("Tim tuan rumah sudah memiliki jadwal pada tanggal & jam yang sama")
	}
	if conflict, err := s.matchRepo.ExistsScheduleConflict(awayID, matchDate, req.MatchTime, match.ID); err != nil {
		return nil, apperror.Internal(err.Error())
	} else if conflict {
		return nil, apperror.Conflict("Tim tamu sudah memiliki jadwal pada tanggal & jam yang sama")
	}

	match.MatchDate = matchDate
	match.MatchTime = req.MatchTime
	match.HomeTeamID = homeID
	match.AwayTeamID = awayID

	if err := s.matchRepo.Update(match); err != nil {
		return nil, apperror.Internal(err.Error())
	}

	updated, err := s.matchRepo.FindByUUID(match.UUID)
	if err != nil {
		return nil, apperror.Internal(err.Error())
	}
	return toMatchResponse(updated), nil
}

func (s *matchService) UpdateStatus(id uuid.UUID, req dto.UpdateMatchStatusRequest) (*dto.MatchResponse, error) {
	match, err := s.matchRepo.FindByUUID(id)
	if err != nil {
		return nil, wrapNotFound(err, "Jadwal pertandingan tidak ditemukan")
	}

	if req.Status == entity.MatchStatusFinished {
		return nil, apperror.BusinessRule("Status finished hanya dapat diset melalui endpoint pelaporan hasil (/matches/{uuid}/result)")
	}

	allowed := entity.AllowedMatchTransitions[match.Status]
	valid := false
	for _, s := range allowed {
		if s == req.Status {
			valid = true
			break
		}
	}
	if !valid {
		return nil, apperror.BusinessRule("Transisi status dari '" + match.Status + "' ke '" + req.Status + "' tidak diizinkan")
	}

	match.Status = req.Status
	if err := s.matchRepo.Update(match); err != nil {
		return nil, apperror.Internal(err.Error())
	}

	return toMatchResponse(match), nil
}

func (s *matchService) Delete(id uuid.UUID) error {
	match, err := s.matchRepo.FindByUUID(id)
	if err != nil {
		return wrapNotFound(err, "Jadwal pertandingan tidak ditemukan")
	}
	return s.matchRepo.Delete(match)
}

func (s *matchService) Detail(id uuid.UUID) (*dto.MatchResponse, error) {
	match, err := s.matchRepo.FindByUUIDWithEvents(id)
	if err != nil {
		return nil, wrapNotFound(err, "Jadwal pertandingan tidak ditemukan")
	}
	return toMatchResponseWithEvents(match), nil
}

func (s *matchService) List(c *gin.Context) ([]dto.MatchResponse, pagination.Meta, error) {
	p := pagination.Parse(c, matchAllowedSort, "match_date")

	filter := repository.MatchListFilter{
		Status:   c.Query("status"),
		DateFrom: c.Query("date_from"),
		DateTo:   c.Query("date_to"),
	}
	if teamUUIDStr := c.Query("team_uuid"); teamUUIDStr != "" {
		teamUUID, err := parseUUID(teamUUIDStr)
		if err != nil {
			return nil, pagination.Meta{}, apperror.Validation("team_uuid tidak valid")
		}
		teamID, err := s.teamRepo.FindIDByUUID(teamUUID)
		if err != nil {
			return nil, pagination.Meta{}, wrapNotFound(err, "Tim tidak ditemukan")
		}
		filter.TeamID = teamID
	}

	matches, total, err := s.matchRepo.List(p, filter)
	if err != nil {
		return nil, pagination.Meta{}, apperror.Internal(err.Error())
	}

	result := make([]dto.MatchResponse, 0, len(matches))
	for i := range matches {
		result = append(result, *toMatchResponse(&matches[i]))
	}
	return result, pagination.BuildMeta(p.Page, p.Limit, total), nil
}

func toMatchResponse(m *entity.Match) *dto.MatchResponse {
	return &dto.MatchResponse{
		UUID:      m.UUID,
		MatchDate: m.MatchDate.Format(dateLayout),
		MatchTime: m.MatchTime,
		Status:    m.Status,
		HomeScore: m.HomeScore,
		AwayScore: m.AwayScore,
		HomeTeam:  dto.TeamRef{UUID: m.HomeTeam.UUID, Name: m.HomeTeam.Name},
		AwayTeam:  dto.TeamRef{UUID: m.AwayTeam.UUID, Name: m.AwayTeam.Name},
	}
}

func toMatchResponseWithEvents(m *entity.Match) *dto.MatchResponse {
	resp := toMatchResponse(m)
	resp.Events = make([]dto.MatchEventResponse, 0, len(m.Events))
	for _, e := range m.Events {
		resp.Events = append(resp.Events, toMatchEventResponse(&e, m))
	}
	return resp
}

func toMatchEventResponse(e *entity.MatchEvent, m *entity.Match) dto.MatchEventResponse {
	resp := dto.MatchEventResponse{
		UUID:      e.UUID,
		EventType: e.EventType,
		Minute:    e.Minute,
		Player:    dto.PlayerRef{UUID: e.Player.UUID, Name: e.Player.Name},
	}
	if e.AssistPlayer != nil {
		resp.AssistPlayer = &dto.PlayerRef{UUID: e.AssistPlayer.UUID, Name: e.AssistPlayer.Name}
	}
	// Team pelaku diturunkan dari players.team_id (lihat TRD Bagian 13, risiko perpindahan pemain).
	if e.Player.TeamID == m.HomeTeamID {
		resp.Team = &dto.TeamRef{UUID: m.HomeTeam.UUID, Name: m.HomeTeam.Name}
	} else if e.Player.TeamID == m.AwayTeamID {
		resp.Team = &dto.TeamRef{UUID: m.AwayTeam.UUID, Name: m.AwayTeam.Name}
	}
	return resp
}
