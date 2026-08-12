package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"Piclo/internal/config"
	"Piclo/internal/handler"
	"Piclo/internal/repository"
	"Piclo/internal/service"
	"Piclo/internal/storage"
)

func main() {
	cfg := config.Load()

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	// пингуем БД
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("DB ping failed: %v", err)
	}
	defer pool.Close()

	// Инициализация MinIO
	minioStore, err := storage.NewMinIOStorage(
		cfg.MinIOEndpoint,
		cfg.MinIOAccessKey,
		cfg.MinIOSecretKey,
		cfg.MinIOBucket,
		false, // useSSL = false для локального dev
	)
	if err != nil {
		log.Fatalf("MinIO connection failed: %v", err)
	}

	repo := repository.NewImageRepository(pool)
	imgService := service.NewImageService(repo, minioStore)

	router := handler.NewRouter(cfg, imgService, minioStore)

	serverAddr := ":" + cfg.Port
	log.Printf("Starting server on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
