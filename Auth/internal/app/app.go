package app

import (
	grpcapp "auth/internal/app/grpc"
	authservice "auth/internal/services/auth"
	mockapp "auth/internal/storage/mock/app"
	"auth/internal/storage/usersmanageservice"
	"auth/pkg/config"
	"log/slog"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, cfg *config.Config) *App {
	usersStorage := usersmanageservice.New(log, cfg.UsersStorageHost, cfg.UsersStoragePort)
	//usersStorage := mockusers.New()
	appProvider := mockapp.New(log)

	authservice := authservice.New(log, usersStorage, appProvider, cfg.TokenTTL)
	grpcapp := grpcapp.New(log, authservice, cfg.Grpc.Port)

	return &App{
		GRPCSrv: grpcapp,
	}
}
