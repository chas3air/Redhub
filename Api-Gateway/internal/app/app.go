package app

import (
	authcontroller "apigateway/internal/controllers/auth"
	userscontroller "apigateway/internal/controllers/usersManager"
	authservice "apigateway/internal/services/auth"
	usersmanagerservice "apigateway/internal/services/usersManager"
	authstorage "apigateway/internal/storage/real/auth"
	usersmanagerstorage "apigateway/internal/storage/real/usersManager"
	"apigateway/pkg/config"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	log *slog.Logger
	cfg *config.Config
}

func New(log *slog.Logger, cfg *config.Config) *App {
	return &App{
		log: log,
		cfg: cfg,
	}
}

func (a *App) Start() {
	authstorage := authstorage.New(a.log, a.cfg.AuthHost, a.cfg.AuthPort)
	auth_service := authservice.New(a.log, authstorage)
	authcontroller := authcontroller.New(a.log, auth_service)

	usersmanagerstorage := usersmanagerstorage.New(a.log, a.cfg.UsersStorageHost, a.cfg.UsersStoragePort)
	usersmanager_service := usersmanagerservice.New(a.log, usersmanagerstorage)
	userscontroller := userscontroller.New(a.log, usersmanager_service)

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/health-check", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.HandleFunc("/api/v1/login", authcontroller.Login).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/register", authcontroller.Register).Methods(http.MethodPost)

	r.HandleFunc("/api/v1/users", userscontroller.GetUsers).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/users/{id}", userscontroller.GetUserById).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/users", userscontroller.Insert).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/users/{id}", userscontroller.Update).Methods(http.MethodPut)
	r.HandleFunc("/api/v1/users/{id}", userscontroller.Delete).Methods(http.MethodDelete)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", a.cfg.API.Port), r); err != nil {
		panic(err)
	}

}
