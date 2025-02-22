package mockapp

import (
	"auth/internal/domain/models"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type AppProvider struct {
	log  *slog.Logger
	apps []models.App
}

func New(log *slog.Logger) *AppProvider {
	apps := make([]models.App, 0, 10)
	apps = append(apps, models.App{
		Id:     uuid.New(),
		Alias:  "Postman",
		Secret: "secret",
	})

	return &AppProvider{
		log:  log,
		apps: apps,
	}
}

// App implements interfaces.AppProvider.
func (a *AppProvider) App(ctx context.Context, app_id uuid.UUID) (models.App, error) {
	panic("unimplemented")
}

func (a *AppProvider) Insert(ctx context.Context, app models.App) error {
	panic("uimplemented")
}

func (a *AppProvider) GetAll(ctx context.Context) ([]models.App, error) {
	return a.apps, nil
}
