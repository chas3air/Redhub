package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"usersManageService/internal/app"

	"usersManageService/pkg/config"
	"usersManageService/pkg/lib/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	log.Info("starting application", slog.Any("config:", cfg))

	application := app.New(log, cfg.Grpc.Port)

	go func() {
		application.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	application.GRPCServer.Stop()
	log.Info("application stopped")
}
