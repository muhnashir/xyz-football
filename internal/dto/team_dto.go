package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateTeamRequest struct {
	Name        string  `json:"name" binding:"required,min=3,max=100"`
	FoundedYear int16   `json:"founded_year" binding:"required,gte=1800,lte=2100"`
	LogoURL     *string `json:"logo_url" binding:"omitempty,max=255" example:"teams/logo/3fa85f64-....jpg"`
	Address     string  `json:"address" binding:"required,max=500"`
	City        string  `json:"city" binding:"required,max=100"`
}

type UpdateTeamRequest struct {
	Name        string  `json:"name" binding:"required,min=3,max=100"`
	FoundedYear int16   `json:"founded_year" binding:"required,gte=1800,lte=2100"`
	LogoURL     *string `json:"logo_url" binding:"omitempty,max=255" example:"teams/logo/3fa85f64-....jpg"`
	Address     string  `json:"address" binding:"required,max=500"`
	City        string  `json:"city" binding:"required,max=100"`
}

type TeamResponse struct {
	UUID         uuid.UUID `json:"uuid"`
	Name         string    `json:"name"`
	FoundedYear  int16     `json:"founded_year"`
	LogoURL      *string   `json:"logo_url"`
	Address      string    `json:"address,omitempty"`
	City         string    `json:"city"`
	TotalPlayers int64     `json:"total_players"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
}

type TeamDetailResponse struct {
	TeamResponse
	Players []PlayerResponse `json:"players"`
}

type UploadImageResponse struct {
	Path string `json:"path" example:"teams/logo/3fa85f64-5717-4562-b3fc-2c963f66afa6.jpg"`
}
