package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"Piclo/internal/config"
	"Piclo/internal/handler"
	"Piclo/internal/repository"
	"Piclo/internal/service"
	"Piclo/internal/storage"
	"Piclo/internal/worker"

)

func main() {
	cfg := config.Load()

	// Database
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("DB ping failed: %v", err)
	}
	defer pool.Close()

	// MinIO
	minioStore, err := storage.NewMinIOStorage(
		cfg.MinIOEndpoint, cfg.MinIOAccessKey, cfg.MinIOSecretKey, cfg.MinIOBucket, false,
	)
	if err != nil {
		log.Fatalf("MinIO connection failed: %v", err)
	}

	// Dependencies
	repo := repository.NewImageRepository(pool)
	imgService := service.NewImageService(repo, minioStore)
	router := handler.NewRouter(cfg, imgService, minioStore)

	// Start TTL Worker in a separate goroutine
	// Create root context with cancellation capability
	ctx, cancel := context.WithCancel(context.Background())

	// Worker checks database every 1 minute
	ttlWorker := worker.NewTTLWorker(repo, minioStore, 1*time.Minute)
	go ttlWorker.Run(ctx)

	// Start HTTP server
	serverAddr := ":" + cfg.Port
	log.Printf("Starting server on %s", serverAddr)

	// Run server in goroutine to handle OS signals
	go func() {
		if err := router.Run(serverAddr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for shutdown signal (Ctrl+C)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Cancel context to stop worker gracefully
	cancel()

	// Give worker a couple of seconds to finish current iteration (optional but good practice)
	time.Sleep(2 * time.Second)
	log.Println("Server and workers stopped gracefully")
}

