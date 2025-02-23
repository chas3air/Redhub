package usersmanagerservice

import (
	"apigateway/internal/domain/interfaces"
	"apigateway/internal/domain/models"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type UsersManager struct {
	log     *slog.Logger
	storage interfaces.UsersStorage
}

func New(log *slog.Logger, storage interfaces.UsersStorage) *UsersManager {
	return &UsersManager{
		log:     log,
		storage: storage,
	}
}

func (um *UsersManager) GetUsers(ctx context.Context) ([]models.User, error) {
	panic("unimplemented")
}
func (um *UsersManager) GetUserById(ctx context.Context, uid uuid.UUID) (models.User, error) {
	panic("unimplemented")
}
func (um *UsersManager) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	panic("unimplemented")
}
func (um *UsersManager) Insert(ctx context.Context, user models.User) error {
	panic("unimplemented")
}
func (um *UsersManager) Update(ctx context.Context, uid uuid.UUID, user models.User) error {
	panic("unimplemented")
}
func (um *UsersManager) Delete(ctx context.Context, uid uuid.UUID) (models.User, error) {
	panic("unimplemented")
}
