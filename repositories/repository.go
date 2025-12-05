package repositories

import "gorm.io/gorm"

type Repository struct {
	User UserRepository
	File FileRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User: NewUserRepository(db),
		File: NewFileRepository(db),
	}
}
