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
	Seed        SeedConfig
}

type SeedConfig struct {
	Requester SeedAccount
	Reviewer  SeedAccount
	Admin     SeedAccount
}

type SeedAccount struct {
	Name     string
	Email    string
	Password string
}

func Load() Config {
	return Config{
		Port:        env("PORT", "8080"),
		DatabaseURL: env("DATABASE_URL", "postgres://workflow:workflow@localhost:5432/workflow?sslmode=disable"),
		RedisURL:    env("REDIS_URL", "redis://localhost:6379/0"),
		JWTSecret:   env("JWT_SECRET", "dev-only-change-me"),
		CORSOrigin:  env("CORS_ORIGIN", "http://localhost:5173"),
		LogLevel:    slog.LevelInfo,
		Seed: SeedConfig{
			Requester: SeedAccount{
				Name:     env("SEED_REQUESTER_NAME", "Amina Requester"),
				Email:    env("SEED_REQUESTER_EMAIL", "requester@example.com"),
				Password: env("SEED_REQUESTER_PASSWORD", "password"),
			},
			Reviewer: SeedAccount{
				Name:     env("SEED_REVIEWER_NAME", "Noah Reviewer"),
				Email:    env("SEED_REVIEWER_EMAIL", "reviewer@example.com"),
				Password: env("SEED_REVIEWER_PASSWORD", "password"),
			},
			Admin: SeedAccount{
				Name:     env("SEED_ADMIN_NAME", "Sam Admin"),
				Email:    env("SEED_ADMIN_EMAIL", "admin@example.com"),
				Password: env("SEED_ADMIN_PASSWORD", "password"),
			},
		},
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
