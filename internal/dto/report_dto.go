package dto

import "github.com/google/uuid"

type TopScorer struct {
	PlayerUUID uuid.UUID `json:"player_uuid"`
	PlayerName string    `json:"player_name"`
	TeamName   string    `json:"team_name"`
	Goals      int64     `json:"goals"`
}

type FinalScore struct {
	Home int16 `json:"home"`
	Away int16 `json:"away"`
}

type AccumulatedWins struct {
	HomeTeam int64 `json:"home_team"`
	AwayTeam int64 `json:"away_team"`
}

type MatchReportResponse struct {
	UUID              uuid.UUID       `json:"uuid"`
	MatchDate         string          `json:"match_date"`
	MatchTime         string          `json:"match_time"`
	Status            string          `json:"status"`
	HomeTeam          TeamRef         `json:"home_team"`
	AwayTeam          TeamRef         `json:"away_team"`
	FinalScore        FinalScore      `json:"final_score"`
	ResultStatus      string          `json:"result_status"`
	ResultDescription string          `json:"result_description"`
	TopScorers        []TopScorer     `json:"top_scorers"`
	AccumulatedWins   AccumulatedWins `json:"accumulated_wins"`
}
