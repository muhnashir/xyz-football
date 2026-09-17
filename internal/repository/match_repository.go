package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/xyz-corp/xyz-football-api/internal/entity"
	"github.com/xyz-corp/xyz-football-api/pkg/pagination"
	"gorm.io/gorm"
)

type MatchListFilter struct {
	Status   string
	TeamID   uint64
	DateFrom string
	DateTo   string
}

type MatchRepository interface {
	DB() *gorm.DB
	Create(match *entity.Match) error
	Update(match *entity.Match) error
	Delete(match *entity.Match) error
	FindByUUID(id uuid.UUID) (*entity.Match, error)
	FindByUUIDWithEvents(id uuid.UUID) (*entity.Match, error)
	FindIDByUUID(id uuid.UUID) (uint64, error)
	ExistsScheduleConflict(teamID uint64, date time.Time, matchTime string, excludeID uint64) (bool, error)
	List(p pagination.Params, f MatchListFilter) ([]entity.Match, int64, error)
}

type matchRepository struct {
	db *gorm.DB
}

func NewMatchRepository(db *gorm.DB) MatchRepository {
	return &matchRepository{db: db}
}

func (r *matchRepository) DB() *gorm.DB {
	return r.db
}

func (r *matchRepository) Create(match *entity.Match) error {
	return r.db.Create(match).Error
}

func (r *matchRepository) Update(match *entity.Match) error {
	return r.db.Save(match).Error
}

func (r *matchRepository) Delete(match *entity.Match) error {
	return r.db.Delete(match).Error
}

func (r *matchRepository) FindByUUID(id uuid.UUID) (*entity.Match, error) {
	var match entity.Match
	if err := r.db.Preload("HomeTeam").Preload("AwayTeam").Where("uuid = ?", id).First(&match).Error; err != nil {
		return nil, err
	}
	return &match, nil
}

func (r *matchRepository) FindByUUIDWithEvents(id uuid.UUID) (*entity.Match, error) {
	var match entity.Match
	err := r.db.Preload("HomeTeam").Preload("AwayTeam").
		Preload("Events").Preload("Events.Player").Preload("Events.AssistPlayer").
		Where("uuid = ?", id).First(&match).Error
	if err != nil {
		return nil, err
	}
	return &match, nil
}

func (r *matchRepository) FindIDByUUID(id uuid.UUID) (uint64, error) {
	var match entity.Match
	if err := r.db.Select("id").Where("uuid = ?", id).First(&match).Error; err != nil {
		return 0, err
	}
	return match.ID, nil
}

func (r *matchRepository) ExistsScheduleConflict(teamID uint64, date time.Time, matchTime string, excludeID uint64) (bool, error) {
	var count int64
	q := r.db.Model(&entity.Match{}).
		Where("match_date = ? AND match_time = ?", date, matchTime).
		Where("(home_team_id = ? OR away_team_id = ?)", teamID, teamID)
	if excludeID > 0 {
		q = q.Where("id <> ?", excludeID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *matchRepository) List(p pagination.Params, f MatchListFilter) ([]entity.Match, int64, error) {
	var matches []entity.Match
	var total int64

	apply := func(q *gorm.DB) *gorm.DB {
		if f.Status != "" {
			q = q.Where("status = ?", f.Status)
		}
		if f.TeamID > 0 {
			q = q.Where("home_team_id = ? OR away_team_id = ?", f.TeamID, f.TeamID)
		}
		if f.DateFrom != "" {
			q = q.Where("match_date >= ?", f.DateFrom)
		}
		if f.DateTo != "" {
			q = q.Where("match_date <= ?", f.DateTo)
		}
		return q
	}

	countQ := apply(r.db.Model(&entity.Match{}))
	if err := countQ.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	listQ := apply(r.db.Model(&entity.Match{}).Preload("HomeTeam").Preload("AwayTeam"))
	err := listQ.Order(p.SortBy + " " + p.Order).
		Offset(p.Offset).Limit(p.Limit).
		Find(&matches).Error
	return matches, total, err
}
