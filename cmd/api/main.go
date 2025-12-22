package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ethaan/miracle74-api/internal/api"
	"github.com/ethaan/miracle74-api/internal/handlers"
	"github.com/ethaan/miracle74-api/internal/repo"
	"github.com/ethaan/miracle74-api/internal/services"
	"github.com/ethaan/miracle74-api/pkg/cache"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	cacheURL := os.Getenv("CACHE_URL")
	if cacheURL == "" {
		cacheURL = "localhost:6379"
	}

	cacheClient, err := cache.NewClient(cacheURL, cache.DefaultTTL)
	if err != nil {
		log.Fatalf("Failed to connect to cache: %v", err)
	}
	defer cacheClient.Close()
	log.Printf("Connected to Valkey cache at %s", cacheURL)

	// Repos
	characterRepo := repo.NewCharacterRepo(cacheClient)
	powerGamersRepo := repo.NewPowerGamersRepo(cacheClient)
	insomniacsRepo := repo.NewInsomniacsRepo(cacheClient)
	guildRepo := repo.NewGuildRepo(cacheClient)
	whoIsOnlineRepo := repo.NewWhoIsOnlineRepo(cacheClient)

	// Services
	characterService := services.NewCharacterService(characterRepo)
	powerGamersService := services.NewPowerGamersService(powerGamersRepo)
	insomniacsService := services.NewInsomniacsService(insomniacsRepo)
	guildService := services.NewGuildService(guildRepo)
	whoIsOnlineService := services.NewWhoIsOnlineService(whoIsOnlineRepo)

	// Handlers
	handler := handlers.NewHandler(characterService, powerGamersService, insomniacsService, guildService, whoIsOnlineService)

	srv, err := api.NewServer(handler)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	httpServer := &http.Server{
		Addr:         ":" + port,
		Handler:      srv,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Starting server on port %s", port)
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
