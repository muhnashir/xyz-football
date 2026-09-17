package repository

import (
	"github.com/google/uuid"
	"github.com/xyz-corp/xyz-football-api/internal/entity"
	"gorm.io/gorm"
)

type MatchEventRepository interface {
	CreateMany(tx *gorm.DB, events []entity.MatchEvent) error
	SoftDeleteByMatchID(tx *gorm.DB, matchID uint64) error
	ListByMatch(matchID uint64) ([]entity.MatchEvent, error)
	CountByTypeAndMatch(matchID uint64, eventType string) (int64, error)
	FindByUUID(matchID uint64, eventUUID uuid.UUID) (*entity.MatchEvent, error)
	Delete(event *entity.MatchEvent) error
	ExistsDuplicate(matchID uint64, playerID uint64, eventType string, minute int16) (bool, error)
}

type matchEventRepository struct {
	db *gorm.DB
}

func NewMatchEventRepository(db *gorm.DB) MatchEventRepository {
	return &matchEventRepository{db: db}
}

func (r *matchEventRepository) CreateMany(tx *gorm.DB, events []entity.MatchEvent) error {
	if len(events) == 0 {
		return nil
	}
	return tx.Create(&events).Error
}

func (r *matchEventRepository) SoftDeleteByMatchID(tx *gorm.DB, matchID uint64) error {
	return tx.Where("match_id = ?", matchID).Delete(&entity.MatchEvent{}).Error
}

func (r *matchEventRepository) ListByMatch(matchID uint64) ([]entity.MatchEvent, error) {
	var events []entity.MatchEvent
	err := r.db.Preload("Player").Preload("AssistPlayer").
		Where("match_id = ?", matchID).Order("minute asc").Find(&events).Error
	return events, err
}

func (r *matchEventRepository) CountByTypeAndMatch(matchID uint64, eventType string) (int64, error) {
	var count int64
	err := r.db.Model(&entity.MatchEvent{}).
		Where("match_id = ? AND event_type = ?", matchID, eventType).
		Count(&count).Error
	return count, err
}

func (r *matchEventRepository) FindByUUID(matchID uint64, eventUUID uuid.UUID) (*entity.MatchEvent, error) {
	var event entity.MatchEvent
	err := r.db.Where("match_id = ? AND uuid = ?", matchID, eventUUID).First(&event).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *matchEventRepository) Delete(event *entity.MatchEvent) error {
	return r.db.Delete(event).Error
}

func (r *matchEventRepository) ExistsDuplicate(matchID uint64, playerID uint64, eventType string, minute int16) (bool, error) {
	var count int64
	err := r.db.Model(&entity.MatchEvent{}).
		Where("match_id = ? AND player_id = ? AND event_type = ? AND minute = ?", matchID, playerID, eventType, minute).
		Count(&count).Error
	return count > 0, err
}
