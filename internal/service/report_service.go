package service

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xyz-corp/xyz-football-api/config"
	"github.com/xyz-corp/xyz-football-api/internal/dto"
	"github.com/xyz-corp/xyz-football-api/internal/entity"
	"github.com/xyz-corp/xyz-football-api/internal/repository"
	"github.com/xyz-corp/xyz-football-api/pkg/apperror"
	"github.com/xyz-corp/xyz-football-api/pkg/pagination"
)

type ReportService interface {
	MatchReport(id uuid.UUID) (*dto.MatchReportResponse, error)
	MatchReportList(c *gin.Context) ([]dto.MatchReportResponse, pagination.Meta, error)
}

type reportService struct {
	matchRepo  repository.MatchRepository
	reportRepo repository.ReportRepository
	cfg        *config.Config
}

func NewReportService(matchRepo repository.MatchRepository, reportRepo repository.ReportRepository, cfg *config.Config) ReportService {
	return &reportService{matchRepo: matchRepo, reportRepo: reportRepo, cfg: cfg}
}

func (s *reportService) buildReport(m *entity.Match) (*dto.MatchReportResponse, error) {
	topScorers, err := s.reportRepo.TopScorers(m.ID)
	if err != nil {
		return nil, apperror.Internal(err.Error())
	}
	if topScorers == nil {
		topScorers = []dto.TopScorer{}
	}

	homeWins, err := s.reportRepo.AccumulatedWins(m.HomeTeamID, s.cfg.SeasonStartDate, m.MatchDate, m.MatchTime, m.ID)
	if err != nil {
		return nil, apperror.Internal(err.Error())
	}
	awayWins, err := s.reportRepo.AccumulatedWins(m.AwayTeamID, s.cfg.SeasonStartDate, m.MatchDate, m.MatchTime, m.ID)
	if err != nil {
		return nil, apperror.Internal(err.Error())
	}

	resultStatus, resultDesc := deriveResult(m.HomeScore, m.AwayScore)

	return &dto.MatchReportResponse{
		UUID:              m.UUID,
		MatchDate:         m.MatchDate.Format(dateLayout),
		MatchTime:         formatMatchTime(m.MatchTime),
		Status:            m.Status,
		HomeTeam:          dto.TeamRef{UUID: m.HomeTeam.UUID, Name: m.HomeTeam.Name},
		AwayTeam:          dto.TeamRef{UUID: m.AwayTeam.UUID, Name: m.AwayTeam.Name},
		FinalScore:        dto.FinalScore{Home: m.HomeScore, Away: m.AwayScore},
		ResultStatus:      resultStatus,
		ResultDescription: resultDesc,
		TopScorers:        topScorers,
		AccumulatedWins:   dto.AccumulatedWins{HomeTeam: homeWins, AwayTeam: awayWins},
	}, nil
}

func deriveResult(home, away int16) (string, string) {
	switch {
	case home > away:
		return "HOME_WIN", "Tim Home Menang"
	case home < away:
		return "AWAY_WIN", "Tim Away Menang"
	default:
		return "DRAW", "Draw"
	}
}

func (s *reportService) MatchReport(id uuid.UUID) (*dto.MatchReportResponse, error) {
	match, err := s.matchRepo.FindByUUID(id)
	if err != nil {
		return nil, wrapNotFound(err, "Pertandingan tidak ditemukan")
	}
	if match.Status != entity.MatchStatusFinished {
		return nil, apperror.BusinessRule("Report hanya tersedia untuk pertandingan berstatus finished")
	}
	return s.buildReport(match)
}

var reportAllowedSort = []string{"match_date", "created_at"}

func (s *reportService) MatchReportList(c *gin.Context) ([]dto.MatchReportResponse, pagination.Meta, error) {
	p := pagination.Parse(c, reportAllowedSort, "match_date")

	filter := repository.MatchListFilter{Status: entity.MatchStatusFinished}
	matches, total, err := s.matchRepo.List(p, filter)
	if err != nil {
		return nil, pagination.Meta{}, apperror.Internal(err.Error())
	}

	result := make([]dto.MatchReportResponse, 0, len(matches))
	for i := range matches {
		report, err := s.buildReport(&matches[i])
		if err != nil {
			return nil, pagination.Meta{}, err
		}
		result = append(result, *report)
	}

	return result, pagination.BuildMeta(p.Page, p.Limit, total), nil
}
