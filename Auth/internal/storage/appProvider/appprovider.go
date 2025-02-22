package appprovider

import (
	"auth/internal/domain/models"
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type AppProvider struct {
	db *sql.DB
}

func New() *AppProvider {
	return &AppProvider{
		db: nil,
	}
}

// App implements interfaces.AppProvider.
func (a *AppProvider) App(ctx context.Context, app_id uuid.UUID) (models.App, error) {
	panic("unimplemented")
}

func (a *AppProvider) Insert(ctx context.Context, app models.App) error {
	panic("uimplemented")
}
