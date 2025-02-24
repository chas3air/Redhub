package mockapp

import (
	"auth/internal/domain/models"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type AppProvider struct {
	log  *slog.Logger
	apps map[uuid.UUID]models.App
}

func New(log *slog.Logger) *AppProvider {
	apps := make(map[uuid.UUID]models.App, 5)
	generatedID, err := uuid.Parse("395ac0be-afd9-4434-a2d1-65d2472ad009")
	if err != nil {
		panic(err)
	}
	log.Info("Showing app-id:", slog.String("app_id", generatedID.String()))
	apps[generatedID] = models.App{
		Id:     generatedID,
		Alias:  "Postman",
		Secret: "secret",
	}

	return &AppProvider{
		log:  log,
		apps: apps,
	}
}

func (a *AppProvider) App(ctx context.Context, app_id uuid.UUID) (models.App, error) {
	return a.apps[app_id], nil
}

func (a *AppProvider) Insert(ctx context.Context, app models.App) error {
	a.apps[app.Id] = app
	return nil
}

func (a *AppProvider) GetAll(ctx context.Context) ([]models.App, error) {
	apps := make([]models.App, 0, len(a.apps))
	for _, v := range a.apps {
		apps = append(apps, v)
	}

	return apps, nil
}
