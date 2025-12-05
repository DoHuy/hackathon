package repositories

import (
	"hackathon/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByUsername(username string) (*models.User, error)
	FindByID(id uint) (*models.User, error)
	UpdateRevokeTokensBefore(user *models.User, timestamp int64) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateRevokeTokensBefore(user *models.User, timestamp int64) error {
	return r.db.Model(user).Update("revoke_tokens_before", timestamp).Error
}
