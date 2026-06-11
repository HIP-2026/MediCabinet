package main

import (
	"context"
	"log"

	"github.com/HIP-2026/MediCabinet/server/internal/api"
	"github.com/HIP-2026/MediCabinet/server/internal/config"
	"github.com/HIP-2026/MediCabinet/server/internal/db"
	"github.com/HIP-2026/MediCabinet/server/internal/inventory"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Printf("warning: could not create database pool: %v", err)
	}
	if pool != nil {
		defer pool.Close()
	}

	var invStore *inventory.Store
	if pool != nil {
		invStore = inventory.New(db.New(pool))
	}

	router := api.NewRouter(pool, invStore)

	log.Printf("MediCabinet server listening on :%s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
