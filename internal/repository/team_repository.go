package repository

import (
	"github.com/google/uuid"
	"github.com/xyz-corp/xyz-football-api/internal/entity"
	"github.com/xyz-corp/xyz-football-api/pkg/pagination"
	"gorm.io/gorm"
)

type TeamRepository interface {
	Create(team *entity.Team) error
	Update(team *entity.Team) error
	Delete(team *entity.Team) error
	FindByUUID(id uuid.UUID) (*entity.Team, error)
	FindByUUIDWithPlayers(id uuid.UUID) (*entity.Team, error)
	FindIDByUUID(id uuid.UUID) (uint64, error)
	ExistsByNameActive(name string, excludeID uint64) (bool, error)
	List(p pagination.Params, city string) ([]entity.Team, int64, error)
	CountPlayers(teamID uint64) (int64, error)
	CountActiveMatches(teamID uint64) (int64, error)
}

type teamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &teamRepository{db: db}
}

func (r *teamRepository) Create(team *entity.Team) error {
	return r.db.Create(team).Error
}

func (r *teamRepository) Update(team *entity.Team) error {
	return r.db.Save(team).Error
}

func (r *teamRepository) Delete(team *entity.Team) error {
	return r.db.Delete(team).Error
}

func (r *teamRepository) FindByUUID(id uuid.UUID) (*entity.Team, error) {
	var team entity.Team
	if err := r.db.Where("uuid = ?", id).First(&team).Error; err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *teamRepository) FindByUUIDWithPlayers(id uuid.UUID) (*entity.Team, error) {
	var team entity.Team
	if err := r.db.Preload("Players").Where("uuid = ?", id).First(&team).Error; err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *teamRepository) FindIDByUUID(id uuid.UUID) (uint64, error) {
	var team entity.Team
	if err := r.db.Select("id").Where("uuid = ?", id).First(&team).Error; err != nil {
		return 0, err
	}
	return team.ID, nil
}

func (r *teamRepository) ExistsByNameActive(name string, excludeID uint64) (bool, error) {
	var count int64
	q := r.db.Model(&entity.Team{}).Where("LOWER(name) = LOWER(?)", name)
	if excludeID > 0 {
		q = q.Where("id <> ?", excludeID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *teamRepository) List(p pagination.Params, city string) ([]entity.Team, int64, error) {
	var teams []entity.Team
	var total int64

	q := r.db.Model(&entity.Team{})
	if p.Search != "" {
		q = q.Where("LOWER(name) LIKE ?", "%"+p.Search+"%")
	}
	if city != "" {
		q = q.Where("LOWER(city) = LOWER(?)", city)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := q.Order(p.SortBy + " " + p.Order).
		Offset(p.Offset).Limit(p.Limit).
		Find(&teams).Error
	return teams, total, err
}

func (r *teamRepository) CountPlayers(teamID uint64) (int64, error) {
	var count int64
	err := r.db.Model(&entity.Player{}).Where("team_id = ?", teamID).Count(&count).Error
	return count, err
}

func (r *teamRepository) CountActiveMatches(teamID uint64) (int64, error) {
	var count int64
	err := r.db.Model(&entity.Match{}).
		Where("(home_team_id = ? OR away_team_id = ?)", teamID, teamID).
		Where("status IN ?", []string{entity.MatchStatusScheduled, entity.MatchStatusInProgress}).
		Count(&count).Error
	return count, err
}
