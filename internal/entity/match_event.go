package entity

const (
	EventTypeGoal       = "GOAL"
	EventTypeYellowCard = "YELLOW_CARD"
	EventTypeRedCard    = "RED_CARD"
	EventTypeOwnGoal    = "OWN_GOAL"
)

type MatchEvent struct {
	Base
	MatchID        uint64  `gorm:"not null" json:"-"`
	PlayerID       uint64  `gorm:"not null" json:"-"`
	AssistPlayerID *uint64 `json:"-"`
	EventType      string  `gorm:"size:20;not null" json:"event_type"`
	Minute         int16   `gorm:"not null" json:"minute"`

	Player       Player  `gorm:"foreignKey:PlayerID" json:"player,omitempty"`
	AssistPlayer *Player `gorm:"foreignKey:AssistPlayerID" json:"assist_player,omitempty"`
}

func (MatchEvent) TableName() string {
	return "match_events"
}
