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
	ErrInvalidType  = errors.New("invalid file type")
	ErrSaveFile     = errors.New("failed to save file to disk")
	ErrSaveDB       = errors.New("failed to save metadata")
)

type FileService struct {
	fileRepo     repositories.FileRepository
	uploadDir    string
	maxSize      int64
	allowedTypes []string
}

func NewFileService(repo repositories.FileRepository, uploadDir string, maxSizeMB int64, allowedTypes []string) *FileService {
	return &FileService{
		fileRepo:     repo,
		uploadDir:    uploadDir,
		maxSize:      maxSizeMB * 1024 * 1024,
		allowedTypes: allowedTypes,
	}
}

func (s *FileService) UploadFile(fileHeader *multipart.FileHeader) (*models.FileMetadata, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()
	return s.UploadFileStream(src, fileHeader.Filename, fileHeader.Size)
}

func (s *FileService) UploadFileStream(reader io.Reader, filename string, size int64) (*models.FileMetadata, error) {
	if size > s.maxSize {
		return nil, ErrFileTooLarge
	}

	buffer := make([]byte, 512)
	n, err := reader.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, err
	}
	buffer = buffer[:n]

	contentType := http.DetectContentType(buffer)

	isValidType := false
	if len(s.allowedTypes) > 0 {
		for _, t := range s.allowedTypes {
			if t == contentType {
				isValidType = true
				break
			}
		}
	} else {
		isValidType = true // No types specified, allow all
	}

	if !isValidType {
		return nil, ErrInvalidType
	}

	dstPath := filepath.Join(s.uploadDir, "upload-"+filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return nil, ErrSaveFile
	}
	defer dst.Close()

	if _, err := dst.Write(buffer); err != nil {
		return nil, ErrSaveFile
	}

	if _, err := io.Copy(dst, reader); err != nil {
		return nil, ErrSaveFile
	}

	metadata := &models.FileMetadata{Filename: filename, Size: size, ContentType: contentType}
	if err := s.fileRepo.Create(metadata); err != nil {
		os.Remove(dstPath) // Clean up file if DB save fails
		return nil, ErrSaveDB
	}

	return metadata, nil
}