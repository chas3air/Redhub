package main

import (
	"auth/pkg/config"
	"auth/pkg/lib/logger"
)

func main() {
	cfg := config.MustLoad()

	logger := logger.SetupLogger(cfg.Env)

	_ = logger
}
