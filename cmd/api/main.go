package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/willherrera/itau-backend-challenge/internal/api/handlers"
	"github.com/willherrera/itau-backend-challenge/internal/api/middleware"
	"github.com/willherrera/itau-backend-challenge/internal/application"
	"github.com/willherrera/itau-backend-challenge/internal/domain"
	"github.com/willherrera/itau-backend-challenge/internal/domain/rules"

	_ "github.com/willherrera/itau-backend-challenge/docs"
)

const (
	// Password validation constraints
	MinPasswordLength   = 9
	AllowedSpecialChars = "!@#$%^&*()-+"

	// Server configuration
	DefaultPort = "8080"
)

// @title Password Validator API
// @version 1.0
// @description API para validação de senhas com regras específicas de segurança
// @description
// @description Regras de validação:
// @description - Mínimo de 9 caracteres
// @description - Pelo menos 1 dígito
// @description - Pelo menos 1 letra minúscula
// @description - Pelo menos 1 letra maiúscula
// @description - Pelo menos 1 caractere especial (!@#$%^&*()-+)
// @description - Não deve conter caracteres repetidos

// @contact.url https://github.com/willherrera/itau-backend-challenge

// @host localhost:8080
// @schemes http

func main() {
	validators := []domain.PasswordValidator{
		rules.NewMinLengthValidator(MinPasswordLength),
		rules.NewDigitValidator(),
		rules.NewLowercaseValidator(),
		rules.NewUppercaseValidator(),
		rules.NewSpecialCharValidator(AllowedSpecialChars),
		rules.NewNoDuplicatesValidator(),
	}

	service := application.NewPasswordService(validators)
	handler := handlers.NewPasswordHandler(service)

	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/api/v1").Subrouter()
	apiRouter.HandleFunc("/validate-password", handler.ValidatePassword).Methods("POST", "OPTIONS")

	router.HandleFunc("/health", handler.Health).Methods("GET")
	router.Handle("/metrics", promhttp.Handler()).Methods("GET")

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.CORSMiddleware)

	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}

	addr := ":" + port
	log.Printf("Starting password validation API on %s", addr)
	log.Printf("Endpoints:")
	log.Printf("  POST   http://localhost%s/api/v1/validate-password", addr)
	log.Printf("  GET    http://localhost%s/health", addr)
	log.Printf("  GET    http://localhost%s/metrics", addr)
	log.Printf("  GET    http://localhost%s/swagger/index.html", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
