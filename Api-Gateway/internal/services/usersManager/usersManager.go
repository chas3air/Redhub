package usersmanagerservice

import (
	"apigateway/internal/domain/interfaces/users"
	"apigateway/internal/domain/models"
	"apigateway/pkg/lib/logger/sl"
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

type UsersManager struct {
	log     *slog.Logger
	storage users.UsersStorage
}

func New(log *slog.Logger, storage users.UsersStorage) *UsersManager {
	return &UsersManager{
		log:     log,
		storage: storage,
	}
}

func (um *UsersManager) GetUsers(ctx context.Context) ([]models.User, error) {
	const op = "services.usersManager.getUsers"
	log := um.log.With(
		slog.String("op", op),
	)

	users, err := um.storage.GetUsers(ctx)
	if err != nil {
		log.Error("error retrieving users", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return users, nil
}

func (um *UsersManager) GetUserById(ctx context.Context, uid uuid.UUID) (models.User, error) {
	const op = "services.usersManager.getUserById"
	log := um.log.With(
		slog.String("op", op),
	)

	user, err := um.storage.GetUserById(ctx, uid)
	if err != nil {
		log.Error("error retrieving user by id", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (um *UsersManager) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	const op = "services.usersManager.getUserByEmail"
	log := um.log.With(
		slog.String("op", op),
	)

	user, err := um.storage.GetUserByEmail(ctx, email)
	if err != nil {
		log.Error("error retrieving user by email", sl.Err(err))
		return models.User{}, err
	}

	return user, nil
}

func (um *UsersManager) Insert(ctx context.Context, user models.User) (models.User, error) {
	const op = "services.usersManager.insert"
	log := um.log.With(
		slog.String("op", op),
	)

	user, err := um.storage.Insert(ctx, user)
	if err != nil {
		log.Error("error inserting user", sl.Err(err))
		return models.User{}, err
	}

	return user, nil
}

func (um *UsersManager) Update(ctx context.Context, uid uuid.UUID, user models.User) (models.User, error) {
	const op = "services.usersManager.update"
	log := um.log.With(
		slog.String("op", op),
	)

	user, err := um.storage.Update(ctx, uid, user)
	if err != nil {
		log.Error("error updating user", sl.Err(err))
		return models.User{}, err
	}

	return user, nil
}

func (um *UsersManager) Delete(ctx context.Context, uid uuid.UUID) (models.User, error) {
	const op = "services.usersManager.delete"
	log := um.log.With(
		slog.String("op", op),
	)

	user, err := um.storage.Delete(ctx, uid)
	if err != nil {
		log.Error("error deleting user by id", sl.Err(err))
		return models.User{}, err
	}

	return user, nil
}
