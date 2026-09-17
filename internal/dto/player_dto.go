package dto

import "github.com/google/uuid"

type CreatePlayerRequest struct {
	TeamUUID     uuid.UUID `json:"team_uuid" binding:"required,uuid"`
	Name         string    `json:"name" binding:"required,min=3,max=100"`
	HeightCM     int16     `json:"height_cm" binding:"required,gte=100,lte=250"`
	WeightKG     int16     `json:"weight_kg" binding:"required,gte=30,lte=200"`
	Position     string    `json:"position" binding:"required,oneof=PENYERANG GELANDANG BERTAHAN PENJAGA_GAWANG"`
	JerseyNumber int16     `json:"jersey_number" binding:"required,gte=1,lte=99"`
}

type UpdatePlayerRequest struct {
	TeamUUID     uuid.UUID `json:"team_uuid" binding:"required,uuid"`
	Name         string    `json:"name" binding:"required,min=3,max=100"`
	HeightCM     int16     `json:"height_cm" binding:"required,gte=100,lte=250"`
	WeightKG     int16     `json:"weight_kg" binding:"required,gte=30,lte=200"`
	Position     string    `json:"position" binding:"required,oneof=PENYERANG GELANDANG BERTAHAN PENJAGA_GAWANG"`
	JerseyNumber int16     `json:"jersey_number" binding:"required,gte=1,lte=99"`
}

type TeamRef struct {
	UUID uuid.UUID `json:"uuid"`
	Name string    `json:"name"`
}

type PlayerResponse struct {
	UUID         uuid.UUID `json:"uuid"`
	Name         string    `json:"name"`
	HeightCM     int16     `json:"height_cm"`
	WeightKG     int16     `json:"weight_kg"`
	Position     string    `json:"position"`
	JerseyNumber int16     `json:"jersey_number"`
	Team         *TeamRef  `json:"team,omitempty"`
}
