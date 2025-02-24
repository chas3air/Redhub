package authcontroller

import (
	"apigateway/internal/domain/models"
	authservice "apigateway/internal/services/auth"
	"apigateway/pkg/lib/logger/sl"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type AuthController struct {
	log          *slog.Logger
	auth_service *authservice.AuthService
}

func New(log *slog.Logger, auth_service *authservice.AuthService) *AuthController {
	return &AuthController{
		log:          log,
		auth_service: auth_service,
	}
}

func (ac *AuthController) handleError(w http.ResponseWriter, err error, log *slog.Logger) {
	if errors.Is(err, context.Canceled) {
		log.Error("Request was canceled by the user.")
		http.Error(w, "Request canceled", http.StatusRequestTimeout)
	} else if errors.Is(err, context.DeadlineExceeded) {
		log.Error("Request timed out.")
		http.Error(w, "Request timeout", http.StatusRequestTimeout)
	} else {
		log.Error("Error occurred", sl.Err(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ac *AuthController) readRequestBody(r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	const op = "controllers.auth.login"
	log := ac.log.With(slog.String("op", op))

	body, err := ac.readRequestBody(r)
	if err != nil {
		log.Error("Bad request", sl.Err(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user_credentials struct {
		Login    string `json:"email"`
		Password string `json:"password"`
		AppID    string `json:"app_id"`
	}
	if err := json.Unmarshal(body, &user_credentials); err != nil {
		log.Error("Bad request", sl.Err(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	app_id, err := uuid.Parse(user_credentials.AppID)
	if err != nil {
		log.Error("Invalid UUID format", sl.Err(err))
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	token, err := ac.auth_service.Login(r.Context(), user_credentials.Login, user_credentials.Password, app_id)
	if err != nil {
		ac.handleError(w, err, log)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
	log.Info("Login succeeded")
}

func (ac *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	const op = "controllers.auth.register"
	log := ac.log.With(slog.String("op", op))

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Error("Bad request", sl.Err(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := ac.auth_service.Register(r.Context(), user); err != nil {
		ac.handleError(w, err, log)
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Info("Registration succeeded")
}
