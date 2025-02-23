package interfaces

import (
	"context"

	"github.com/google/uuid"
)

type Auth interface {
	Login(ctx context.Context, email string, password string, appID uuid.UUID) (token string, err error)
	Register(ctx context.Context, email string, password string) (user_id uuid.UUID, err error)
	IsAdmin(ctx context.Context, user_id uuid.UUID) (isAdmin bool, err error)
}
