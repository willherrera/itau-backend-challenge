package models

type ValidatePasswordRequest struct {
	Password string `json:"password" example:"AbTp9!fok" binding:"required"`
}

type ValidatePasswordResponse struct {
	IsValid bool     `json:"isValid" example:"true"`
	Errors  []string `json:"errors,omitempty" example:""`
}

type ErrorResponse struct {
	Error   string `json:"error" example:"Bad Request"`
	Message string `json:"message,omitempty" example:"Invalid request body"`
}

type HealthResponse struct {
	Status  string `json:"status" example:"healthy"`
	Service string `json:"service" example:"password-validator"`
}
