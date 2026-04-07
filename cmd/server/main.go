package main

import (
	"log"
	"net/http"

	"github.com/jesusthecreator017/PicoURL/cmd/server/api"
	"github.com/jesusthecreator017/PicoURL/internal/config"
	"github.com/jesusthecreator017/PicoURL/internal/service"
	"github.com/jesusthecreator017/PicoURL/internal/store"
	"github.com/joho/godotenv"
)

func main() {
	// Load ENV variables
	_ = godotenv.Load()

	// Load configs
	cfg := config.LoadConfig()

	// Initialize Postgres
	pg, err := store.NewPostgresStore(&cfg.Postgres)
	if err != nil {
		log.Fatal("Failed to connect to postgres:", err)
	}

	// Initialize Redis
	cache, err := store.NewRedisCache(&cfg.Redis, cfg.Redis.CacheTTL)
	if err != nil {
		log.Fatal("Failed to connect to redis:", err)
	}

	// Combine into the cache store
	cachedStore := store.NewCachedStore(pg, cache)
	defer cachedStore.Close()

	// Create a new service
	svc := service.NewURLService(cachedStore)

	// Create the api end and start the server
	app := api.NewApplication(svc)

	log.Printf("Server starting on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, app.Handler()); err != nil {
		log.Fatal("Server Failed:", err)
	}
}
