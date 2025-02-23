package authcontroller

import (
	authservice "apigateway/internal/services/auth"
	"log/slog"
	"net/http"
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

func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {

}
func (ac *AuthController) Register(w http.ResponseWriter, r *http.Request) {}
