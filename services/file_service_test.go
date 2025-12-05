
package services

import (
	"bytes"
	"errors"
	"hackathon/models"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockFileRepository is a mock implementation of FileRepository for testing
type MockFileRepository struct {
	err error
}

func (m *MockFileRepository) Create(metadata *models.FileMetadata) error {
	if m.err != nil {
		return m.err
	}
	metadata.ID = 1
	return nil
}

func newTestFileService(repo *MockFileRepository, uploadDir string, maxSizeMB int64, allowedTypes []string) *FileService {
	return NewFileService(repo, uploadDir, maxSizeMB, allowedTypes)
}

func TestFileService_UploadFileStream(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "upload-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	t.Run("successful upload", func(t *testing.T) {
		repo := &MockFileRepository{}
		service := newTestFileService(repo, tmpDir, 1, []string{"text/plain; charset=utf-8"})
		content := "fake jpeg content"
		reader := strings.NewReader(content)
		filename := "test.jpg"

		metadata, err := service.UploadFileStream(reader, filename, int64(len(content)))
		assert.NoError(t, err)
		assert.NotNil(t, metadata)
		assert.Equal(t, filename, metadata.Filename)
		assert.Equal(t, int64(len(content)), metadata.Size)
		// This is tricky to test without a real image, http.DetectContentType will return "text/plain; charset=utf-8"
		// For a real jpeg it would be "image/jpeg"
		// We will assert that it's not empty
		assert.NotEmpty(t, metadata.ContentType)

		filePath := filepath.Join(tmpDir, "upload-"+filename)
		_, err = os.Stat(filePath)
		assert.NoError(t, err)
	})

	t.Run("file too large", func(t *testing.T) {
		repo := &MockFileRepository{}
		service := newTestFileService(repo, tmpDir, 1, []string{"image/jpeg"})
		content := "a very large content that exceeds the limit"
		reader := strings.NewReader(content)
		filename := "large.jpg"

		_, err := service.UploadFileStream(reader, filename, 2*1024*1024) // 2MB
		assert.Error(t, err)
		assert.Equal(t, ErrFileTooLarge, err)
	})

	t.Run("invalid file type", func(t *testing.T) {
		repo := &MockFileRepository{}
		service := newTestFileService(repo, tmpDir, 1, []string{"image/png"})
		// A simple string will be detected as text/plain
		content := "this is not a png"
		reader := strings.NewReader(content)
		filename := "test.txt"

		_, err := service.UploadFileStream(reader, filename, int64(len(content)))
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidType, err)
	})

	t.Run("database error", func(t *testing.T) {
		repo := &MockFileRepository{err: errors.New("database error")}
		service := newTestFileService(repo, tmpDir, 1, []string{"text/plain; charset=utf-8"})
		content := "fake jpeg content"
		reader := bytes.NewReader([]byte(content))
		filename := "test.jpg"

		_, err := service.UploadFileStream(reader, filename, int64(len(content)))
		assert.Error(t, err)
		assert.Equal(t, ErrSaveDB, err)

		filePath := filepath.Join(tmpDir, "upload-"+filename)
		_, err = os.Stat(filePath)
		assert.True(t, os.IsNotExist(err), "File should be deleted after db error")
	})
}
