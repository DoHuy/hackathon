package handlers

import (
	"hackathon/config"
	"hackathon/dto"
	"hackathon/middleware"
	"hackathon/models"
	"hackathon/repositories"
	"hackathon/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	service  *services.AuthService
	userRepo repositories.UserRepository
	cfg      *config.Config
}

func NewAuthHandler(g *echo.Group, s *services.AuthService, userRepo repositories.UserRepository, cfg *config.Config) *AuthHandler {
	h := &AuthHandler{service: s, userRepo: userRepo, cfg: cfg}
	authGroup := g.Group("/auth")
	authGroup.POST("/register", h.Register)
	authGroup.POST("/login", h.Login)

	authMiddleware := middleware.NewAuthMiddleware(userRepo, cfg.JWT.Secret)
	authGroup.POST("/revoke", h.Revoke, authMiddleware)
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
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
	}
	if err := h.service.Register(req.Username, req.Password); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error(), StatusCode: http.StatusInternalServerError})
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
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
	}
	tokenResponse, err := h.service.Login(req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Message: "Invalid credentials", StatusCode: http.StatusUnauthorized})
	}
	return c.JSON(http.StatusOK, tokenResponse)
}

// @Summary Revoke user token by time
// @Tags auth
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Router /api/auth/revoke [post]
func (h *AuthHandler) Revoke(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to get user from context", StatusCode: http.StatusInternalServerError})
	}

	if err := h.service.RevokeToken(user.ID); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error(), StatusCode: http.StatusInternalServerError})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "Tokens revoked"})
}
