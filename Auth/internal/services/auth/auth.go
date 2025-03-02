package authservice

import (
	"auth/internal/domain/interfaces"
	"auth/internal/domain/models"
	"auth/internal/lib/jwt"
	"auth/internal/storage"
	"auth/pkg/lib/logger/sl"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type AuthService struct {
	log             *slog.Logger
	usersstorage    interfaces.UsersStorage
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func New(log *slog.Logger, usersStorage interfaces.UsersStorage, accessTokenTTL time.Duration, refreshTokenTTL time.Duration) *AuthService {
	return &AuthService{
		log:             log,
		usersstorage:    usersStorage,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID       = errors.New("invalid app id")
	ErrUserExists         = errors.New("user already exists")
)

// Login implements interfaces.Auth.
func (a AuthService) Login(ctx context.Context, email string, password string) (string, string, error) {
	const op = "service.auth.login"
	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("attempt to login")

	select {
	case <-ctx.Done():
		return "", "", fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	user, err := a.usersstorage.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found")
			return "", "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
	}
	log.Info("fetched user:", slog.Any("user", user))

	if user.Password != password {
		log.Warn("user not found")
		return "", "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	accessToken, refreshToken, err := jwt.NewTokens(user, a.accessTokenTTL, a.refreshTokenTTL)
	if err != nil {
		log.Error("failed to generate token", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	return accessToken, refreshToken, nil
}

// Register implements interfaces.Auth.
func (a AuthService) Register(ctx context.Context, user models.User) (err error) {
	const op = "service.auth.register"
	log := a.log.With(
		slog.String("op", op),
		slog.String("email", user.Email),
	)

	select {
	case <-ctx.Done():
		return fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	err = a.usersstorage.Insert(ctx, user)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", sl.Err(err))
			return fmt.Errorf("%s: %w", op, err)
		}

		log.Error("failed to save user", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// IsAdmin implements interfaces.Auth.
func (a AuthService) IsAdmin(ctx context.Context, user_id uuid.UUID) (isAdmin bool, err error) {
	const op = "service.auth.isAdmin"
	log := a.log.With(
		slog.String("op", op),
		slog.String("uid", user_id.String()),
	)

	select {
	case <-ctx.Done():
		return false, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	user, err := a.usersstorage.GetUserById(ctx, user_id)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found")
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		log.Error("failed get user by id", sl.Err(err))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if user.Role == "admin" {
		return true, nil
	} else {
		return false, nil
	}
}
