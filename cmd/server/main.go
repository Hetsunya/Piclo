package main

import (
	"log"

	"Piclo/internal/config"
	"Piclo/internal/handler"

)

func main() {
	cfg := config.Load()
	router := handler.NewRouter(cfg)

	serverAddr := ":" + cfg.Port
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}