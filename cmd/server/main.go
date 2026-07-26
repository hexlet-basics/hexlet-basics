package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"hexletbasics/internal/api"
	"hexletbasics/internal/handlers"
)

func main() {
	ctx := context.Background()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@127.0.0.1:54330/code_basics_development"
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("db ping: %v", err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// Dev CORS so the Vite frontend (any localhost port) can call the API.
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:*", "http://127.0.0.1:*"},
	}))

	srv := handlers.NewServer(pool)
	strict := api.NewStrictHandler(srv, nil)
	api.RegisterHandlers(e, strict)

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":3001"
	}
	log.Printf("backend listening on %s", addr)
	if err := e.Start(addr); err != nil {
		log.Fatal(err)
	}
}
