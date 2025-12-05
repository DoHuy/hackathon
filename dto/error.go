package dto

type ErrorResponse struct {
	Message    string `json:"message" example:"Error description"`
	StatusCode int    `json:"status_code" example:"500"`
}