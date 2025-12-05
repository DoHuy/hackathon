package dto

type RegisterRequest struct {
	Username string `json:"username" example:"dev"`
	Password string `json:"password" example:"123456"`
}

type LoginRequest struct {
	Username string `json:"username" example:"dev"`
	Password string `json:"password" example:"123456"`
}

type TokenResponse struct {
	Token string `json:"token"`
	ExpiredTime int64 `json:"expired_time"`
}
