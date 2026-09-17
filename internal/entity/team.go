package entity

type Team struct {
	Base
	Name        string  `gorm:"size:100;not null" json:"name"`
	FoundedYear int16   `gorm:"not null" json:"founded_year"`
	LogoURL     *string `gorm:"size:255" json:"logo_url"`
	Address     string  `gorm:"not null" json:"address"`
	City        string  `gorm:"size:100;not null" json:"city"`

	Players []Player `gorm:"foreignKey:TeamID" json:"-"`
}

func (Team) TableName() string {
	return "teams"
}
