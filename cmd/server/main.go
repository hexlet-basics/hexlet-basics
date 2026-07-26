package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"

	"hexletbasics/internal/api"
	"hexletbasics/internal/config"
	"hexletbasics/internal/handlers"
	"hexletbasics/internal/store"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	client, err := store.NewClient(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer client.Close()

	srv := handlers.NewServer(client)
	apiServer, err := api.NewServer(srv)
	if err != nil {
		log.Fatalf("api server: %v", err)
	}

	// Dev CORS so the Vite frontend (any localhost port) can call the API.
	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:*", "http://127.0.0.1:*"},
	}).Handler(apiServer)

	log.Printf("backend listening on %s", cfg.Addr)
	if err := http.ListenAndServe(cfg.Addr, handler); err != nil {
		log.Fatal(err)
	}
}
