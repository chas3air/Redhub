package authservice

import (
	"apigateway/internal/domain/interfaces"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type AuthService struct {
	log     *slog.Logger
	storage interfaces.Auth
}

func New(log *slog.Logger, storage interfaces.Auth) *AuthService {
	return &AuthService{
		log:     log,
		storage: storage,
	}
}

func (as *AuthService) Login(ctx context.Context, email string, password string, appID uuid.UUID) (token string, err error) {
	panic("unimplemented")
}
func (as *AuthService) Register(ctx context.Context, email string, password string) (user_id uuid.UUID, err error) {
	panic("unimplemented")
}
func (as *AuthService) IsAdmin(ctx context.Context, user_id uuid.UUID) (isAdmin bool, err error) {
	panic("unimplemented")
}
