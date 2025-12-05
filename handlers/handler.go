package handlers

import (
	"hackathon/config"
	"hackathon/repositories"
	"hackathon/services"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	group    *echo.Group
	services *services.Service
	cfg      *config.Config
	repos    *repositories.Repository
}

func NewHandler(g *echo.Group, srv *services.Service, cfg *config.Config, repos *repositories.Repository) *Handler {
	return &Handler{group: g, services: srv, cfg: cfg, repos: repos}
}

func (h *Handler) RegisterRoutes() {
	NewAuthHandler(h.group, h.services.Auth, h.repos.User, h.cfg)
	NewFileHandler(h.group, *h.services.File, h.repos.User, h.cfg)
}
