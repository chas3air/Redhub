package app

import (
	grpcapp "auth/internal/app/grpc"
	authservice "auth/internal/services/auth"
	"auth/internal/storage/usersmanageservice"
	"log/slog"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, port int) *App {
	storage := usersmanageservice.New(log, "server", 50051)

	authservice := authservice.New(log, storage)
	grpcsrv := grpcapp.New(log, authservice, port)

	return &App{
		GRPCSrv: grpcsrv,
	}
}
