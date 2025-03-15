package app

import (
	"log/slog"
	"os"
	grpcapp "usersManageService/internal/app/grpc"
	usermanager "usersManageService/internal/services/usersManager"
	psqlstorage "usersManageService/internal/storage/real/psql"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, port int) *App {
	storage := psqlstorage.New(log, os.Getenv("DATABASE_URL"))
	//storage := mock.New()
	usermanager := usermanager.New(log, storage)

	grpcapp := grpcapp.New(log, usermanager, port)
	return &App{
		GRPCServer: grpcapp,
	}
}
