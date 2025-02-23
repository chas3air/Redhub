package app

import (
	grpcapp "auth/internal/app/grpc"
	authservice "auth/internal/services/auth"
	mockapp "auth/internal/storage/mock/app"
	"auth/internal/storage/usersmanageservice"
	"log/slog"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, tokenTTL time.Duration, port int) *App {
	usersStorage := usersmanageservice.New(log, "user_service", 50051)
	//usersStorage := mockusers.New()
	appProvider := mockapp.New(log)

	authservice := authservice.New(log, usersStorage, appProvider, tokenTTL)
	grpcapp := grpcapp.New(log, authservice, port)

	return &App{
		GRPCSrv: grpcapp,
	}
}
