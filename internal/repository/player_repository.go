package repository

import (
	"github.com/google/uuid"
	"github.com/xyz-corp/xyz-football-api/internal/entity"
	"github.com/xyz-corp/xyz-football-api/pkg/pagination"
	"gorm.io/gorm"
)

type PlayerRepository interface {
	Create(player *entity.Player) error
	Update(player *entity.Player) error
	Delete(player *entity.Player) error
	FindByUUID(id uuid.UUID) (*entity.Player, error)
	FindByID(id uint64) (*entity.Player, error)
	FindIDByUUID(id uuid.UUID) (uint64, error)
	ExistsByTeamJerseyActive(teamID uint64, jersey int16, excludeID uint64) (bool, error)
	ListByTeam(teamID uint64) ([]entity.Player, error)
	List(p pagination.Params, teamID uint64, position string) ([]entity.Player, int64, error)
}

type playerRepository struct {
	db *gorm.DB
}

func NewPlayerRepository(db *gorm.DB) PlayerRepository {
	return &playerRepository{db: db}
}

func (r *playerRepository) Create(player *entity.Player) error {
	return r.db.Create(player).Error
}

func (r *playerRepository) Update(player *entity.Player) error {
	return r.db.Save(player).Error
}

func (r *playerRepository) Delete(player *entity.Player) error {
	return r.db.Delete(player).Error
}

func (r *playerRepository) FindByUUID(id uuid.UUID) (*entity.Player, error) {
	var player entity.Player
	if err := r.db.Preload("Team").Where("uuid = ?", id).First(&player).Error; err != nil {
		return nil, err
	}
	return &player, nil
}

func (r *playerRepository) FindByID(id uint64) (*entity.Player, error) {
	var player entity.Player
	if err := r.db.Preload("Team").Where("id = ?", id).First(&player).Error; err != nil {
		return nil, err
	}
	return &player, nil
}

func (r *playerRepository) FindIDByUUID(id uuid.UUID) (uint64, error) {
	var player entity.Player
	if err := r.db.Select("id").Where("uuid = ?", id).First(&player).Error; err != nil {
		return 0, err
	}
	return player.ID, nil
}

func (r *playerRepository) ExistsByTeamJerseyActive(teamID uint64, jersey int16, excludeID uint64) (bool, error) {
	var count int64
	q := r.db.Model(&entity.Player{}).Where("team_id = ? AND jersey_number = ?", teamID, jersey)
	if excludeID > 0 {
		q = q.Where("id <> ?", excludeID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *playerRepository) ListByTeam(teamID uint64) ([]entity.Player, error) {
	var players []entity.Player
	err := r.db.Where("team_id = ?", teamID).Order("jersey_number asc").Find(&players).Error
	return players, err
}

func (r *playerRepository) List(p pagination.Params, teamID uint64, position string) ([]entity.Player, int64, error) {
	var players []entity.Player
	var total int64

	q := r.db.Model(&entity.Player{}).Preload("Team")
	if p.Search != "" {
		q = q.Where("LOWER(name) LIKE ?", "%"+p.Search+"%")
	}
	if teamID > 0 {
		q = q.Where("team_id = ?", teamID)
	}
	if position != "" {
		q = q.Where("position = ?", position)
	}

	countQ := r.db.Model(&entity.Player{})
	if p.Search != "" {
		countQ = countQ.Where("LOWER(name) LIKE ?", "%"+p.Search+"%")
	}
	if teamID > 0 {
		countQ = countQ.Where("team_id = ?", teamID)
	}
	if position != "" {
		countQ = countQ.Where("position = ?", position)
	}
	if err := countQ.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := q.Order(p.SortBy + " " + p.Order).
		Offset(p.Offset).Limit(p.Limit).
		Find(&players).Error
	return players, total, err
}
