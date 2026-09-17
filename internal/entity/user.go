package entity

type User struct {
	Base
	Email    string `gorm:"size:150;not null" json:"email"`
	Password string `gorm:"size:255;not null" json:"-"`
}

func (User) TableName() string {
	return "users"
}
