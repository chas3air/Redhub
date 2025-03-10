package usermanager

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"usersManageService/internal/domain/interfaces/storage"
	"usersManageService/internal/domain/models"
	"usersManageService/internal/services"
	storage_errors "usersManageService/internal/storage"
	"usersManageService/pkg/lib/logger/sl"

	"github.com/google/uuid"
)

type UserManager struct {
	log     *slog.Logger
	storage storage.Storage
}

var ErrInvalidCredentials = errors.New("invalid credentials")

func New(log *slog.Logger, storage storage.Storage) *UserManager {
	return &UserManager{
		log:     log,
		storage: storage,
	}
}

func (um *UserManager) GetUsers(ctx context.Context) ([]models.User, error) {
	const op = "services.userManager.GetUsers"
	log := um.log.With(slog.String("operation", op))

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	users, err := um.storage.GetUsers(ctx)
	if err != nil {
		log.Error("Failed to retrieve users", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully retrieved users")
	return users, nil
}

func (um *UserManager) GetUserById(ctx context.Context, uid uuid.UUID) (models.User, error) {
	const op = "services.userManager.GetUserById"
	log := um.log.With(slog.String("operation", op))

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	user, err := um.storage.GetUserById(ctx, uid)
	if err != nil {
		if errors.Is(err, storage_errors.ErrNotFound) {
			log.Error("User not found", sl.Err(err))
			return models.User{}, fmt.Errorf("%s: %w", op, services.ErrNotFound)
		}

		log.Error("Failed to retrieve user by id", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully retrieved user")
	return user, nil
}

func (um *UserManager) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	const op = "services.userManager.GetUserByEmail"
	log := um.log.With(slog.String("operation", op))

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	user, err := um.storage.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, storage_errors.ErrNotFound) {
			log.Error("User not found", sl.Err(err))
			return models.User{}, fmt.Errorf("%s: %w", op, services.ErrNotFound)
		}

		log.Error("Failed to retrieve user by email", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully retrieved user")
	return user, nil
}

func (um *UserManager) Insert(ctx context.Context, user models.User) (models.User, error) {
	const op = "services.userManager.Insert"
	log := um.log.With(slog.String("operation", op))

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	user, err := um.storage.Insert(ctx, user)
	if err != nil {
		if errors.Is(err, storage_errors.ErrAlreadyExists) {
			log.Warn("User already exists", sl.Err(err))
			return models.User{}, fmt.Errorf("%s: %w", op, services.ErrAlreadyExists)
		}

		log.Error("Failed to insert user", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("User inserted successfully")
	return user, nil
}

func (um *UserManager) Update(ctx context.Context, uid uuid.UUID, user models.User) (models.User, error) {
	const op = "services.userManager.Update"
	log := um.log.With(slog.String("operation", op))

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	user, err := um.storage.Update(ctx, uid, user)
	if err != nil {
		if errors.Is(err, storage_errors.ErrNotFound) {
			log.Warn("User not found for update", sl.Err(err))
			return models.User{}, fmt.Errorf("%s: %w", op, services.ErrNotFound)
		}

		log.Error("Failed to update user", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("User updated successfully")
	return user, nil
}

func (um *UserManager) Delete(ctx context.Context, uid uuid.UUID) (models.User, error) {
	const op = "services.userManager.Delete"
	log := um.log.With(slog.String("operation", op))

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	user, err := um.storage.Delete(ctx, uid)
	if err != nil {
		if errors.Is(err, storage_errors.ErrNotFound) {
			log.Warn("User not found for delete", sl.Err(err))
			return models.User{}, fmt.Errorf("%s: %w", op, services.ErrNotFound)
		}

		log.Error("Failed to delete user by id", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("User deleted successfully")
	return user, nil
}
