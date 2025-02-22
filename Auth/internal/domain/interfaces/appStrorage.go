package interfaces

import (
	"auth/internal/domain/models"
	"context"

	"github.com/google/uuid"
)

type AppProvider interface {
	GetAll(ctx context.Context) ([]models.App, error)
	App(ctx context.Context, app_id uuid.UUID) (models.App, error)
	Insert(ctx context.Context, app models.App) error
}
