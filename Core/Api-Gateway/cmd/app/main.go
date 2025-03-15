package main

import (
	"apigateway/internal/app"
	"apigateway/pkg/config"
	"apigateway/pkg/lib/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	log.Info(
		"starting application", slog.Any("config:", cfg),
	)

	application := app.New(log, cfg)
	go func() {
		application.Start()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	log.Info("Gracefully stopped")
}
