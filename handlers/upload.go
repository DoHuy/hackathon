package handlers

import (
	"hackathon/dto"
	"hackathon/middleware"
	"hackathon/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type FileHandler struct {
	service *services.FileService
}

func NewFileHandler(g *echo.Group, s *services.FileService, jwtSecret string, maxSizeMB int64) *FileHandler {
	h := &FileHandler{service: s}
	uploadGroup := g.Group("/upload")
	uploadGroup.Use(middleware.JWTMiddleware(jwtSecret))
	uploadGroup.Use(middleware.BodySizeLimit(maxSizeMB))
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
		return c.String(http.StatusBadRequest, "File 'data' is required")
	}
	metadata, err := h.service.UploadFile(file)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, dto.UploadResponse{
		Message: "File uploaded successfully", Filename: metadata.Filename, ID: metadata.ID, ContentType: metadata.ContentType,
	})
}
