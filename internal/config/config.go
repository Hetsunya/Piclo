package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	PublicURL   string
}

func Load() *Config {
	// Эта строка загружает переменные из файла .env в окружение Go
	// Ошибку игнорируем, чтобы код работал и в проде, где .env нет, а есть реальные env vars
	_ = godotenv.Load()

	cfg := &Config{
		Port:        getEnv("PORT", "9090"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://piclo:piclo@localhost:5433/piclo?sslmode=disable"),
		PublicURL:   getEnv("PUBLIC_URL", "http://localhost:9090"),
	}

	log.Printf("Configuration loaded: %+v", cfg)
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
