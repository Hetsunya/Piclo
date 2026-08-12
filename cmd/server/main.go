package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"Piclo/internal/config"
	"Piclo/internal/handler"
	"Piclo/internal/repository"
	"Piclo/internal/service"
)

func main() {
	cfg := config.Load()

	// 1. Подключение к БД
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("DB pool creation failed: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("DB ping failed: %v", err)
	}

	// 2. Инициализация слоев
	repo := repository.NewImageRepository(pool)
	imgService := service.NewImageService(repo)

	// 3. Запуск сервера
	router := handler.NewRouter(cfg, imgService)

	serverAddr := ":" + cfg.Port
	log.Printf("Starting server on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
