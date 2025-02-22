package usersmanageservice

import (
	"auth/internal/domain/models"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type UsersManageService struct {
	log         *slog.Logger
	ServiceHost string
	ServicePort int
}

func New(log *slog.Logger, serviceHost string, servicePort int) *UsersManageService {
	return &UsersManageService{
		log:         log,
		ServiceHost: serviceHost,
		ServicePort: servicePort,
	}
}

// GetUserByEmail implements interfaces.Storage.
func (u UsersManageService) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	panic("unimplemented")
}

// GetUserById implements interfaces.Storage.
func (u UsersManageService) GetUserById(ctx context.Context, uid uuid.UUID) (models.User, error) {
	panic("unimplemented")
}

// GetUsers implements interfaces.Storage.
func (u UsersManageService) GetUsers(ctx context.Context) ([]models.User, error) {
	panic("unimplemented")
}

// Insert implements interfaces.Storage.
func (u UsersManageService) Insert(ctx context.Context, user models.User) error {
	panic("unimplemented")
}

// Update implements interfaces.Storage.
func (u UsersManageService) Update(ctx context.Context, uid uuid.UUID, user models.User) error {
	panic("unimplemented")
}

// Delete implements interfaces.Storage.
func (u UsersManageService) Delete(ctx context.Context, uid uuid.UUID) (models.User, error) {
	panic("unimplemented")
}
