package main

import (
	"articlesManageService/internal/app"
	"articlesManageService/internal/storage/real/psql"
	"articlesManageService/pkg/config"
	"articlesManageService/pkg/lib/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	log.Info("application started")

	storage := psql.New(log, os.Getenv("DATABASE_URL"))
	// storage := psql.New("postgres://postgres:123@psql:5432/postgres?sslmode=disable")

	application := app.New(log, cfg.Grpc.Port, storage)

	go func() {
		application.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	application.GRPCServer.Stop()
	storage.Close()
	log.Info("application stopped")

}
