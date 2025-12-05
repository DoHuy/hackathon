package handlers

import (
	"hackathon/config"
	"hackathon/services"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	group    *echo.Group
	services *services.Service
	cfg      *config.Config
}

func NewHandler(g *echo.Group, srv *services.Service, cfg *config.Config) *Handler {
	return &Handler{group: g, services: srv, cfg: cfg}
}

func (h *Handler) RegisterRoutes() {
	NewAuthHandler(h.group, h.services.Auth)
	NewFileHandler(h.group, h.services.File, h.cfg.JWT.Secret, h.cfg.Storage.MaxSizeMB)
}
