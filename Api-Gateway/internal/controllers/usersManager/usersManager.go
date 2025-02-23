package userscontroller

import (
	usersmanagerservice "apigateway/internal/services/usersManager"
	"apigateway/pkg/lib/logger/sl"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type UsersController struct {
	log          *slog.Logger
	usersService *usersmanagerservice.UsersManager
}

func New(log *slog.Logger, usersService *usersmanagerservice.UsersManager) *UsersController {
	return &UsersController{
		log:          log,
		usersService: usersService,
	}
}

// handleError обрабатывает ошибки и отправляет ответ клиенту
func (uc *UsersController) handleError(w http.ResponseWriter, err error, log *slog.Logger) {
	if errors.Is(err, context.Canceled) {
		log.Error("Request was canceled by the user.")
		http.Error(w, "Request canceled", http.StatusRequestTimeout)
	} else if errors.Is(err, context.DeadlineExceeded) {
		log.Error("Request timed out.")
		http.Error(w, "Request timeout", http.StatusRequestTimeout)
	} else {
		log.Error("Operation failed", sl.Err(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetUsers обрабатывает запрос на получение всех пользователей
func (uc *UsersController) GetUsers(w http.ResponseWriter, r *http.Request) {
	const op = "controllers.usersManager.getUsers"
	log := uc.log.With(slog.String("op", op))

	users, err := uc.usersService.GetUsers(r.Context())
	if err != nil {
		uc.handleError(w, err, log)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Error("Failed to encode response", sl.Err(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Info("Retrieved all users successfully")
}

// GetUserById обрабатывает запрос на получение пользователя по ID
func (uc *UsersController) GetUserById(w http.ResponseWriter, r *http.Request) {
	const op = "controllers.usersManager.getUserById"
	log := uc.log.With(slog.String("op", op))

	idStr := mux.Vars(r)["id"]
	uuidID, err := uuid.Parse(idStr)
	if err != nil {
		log.Error("ID must be a valid UUID", sl.Err(err))
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	user, err := uc.usersService.GetUserById(r.Context(), uuidID)
	if err != nil {
		uc.handleError(w, err, log)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Error("Failed to encode response", sl.Err(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Info("Retrieved user by ID successfully")
}

// GetUserByEmail обрабатывает запрос на получение пользователя по Email
func (uc *UsersController) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	const op = "controllers.usersManager.getUserByEmail"
	log := uc.log.With(slog.String("op", op))

	email := mux.Vars(r)["email"]
	user, err := uc.usersService.GetUserByEmail(r.Context(), email)
	if err != nil {
		uc.handleError(w, err, log)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Error("Failed to encode response", sl.Err(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Info("Retrieved user by email successfully")
}

func (uc *UsersController) Insert(w http.ResponseWriter, r *http.Request) {}
func (uc *UsersController) Update(w http.ResponseWriter, r *http.Request) {}
func (uc *UsersController) Delete(w http.ResponseWriter, r *http.Request) {}
