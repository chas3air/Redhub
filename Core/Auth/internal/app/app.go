package app

import (
	grpcapp "auth/internal/app/grpc"
	authservice "auth/internal/services/auth"
	"auth/internal/storage/real/usersmanageservice"
	"auth/pkg/config"
	"log/slog"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, cfg *config.Config) *App {
	usersStorage := usersmanageservice.New(log, cfg.UsersStorageHost, cfg.UsersStoragePort)

	authservice := authservice.New(log, usersStorage, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)
	grpcapp := grpcapp.New(log, authservice, cfg.Grpc.Port)

	return &App{
		GRPCSrv: grpcapp,
	}
}
