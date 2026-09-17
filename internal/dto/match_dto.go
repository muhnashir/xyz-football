package dto

import "github.com/google/uuid"

type CreateMatchRequest struct {
	MatchDate    string    `json:"match_date" binding:"required,datetime=2006-01-02"`
	MatchTime    string    `json:"match_time" binding:"required,datetime=15:04"`
	HomeTeamUUID uuid.UUID `json:"home_team_uuid" binding:"required,uuid"`
	AwayTeamUUID uuid.UUID `json:"away_team_uuid" binding:"required,uuid"`
}

type UpdateMatchRequest struct {
	MatchDate    string    `json:"match_date" binding:"required,datetime=2006-01-02"`
	MatchTime    string    `json:"match_time" binding:"required,datetime=15:04"`
	HomeTeamUUID uuid.UUID `json:"home_team_uuid" binding:"required,uuid"`
	AwayTeamUUID uuid.UUID `json:"away_team_uuid" binding:"required,uuid"`
}

type UpdateMatchStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=scheduled in_progress finished postpone"`
}

type MatchResultRequest struct {
	HomeScore int16               `json:"home_score" binding:"gte=0"`
	AwayScore int16               `json:"away_score" binding:"gte=0"`
	Events    []MatchEventRequest `json:"events"`
}

type MatchResponse struct {
	UUID      uuid.UUID            `json:"uuid"`
	MatchDate string               `json:"match_date"`
	MatchTime string               `json:"match_time"`
	Status    string               `json:"status"`
	HomeScore int16                `json:"home_score"`
	AwayScore int16                `json:"away_score"`
	HomeTeam  TeamRef              `json:"home_team"`
	AwayTeam  TeamRef              `json:"away_team"`
	Events    []MatchEventResponse `json:"events,omitempty"`
}
