package entity

const (
	PositionPenyerang     = "PENYERANG"
	PositionGelandang     = "GELANDANG"
	PositionBertahan      = "BERTAHAN"
	PositionPenjagaGawang = "PENJAGA_GAWANG"
)

type Player struct {
	Base
	TeamID       uint64 `gorm:"not null" json:"-"`
	Name         string `gorm:"size:100;not null" json:"name"`
	HeightCM     int16  `gorm:"column:height_cm;not null" json:"height_cm"`
	WeightKG     int16  `gorm:"column:weight_kg;not null" json:"weight_kg"`
	Position     string `gorm:"size:20;not null" json:"position"`
	JerseyNumber int16  `gorm:"not null" json:"jersey_number"`

	Team Team `gorm:"foreignKey:TeamID" json:"team,omitempty"`
}

func (Player) TableName() string {
	return "players"
}
