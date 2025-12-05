package handlers

import (
	"hackathon/config"
	"hackathon/dto"
	"hackathon/middleware"
	"hackathon/repositories"
	"hackathon/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type FileHandler struct {
	service  services.FileService
	userRepo repositories.UserRepository
	cfg      *config.Config
}

func NewFileHandler(g *echo.Group, s services.FileService, userRepo repositories.UserRepository, cfg *config.Config) *FileHandler {
	h := &FileHandler{service: s, userRepo: userRepo, cfg: cfg}
	uploadGroup := g.Group("/upload")
	uploadGroup.Use(middleware.NewAuthMiddleware(userRepo, cfg.JWT.Secret))
	uploadGroup.Use(middleware.BodySizeLimit(cfg.Storage.MaxSizeMB))
	uploadGroup.POST("", h.Upload)

	return h
}

// @Summary Upload file
// @Security BearerAuth
// @Tags file
// @Param data formData file true "Image"
// @Success 200 {object} dto.UploadResponse
// @Router /api/upload [post]
func (h *FileHandler) Upload(c echo.Context) error {
	file, err := c.FormFile("data")
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "File 'data' is required", StatusCode: http.StatusBadRequest})
	}
	metadata, err := h.service.UploadFile(file)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
	}
	return c.JSON(http.StatusOK, dto.UploadResponse{
		Message: "File uploaded successfully", Filename: metadata.Filename, ID: metadata.ID, ContentType: metadata.ContentType,
	})
}
