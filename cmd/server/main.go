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

	//База данных
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("DB ping failed: %v", err)
	}
	defer pool.Close()

	//MinIO
	minioStore, err := storage.NewMinIOStorage(
		cfg.MinIOEndpoint, cfg.MinIOAccessKey, cfg.MinIOSecretKey, cfg.MinIOBucket, false,
	)
	if err != nil {
		log.Fatalf("MinIO connection failed: %v", err)
	}

	//Зависимости
	repo := repository.NewImageRepository(pool)
	imgService := service.NewImageService(repo, minioStore)
	router := handler.NewRouter(cfg, imgService, minioStore)

	// Запуск TTL Воркера в отдельной горутине
	// Создаем корневой контекст с возможностью отмены
	ctx, cancel := context.WithCancel(context.Background())

	// Воркер будет проверять базу каждые 1 минуту
	ttlWorker := worker.NewTTLWorker(repo, minioStore, 1*time.Minute)
	go ttlWorker.Run(ctx)

	// Запуск HTTP сервера
	serverAddr := ":" + cfg.Port
	log.Printf("Starting server on %s", serverAddr)

	// Запускаем сервер в горутине, чтобы иметь возможность перехватить сигналы ОС
	go func() {
		if err := router.Run(serverAddr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Ожидание сигнала завершения (Ctrl+C)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Отменяем контекст, чтобы воркер корректно остановился
	cancel()

	// Даем воркеру пару секунд на завершение текущей итерации (опционально, но хорошая практика)
	time.Sleep(2 * time.Second)
	log.Println("Server and workers stopped gracefully")
}
