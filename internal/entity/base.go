package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID        uint64         `gorm:"primaryKey" json:"-"`
	UUID      uuid.UUID      `gorm:"type:uuid;uniqueIndex;default:gen_random_uuid()" json:"uuid"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) error {
	if b.UUID == uuid.Nil {
		b.UUID = uuid.New()
	}
	return nil
}
