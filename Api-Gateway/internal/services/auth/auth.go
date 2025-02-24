package authservice

import (
	"apigateway/internal/domain/interfaces"
	"apigateway/internal/domain/models"
	"context"
	"fmt"
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

func (as *AuthService) Login(ctx context.Context, email string, password string, appID uuid.UUID) (string, error) {
	const op = "service.auth.login"
	log := as.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)
	_ = log

	select {
	case <-ctx.Done():
		return "", fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	token, err := as.storage.Login(ctx, email, password, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (as *AuthService) Register(ctx context.Context, user models.User) error {
	const op = "service.auth.register"
	log := as.log.With(
		slog.String("op", op),
		slog.String("uid", user.Email),
	)
	_ = log

	select {
	case <-ctx.Done():
		return fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	err := as.storage.Register(ctx, user)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
func (as *AuthService) IsAdmin(ctx context.Context, user_id uuid.UUID) (bool, error) {
	const op = "service.auth.isAdmin"
	log := as.log.With(
		slog.String("op", op),
		slog.String("uid", user_id.String()),
	)
	_ = log

	select {
	case <-ctx.Done():
		return false, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	isAdmin, err := as.storage.IsAdmin(ctx, user_id)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}
