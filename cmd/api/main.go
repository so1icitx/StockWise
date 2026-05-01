package main

import (
	"context"
	"fmt"
	"log"

	"github.com/so1icitx/StockWise/internal/application"
	"github.com/so1icitx/StockWise/internal/config"
	"github.com/so1icitx/StockWise/internal/infrastructure/postgres"
	"github.com/so1icitx/StockWise/internal/transport/httpapi"
	websocketapi "github.com/so1icitx/StockWise/internal/transport/websocket"
)

// main starts the StockWise HTTP API.
func main() {
	cfg := config.Load()
	db, err := postgres.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}

	store := postgres.NewStore(db)
	notificationHub := websocketapi.NewHub()
	go notificationHub.Run(context.Background())

	services := application.NewServices(store, store, notificationHub)
	router := httpapi.NewRouter(cfg, services, notificationHub)

	addr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	log.Printf("starting StockWise API on %s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
