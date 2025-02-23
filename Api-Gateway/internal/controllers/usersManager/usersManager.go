package userscontroller

import (
	usersmanagerservice "apigateway/internal/services/usersManager"
	"log/slog"
	"net/http"
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

func (uc *UsersController) GetUsers(w http.ResponseWriter, r *http.Request)       {}
func (uc *UsersController) GetUserById(w http.ResponseWriter, r *http.Request)    {}
func (uc *UsersController) GetUserByEmail(w http.ResponseWriter, r *http.Request) {}
func (uc *UsersController) Insert(w http.ResponseWriter, r *http.Request)         {}
func (uc *UsersController) Update(w http.ResponseWriter, r *http.Request)         {}
func (uc *UsersController) Delete(w http.ResponseWriter, r *http.Request)         {}
