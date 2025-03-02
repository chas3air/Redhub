package interfaces

import (
	"apigateway/internal/domain/models"
	"context"

	"github.com/google/uuid"
)

type Auth interface {
	Login(ctx context.Context, email string, password string) (accessToken string, refreshToken string, err error)
	Register(ctx context.Context, user models.User) (err error)
	IsAdmin(ctx context.Context, user_id uuid.UUID) (isAdmin bool, err error)
}
