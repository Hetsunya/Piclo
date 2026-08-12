package config

import (
	"log"
	"os"
)

type Config struct {
	Port string
}

func Load() *Config {
	cfg := &Config{
		Port: getEnv("PORT", "9090"),
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
