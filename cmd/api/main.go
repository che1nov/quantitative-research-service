package main

import (
	"os"

	"github.com/che1nov/quantitative-research-service/config"
	"github.com/che1nov/quantitative-research-service/internal/app"
	"github.com/che1nov/quantitative-research-service/pkg/logger"
)

func main() {
	cfg := config.Load()
	appLogger := logger.New(cfg.LogLevel)
	application, err := app.New(cfg)
	if err != nil {
		appLogger.Error("Ошибка инициализации приложения", "error", err)
		os.Exit(1)
	}
	if err := application.Start(); err != nil {
		appLogger.Error("Ошибка запуска приложения", "error", err)
		os.Exit(1)
	}
}
