package dto

type RegisterRequest struct {
	Username string `json:"username" validate:"required" example:"dev"`
	Password string `json:"password" validate:"required,min=6" example:"123456"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required" example:"dev"`
	Password string `json:"password" validate:"required,min=6" example:"123456"`
}

type TokenResponse struct {
	Token       string `json:"token"`
	ExpiredTime int64  `json:"expired_time"`
}
