package interfaces

import (
	"apigateway/internal/domain/models"
	"context"

	"github.com/google/uuid"
)

type UsersService interface {
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserById(ctx context.Context, uid uuid.UUID) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	Insert(ctx context.Context, user models.User) error
	Update(ctx context.Context, uid uuid.UUID, user models.User) error
	Delete(ctx context.Context, uid uuid.UUID) (models.User, error)
}
