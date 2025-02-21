package usermanager

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"usersManageService/internal/domain/interfaces/storage"
	"usersManageService/internal/domain/models"
	storage_errors "usersManageService/internal/storage"
	"usersManageService/pkg/lib/logger/sl"

	"github.com/google/uuid"
)

type UserManager struct {
	log     *slog.Logger
	storage storage.Storage
}

var ErrInvalidCredentials = errors.New("invalid credentials")

// TODO: нужно добавить допонительные ошибки (errors.New) для storage
// для более точного анализирования приходящих ошибок
func New(log *slog.Logger, storage storage.Storage) *UserManager {
	return &UserManager{
		log:     log,
		storage: storage,
	}
}

func (um *UserManager) GetUsers(ctx context.Context) ([]models.User, error) {
	const op = "services.userManager.GetUsers"
	log := um.log.With(slog.String("operation", op))

	users, err := um.storage.GetUsers(ctx)
	if err != nil {
		if errors.Is(err, storage_errors.ErrUserNotFound) {
			log.Warn("User not found", sl.Err(err), slog.String("error", err.Error()))

			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("Failed to retrieve users", sl.Err(err), slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully retrieved users", slog.Any("additional info", users), slog.String("error", "nil"))
	return users, nil
}

func (um *UserManager) GetUserById(ctx context.Context, uid uuid.UUID) (models.User, error) {
	const op = "services.userManager.GetUserById"
	log := um.log.With(slog.String("operation", op))

	user, err := um.storage.GetUserById(ctx, uid)
	if err != nil {
		if errors.Is(err, storage_errors.ErrUserNotFound) {
			log.Warn("User not found", sl.Err(err), slog.String("userId", uid.String()), slog.String("error", err.Error()))

			return models.User{}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("Failed to retrieve user by id", sl.Err(err), slog.String("userId", uid.String()), slog.String("error", err.Error()))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully retrieved user", slog.Any("additional info", user), slog.String("error", "nil"))
	return user, nil
}

func (um *UserManager) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	const op = "services.userManager.GetUserByEmail"
	log := um.log.With(slog.String("operation", op))

	user, err := um.storage.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, storage_errors.ErrUserNotFound) {
			log.Warn("User not found", sl.Err(err), slog.String("email", email), slog.String("error", err.Error()))

			return models.User{}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("Failed to retrieve user by email", sl.Err(err), slog.String("email", email), slog.String("error", err.Error()))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully retrieved user", slog.Any("additional info", user), slog.String("error", "nil"))
	return user, nil
}

func (um *UserManager) Insert(ctx context.Context, user models.User) error {
	const op = "services.userManager.Insert"
	log := um.log.With(slog.String("operation", op))

	err := um.storage.Insert(ctx, user)
	if err != nil {
		if errors.Is(err, storage_errors.ErrUserExists) {
			log.Warn("User already exists", sl.Err(err), slog.Any("additional info", user), slog.String("error", err.Error()))

			return fmt.Errorf("%s: %s", op, "user already exists")
		}

		log.Error("Failed to insert user", sl.Err(err), slog.Any("additional info", user), slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("User inserted successfully", slog.Any("additional info", []map[string]interface{}{
		{"user": user},
	}), slog.String("error", "nil"))
	return nil
}

func (um *UserManager) Update(ctx context.Context, uid uuid.UUID, user models.User) error {
	const op = "services.userManager.Update"
	log := um.log.With(slog.String("operation", op))

	err := um.storage.Update(ctx, uid, user)
	if err != nil {
		if errors.Is(err, storage_errors.ErrUserNotFound) {
			log.Warn("User not found", sl.Err(err), slog.String("userId", uid.String()), slog.String("error", err.Error()))

			return fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("Failed to update user", sl.Err(err), slog.String("userId", uid.String()), slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("User updated successfully", slog.Any("additional info", []map[string]interface{}{
		{"user": user},
	}), slog.String("error", "nil"))
	return nil
}

func (um *UserManager) Delete(ctx context.Context, uid uuid.UUID) (models.User, error) {
	const op = "services.userManager.Delete"
	log := um.log.With(slog.String("operation", op))

	user, err := um.storage.Delete(ctx, uid)
	if err != nil {
		if errors.Is(err, storage_errors.ErrUserNotFound) {
			log.Warn("User not found", sl.Err(err), slog.String("userId", uid.String()), slog.String("error", err.Error()))

			return models.User{}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("Failed to delete user by id", sl.Err(err), slog.String("userId", uid.String()), slog.String("error", err.Error()))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("User deleted successfully", slog.Any("additional info", []map[string]interface{}{
		{"user": user},
	}), slog.String("error", "nil"))
	return user, nil
}
