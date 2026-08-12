package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DatabaseURL    string
	PublicURL      string
	MinIOEndpoint  string
	MinIOAccessKey string
	MinIOSecretKey string
	MinIOBucket    string
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		Port:           getEnv("PORT", "9090"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://piclo:piclo@localhost:5433/piclo?sslmode=disable"),
		PublicURL:      getEnv("PUBLIC_URL", "http://localhost:9090"),
		MinIOEndpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinIOAccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinIOSecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
		MinIOBucket:    getEnv("MINIO_BUCKET", "piclo-images"),
	}
	log.Printf("Configuration loaded")
	return cfg
}

// getEnv оставляем как был

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
