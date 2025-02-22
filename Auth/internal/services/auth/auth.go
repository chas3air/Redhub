package authservice

import (
	"auth/internal/domain/interfaces"
	"context"
	"log/slog"
)

type AuthService struct {
	log     *slog.Logger
	storage interfaces.Storage
}

func New(log *slog.Logger, storage interfaces.Storage) *AuthService {
	return &AuthService{
		log:     log,
		storage: storage,
	}
}

// Login implements interfaces.Auth.
func (a AuthService) Login(ctx context.Context, email string, password string, appID string) (token string, err error) {
	panic("unimplemented")
}

// Register implements interfaces.Auth.
func (a AuthService) Register(ctx context.Context, email string, password string) (user_id string, err error) {
	panic("unimplemented")
}

// IsAdmin implements interfaces.Auth.
func (a AuthService) IsAdmin(ctx context.Context, user_id string) (isAdmin bool, err error) {
	panic("unimplemented")
}
