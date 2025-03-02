package main

import (
	"articlesManageService/pkg/config"
	"articlesManageService/pkg/lib/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	log.Info("application started")

	

}
