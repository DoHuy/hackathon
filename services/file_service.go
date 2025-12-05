package services

import (
	"errors"
	"hackathon/models"
	"hackathon/repositories"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

var (
	ErrFileTooLarge = errors.New("file too large")
	ErrInvalidType  = errors.New("invalid file type, only images allowed")
	ErrSaveFile     = errors.New("failed to save file to disk")
	ErrSaveDB       = errors.New("failed to save metadata")
)

type FileService struct {
	fileRepo  repositories.FileRepository
	uploadDir string
	maxSize   int64
}

func NewFileService(repo repositories.FileRepository, uploadDir string, maxSizeMB int64) *FileService {
	return &FileService{fileRepo: repo, uploadDir: uploadDir, maxSize: maxSizeMB * 1024 * 1024}
}

func (s *FileService) UploadFile(fileHeader *multipart.FileHeader) (*models.FileMetadata, error) {
	if fileHeader.Size > s.maxSize {
		return nil, ErrFileTooLarge
	}
	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	buffer := make([]byte, 512)
	if _, err := src.Read(buffer); err != nil {
		return nil, err
	}
	src.Seek(0, 0)

	contentType := http.DetectContentType(buffer)
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/gif" {
		return nil, ErrInvalidType
	}

	dstPath := filepath.Join(s.uploadDir, "upload-"+fileHeader.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return nil, ErrSaveFile
	}
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		return nil, ErrSaveFile
	}

	metadata := &models.FileMetadata{Filename: fileHeader.Filename, Size: fileHeader.Size, ContentType: contentType}
	if err := s.fileRepo.Create(metadata); err != nil {
		return nil, ErrSaveDB
	}
	return metadata, nil
}
