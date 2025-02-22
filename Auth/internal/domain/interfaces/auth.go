package interfaces

import (
	"context"
)

type Auth interface {
	Login(ctx context.Context, email string, password string, appID string) (token string, err error)
	Register(ctx context.Context, email string, password string) (user_id string, err error)
	IsAdmin(ctx context.Context, user_id string) (isAdmin bool, err error)
}
