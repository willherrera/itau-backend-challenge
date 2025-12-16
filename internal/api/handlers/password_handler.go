package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/willherrera/itau-backend-challenge/internal/api/models"
	"github.com/willherrera/itau-backend-challenge/internal/application"
	"github.com/willherrera/itau-backend-challenge/pkg/metrics"
)

type PasswordHandler struct {
	service *application.PasswordService
}

func NewPasswordHandler(service *application.PasswordService) *PasswordHandler {
	return &PasswordHandler{
		service: service,
	}
}

// ValidatePassword handles POST /api/v1/validate-password requests.
// @Summary Valida uma senha
// @Description Valida se uma senha atende a todos os critérios de segurança definidos
// @Tags Password
// @Accept json
// @Produce json
// @Param request body models.ValidatePasswordRequest true "Senha a ser validada"
// @Success 200 {object} models.ValidatePasswordResponse "Resultado da validação"
// @Failure 400 {object} models.ErrorResponse "Requisição inválida"
// @Failure 405 {object} models.ErrorResponse "Método não permitido"
// @Router /api/v1/validate-password [post]
func (h *PasswordHandler) ValidatePassword(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	metrics.InProgress.Inc()
	defer func() {
		metrics.InProgress.Dec()
		metrics.RequestDuration.Observe(time.Since(start).Seconds())
	}()

	var req models.ValidatePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Password == "" {
		h.sendError(w, http.StatusBadRequest, "Password field is required")
		return
	}

	result := h.service.Validate(req.Password)
	metrics.RecordValidation(result.IsValid, result.Errors)

	h.sendJSON(w, http.StatusOK, models.ValidatePasswordResponse{
		IsValid: result.IsValid,
		Errors:  result.Errors,
	})
}

// Health handles GET /health requests.
// @Summary Health check
// @Description Verifica se a API está funcionando corretamente
// @Tags Health
// @Produce json
// @Success 200 {object} models.HealthResponse "API está saudável"
// @Router /health [get]
func (h *PasswordHandler) Health(w http.ResponseWriter, r *http.Request) {
	h.sendJSON(w, http.StatusOK, models.HealthResponse{
		Status:  "healthy",
		Service: "password-validator",
	})
}

func (h *PasswordHandler) sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}

func (h *PasswordHandler) sendError(w http.ResponseWriter, status int, message string) {
	h.sendJSON(w, status, models.ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	})
}
