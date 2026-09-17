package dto

import "github.com/google/uuid"

type MatchEventRequest struct {
	PlayerUUID       uuid.UUID  `json:"player_uuid" binding:"required,uuid"`
	AssistPlayerUUID *uuid.UUID `json:"assist_player_uuid"`
	EventType        string     `json:"event_type" binding:"required,oneof=GOAL YELLOW_CARD RED_CARD OWN_GOAL"`
	Minute           int16      `json:"minute" binding:"required,gte=1,lte=130"`
}

type PlayerRef struct {
	UUID uuid.UUID `json:"uuid"`
	Name string    `json:"name"`
}

type MatchEventResponse struct {
	UUID         uuid.UUID  `json:"uuid"`
	EventType    string     `json:"event_type"`
	Minute       int16      `json:"minute"`
	Player       PlayerRef  `json:"player"`
	AssistPlayer *PlayerRef `json:"assist_player,omitempty"`
	Team         *TeamRef   `json:"team,omitempty"`
}
