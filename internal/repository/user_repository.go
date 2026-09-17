package repository

import (
	"github.com/google/uuid"
	"github.com/xyz-corp/xyz-football-api/internal/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
	FindByUUID(id uuid.UUID) (*entity.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("LOWER(email) = LOWER(?)", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUUID(id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("uuid = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
