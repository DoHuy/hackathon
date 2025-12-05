package handlers

import (
	"hackathon/dto"
	"hackathon/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(g *echo.Group, s *services.AuthService) *AuthHandler {
	h := &AuthHandler{service: s}
	authGroup := g.Group("/auth")
	authGroup.POST("/register", h.Register)
	authGroup.POST("/login", h.Login)
	return h
}

// @Summary Register user
// @Tags auth
// @Param req body dto.RegisterRequest true "Info"
// @Success 201 {object} map[string]string
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	req := new(dto.RegisterRequest)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request")
	}
	if err := h.service.Register(req.Username, req.Password); err != nil {
		if err == services.ErrUserExists {
			return c.String(http.StatusConflict, err.Error())
		}
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, echo.Map{"message": "User registered"})
}

// @Summary Login
// @Tags auth
// @Param req body dto.LoginRequest true "Info"
// @Success 200 {object} dto.TokenResponse
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	req := new(dto.LoginRequest)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request")
	}
	tokenResponse, err := h.service.Login(req.Username, req.Password)
	if err != nil {
		return c.String(http.StatusUnauthorized, "Invalid credentials")
	}
	return c.JSON(http.StatusOK, tokenResponse)
}
