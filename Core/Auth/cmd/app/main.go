package main

import (
	"auth/internal/app"
	"auth/pkg/config"
	"auth/pkg/lib/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("/app/.env")
	if err != nil {
		pwd, _ := os.Getwd()
		panic(pwd)
	}

	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)

	log.Info(
		"starting application", slog.Any("config:", cfg),
	)

	application := app.New(
		log,
		cfg,
	)

	go func() {
		application.GRPCSrv.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCSrv.Stop()
	log.Info("Gracefully stopped")
}
