package main

import (
	"log"
	"os"

	"shorturl.com/internal/handler"
	"shorturl.com/internal/repository"
	"shorturl.com/internal/service"
	"shorturl.com/pkg/config"
	"shorturl.com/pkg/logger"
)

func main() {
	cfg := config.LoadConfig()

	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %s", err)
	}
	defer logFile.Close()

	appLogger := logger.New(
		logger.WithLevel(cfg.LogLevel),
		logger.WithOutput(logFile),
	)
	appLogger.Info("starting URL shortener service")

	repo, err := repository.NewPostgresRepository(cfg.DatabaseURL, appLogger)
	if err != nil {
		appLogger.Error("failed to connect to database", logger.Err(err))
		log.Fatalf("Failed to connect to database: %s", err)
	}
	defer repo.Close()

	service := service.NewShortenerService(repo)
	handler := handler.NewHandler(service, cfg, appLogger)

	router := handler.InitRoutes()

	appLogger.Info("server starting",
		logger.String("address", cfg.ServerAddress),
	)

	// log.Printf("Server starting on %s", cfg.ServerAddress)
	log.Fatal(router.Run(cfg.ServerAddress))
}
