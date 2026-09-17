package repository

import (
	"time"

	"github.com/xyz-corp/xyz-football-api/internal/dto"
	"gorm.io/gorm"
)

type ReportRepository interface {
	TopScorers(matchID uint64) ([]dto.TopScorer, error)
	AccumulatedWins(teamID uint64, seasonStart time.Time, matchDate time.Time, matchTime string, matchID uint64) (int64, error)
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}

// TopScorers returns all players tied for the most goals in a single match (TRD 8.8).
func (r *reportRepository) TopScorers(matchID uint64) ([]dto.TopScorer, error) {
	var scorers []dto.TopScorer
	err := r.db.Raw(`
		WITH goal_count AS (
			SELECT p.uuid AS player_uuid, p.name AS player_name, t.name AS team_name, COUNT(*) AS goals
			FROM match_events e
			JOIN players p ON p.id = e.player_id
			JOIN teams   t ON t.id = p.team_id
			WHERE e.match_id = ?
			  AND e.event_type = 'GOAL'
			  AND e.deleted_at IS NULL
			GROUP BY p.uuid, p.name, t.name
		)
		SELECT * FROM goal_count
		WHERE goals = (SELECT MAX(goals) FROM goal_count)
	`, matchID).Scan(&scorers).Error
	return scorers, err
}

// AccumulatedWins counts wins for a team from season start up to and including the reference match (TRD 8.8).
func (r *reportRepository) AccumulatedWins(teamID uint64, seasonStart time.Time, matchDate time.Time, matchTime string, matchID uint64) (int64, error) {
	var count int64
	err := r.db.Raw(`
		SELECT COUNT(*)
		FROM matches m
		WHERE m.deleted_at IS NULL
		  AND m.status = 'finished'
		  AND m.match_date >= ?
		  AND (m.match_date, m.match_time, m.id) <= (?, ?, ?)
		  AND (
		        (m.home_team_id = ? AND m.home_score > m.away_score)
		     OR (m.away_team_id = ? AND m.away_score > m.home_score)
		      )
	`, seasonStart, matchDate, matchTime, matchID, teamID, teamID).Scan(&count).Error
	return count, err
}
