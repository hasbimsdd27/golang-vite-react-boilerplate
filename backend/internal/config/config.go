package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	Port       string
	JWTSecret  string
	StaticDir  string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "inventory"),
		DBPort:     getEnv("DB_PORT", "5432"),
		Port:       getEnv("PORT", "8080"),
		JWTSecret:  getEnv("JWT_SECRET", "default-secret-change-me"),
		StaticDir:  getEnv("STATIC_DIR", "../frontend/dist"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
