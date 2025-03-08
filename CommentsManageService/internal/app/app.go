package app

import (
	grpcapp "commentsManageService/internal/app/grpc"
	"commentsManageService/internal/domain/interfaces/storage"
	commentservice "commentsManageService/internal/service/commentService"
	"log/slog"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, storage storage.CommentStorage, port int) *App {
	// storage := psqlstorage.New(log, os.Getenv("DATABASE_URL"))
	commentsservice := commentservice.New(log, storage)

	grpcapp := grpcapp.New(log, commentsservice, port)
	return &App{
		GRPCServer: grpcapp,
	}
}
