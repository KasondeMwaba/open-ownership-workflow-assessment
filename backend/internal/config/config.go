package config

import (
	"log/slog"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	RedisURL    string
	JWTSecret   string
	CORSOrigin  string
	LogLevel    slog.Level
}

func Load() Config {
	return Config{
		Port:        env("PORT", "8080"),
		DatabaseURL: env("DATABASE_URL", "postgres://workflow:workflow@localhost:5432/workflow?sslmode=disable"),
		RedisURL:    env("REDIS_URL", "redis://localhost:6379/0"),
		JWTSecret:   env("JWT_SECRET", "dev-only-change-me"),
		CORSOrigin:  env("CORS_ORIGIN", "http://localhost:5173"),
		LogLevel:    slog.LevelInfo,
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
