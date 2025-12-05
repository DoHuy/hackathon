package services

import "hackathon/repositories"

type Service struct {
	Auth *AuthService
	File *FileService
}

func NewService(repos *repositories.Repository, jwtSecret []byte, expHours int, uploadDir string, maxSizeMB int64, allowedTypes []string) *Service {
	return &Service{
		Auth: NewAuthService(repos.User, jwtSecret, expHours),
		File: NewFileService(repos.File, uploadDir, maxSizeMB, allowedTypes),
	}
}
