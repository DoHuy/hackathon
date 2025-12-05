package models

import "time"

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"-"`
}

type FileMetadata struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Filename    string    `gorm:"type:text" json:"filename"`
	Size        int64     `gorm:"type:integer" json:"size"`
	ContentType string    `gorm:"type:text" json:"content_type"`
	UploadedAt  time.Time `gorm:"autoCreateTime" json:"uploaded_at"`
}

func (FileMetadata) TableName() string { return "file_metadata" }
