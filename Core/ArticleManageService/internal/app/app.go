package app

import (
	grpcapp "articlesManageService/internal/app/grpc"
	"articlesManageService/internal/domain/interfaces/storage"
	articlemanager "articlesManageService/internal/services/articleManager"
	"log/slog"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, port int, storage storage.Storage) *App {
	articleManager := articlemanager.New(log, storage)

	grpcapp := grpcapp.New(log, articleManager, port)
	return &App{
		GRPCServer: grpcapp,
	}
}
