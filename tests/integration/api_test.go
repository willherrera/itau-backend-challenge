package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/willherrera/itau-backend-challenge/internal/api/handlers"
	"github.com/willherrera/itau-backend-challenge/internal/api/middleware"
	"github.com/willherrera/itau-backend-challenge/internal/api/models"
	"github.com/willherrera/itau-backend-challenge/internal/application"
	"github.com/willherrera/itau-backend-challenge/internal/domain"
	"github.com/willherrera/itau-backend-challenge/internal/domain/rules"
)

func setupTestServer() *httptest.Server {
	validators := []domain.PasswordValidator{
		rules.NewMinLengthValidator(9),
		rules.NewDigitValidator(),
		rules.NewLowercaseValidator(),
		rules.NewUppercaseValidator(),
		rules.NewSpecialCharValidator("!@#$%^&*()-+"),
		rules.NewNoDuplicatesValidator(),
	}

	service := application.NewPasswordService(validators)
	handler := handlers.NewPasswordHandler(service)

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/validate-password", handler.ValidatePassword).Methods("POST")
	router.HandleFunc("/health", handler.Health).Methods("GET")
	router.Use(middleware.LoggingMiddleware)

	return httptest.NewServer(router)
}

func TestValidatePasswordEndpoint(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	tests := []struct {
		name       string
		password   string
		wantValid  bool
		wantStatus int
	}{
		{
			name:       "valid password",
			password:   "AbTp9!fok",
			wantValid:  true,
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid empty password",
			password:   "",
			wantValid:  false,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid short password",
			password:   "aa",
			wantValid:  false,
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid password with duplicates",
			password:   "AbTp9!foo",
			wantValid:  false,
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid password with whitespace",
			password:   "AbTp9 fok",
			wantValid:  false,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := models.ValidatePasswordRequest{
				Password: tt.password,
			}
			body, _ := json.Marshal(reqBody)

			resp, err := http.Post(
				server.URL+"/api/v1/validate-password",
				"application/json",
				bytes.NewBuffer(body),
			)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("Status code = %d, want %d", resp.StatusCode, tt.wantStatus)
			}

			// Se o status é 400, não tentamos decodificar como ValidatePasswordResponse
			if resp.StatusCode == http.StatusBadRequest {
				return
			}

			var response models.ValidatePasswordResponse
			if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if response.IsValid != tt.wantValid {
				t.Errorf("IsValid = %v, want %v. Errors: %v",
					response.IsValid, tt.wantValid, response.Errors)
			}

			if !tt.wantValid && len(response.Errors) == 0 {
				t.Error("Expected errors for invalid password, got none")
			}
		})
	}
}

func TestAllRequirementExamples(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	examples := []struct {
		password string
		expected bool
	}{
		{"", false},
		{"aa", false},
		{"ab", false},
		{"AAAbbbCc", false},
		{"AbTp9!foo", false},
		{"AbTp9!foA", false},
		{"AbTp9 fok", false},
		{"AbTp9!fok", true},
	}

	for _, ex := range examples {
		t.Run("example: "+ex.password, func(t *testing.T) {
			reqBody := models.ValidatePasswordRequest{
				Password: ex.password,
			}
			body, _ := json.Marshal(reqBody)

			resp, err := http.Post(
				server.URL+"/api/v1/validate-password",
				"application/json",
				bytes.NewBuffer(body),
			)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			// Senha vazia agora retorna 400 Bad Request
			if ex.password == "" {
				if resp.StatusCode != http.StatusBadRequest {
					t.Errorf("Expected status 400 for empty password, got %d", resp.StatusCode)
				}
				return
			}

			var response models.ValidatePasswordResponse
			json.NewDecoder(resp.Body).Decode(&response)

			if response.IsValid != ex.expected {
				t.Errorf("IsValid(%q) = %v, want %v. Errors: %v",
					ex.password, response.IsValid, ex.expected, response.Errors)
			}
		})
	}
}

func TestHealthEndpoint(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var health models.HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if health.Status != "healthy" {
		t.Errorf("Health status = %s, want healthy", health.Status)
	}
}

func TestInvalidRequestBody(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	resp, err := http.Post(
		server.URL+"/api/v1/validate-password",
		"application/json",
		bytes.NewBufferString("invalid json"),
	)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Status code = %d, want %d", resp.StatusCode, http.StatusBadRequest)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/v1/validate-password")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Status code = %d, want %d", resp.StatusCode, http.StatusMethodNotAllowed)
	}
}
