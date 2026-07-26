package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"
	"github.com/samber/do/v2"

	"hexletbasics/ent"
	"hexletbasics/internal/api"
	"hexletbasics/internal/config"
	"hexletbasics/internal/di"
)

func main() {
	injector := di.New()
	defer func() { _ = do.MustInvoke[*ent.Client](injector).Close() }()

	cfg, err := do.Invoke[*config.Config](injector)
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	apiServer, err := do.Invoke[*api.Server](injector)
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
