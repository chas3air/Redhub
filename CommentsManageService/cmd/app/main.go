package main

import (
	"commentsManageService/internal/app"
	psqlstorage "commentsManageService/internal/storage/real/psql"
	"commentsManageService/pkg/config"
	"commentsManageService/pkg/lib/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	log.Info("starting application", slog.Any("config", cfg))

	// storage := mock.New(log)
	storage := psqlstorage.New(log, os.Getenv("DATABASE_URL"))

	application := app.New(log, storage, cfg.Grpc.Port)

	go func() {
		application.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	log.Info("Stoping application")
	application.GRPCServer.Stop()

	log.Info("Closing conn to DB")
	storage.DB.Close()

	log.Info("application stopped")
}
