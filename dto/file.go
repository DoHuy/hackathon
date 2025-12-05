package dto

type UploadResponse struct {
	Message     string `json:"message"`
	Filename    string `json:"filename"`
	ID          uint   `json:"id"`
	ContentType string `json:"content_type"`
}
