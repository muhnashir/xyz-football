package entity

import "time"

const (
	MatchStatusScheduled  = "scheduled"
	MatchStatusInProgress = "in_progress"
	MatchStatusFinished   = "finished"
	MatchStatusPostpone   = "postpone"
)

// AllowedMatchTransitions defines the state machine from TRD 5.4.
var AllowedMatchTransitions = map[string][]string{
	MatchStatusScheduled:  {MatchStatusInProgress, MatchStatusPostpone},
	MatchStatusInProgress: {MatchStatusFinished},
	MatchStatusPostpone:   {MatchStatusScheduled},
	MatchStatusFinished:   {},
}

type Match struct {
	Base
	MatchDate  time.Time `gorm:"type:date;not null" json:"match_date"`
	MatchTime  string    `gorm:"type:time;not null" json:"match_time"`
	HomeTeamID uint64    `gorm:"not null" json:"-"`
	AwayTeamID uint64    `gorm:"not null" json:"-"`
	Status     string    `gorm:"size:20;not null;default:scheduled" json:"status"`
	HomeScore  int16     `gorm:"not null;default:0" json:"home_score"`
	AwayScore  int16     `gorm:"not null;default:0" json:"away_score"`

	HomeTeam Team         `gorm:"foreignKey:HomeTeamID" json:"home_team,omitempty"`
	AwayTeam Team         `gorm:"foreignKey:AwayTeamID" json:"away_team,omitempty"`
	Events   []MatchEvent `gorm:"foreignKey:MatchID" json:"events,omitempty"`
}

func (Match) TableName() string {
	return "matches"
}
