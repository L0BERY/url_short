package main

import (
	"log"

	"shorturl.com/internal/handler"
	"shorturl.com/internal/repository"
	"shorturl.com/internal/service"
	"shorturl.com/pkg/config"
)

func main() {
	cfg := config.LoadConfig()

	repo, err := repository.NewPostgresRepository(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err)
	}
	defer repo.Close()

	service := service.NewShortenerService(repo)
	handler := handler.NewHandler(service, cfg)

	router := handler.InitRoutes()

	log.Printf("Server starting on %s", cfg.ServerAddress)
	log.Fatal(router.Run(cfg.ServerAddress))
}
