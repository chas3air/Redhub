package app

import (
	"log/slog"
	grpcapp "usersManageService/internal/app/grpc"
	usermanager "usersManageService/internal/services/usersManager"
	"usersManageService/internal/storage/mock"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, port int) *App {
	// storage := psqlstorage.New(os.Getenv("DATABASE_URL"))
	storage := mock.New()
	usermanager := usermanager.New(log, storage)

	grpcapp := grpcapp.New(log, usermanager, port)
	return &App{
		GRPCServer: grpcapp,
	}
}
