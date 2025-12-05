package repositories

import (
	"hackathon/models"

	"gorm.io/gorm"
)

type FileRepository interface {
	Create(metadata *models.FileMetadata) error
}

type fileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepository{db: db}
}

func (r *fileRepository) Create(metadata *models.FileMetadata) error {
	return r.db.Create(metadata).Error
}
