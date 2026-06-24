package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"openownership-workflow/backend/internal/cache"
	"openownership-workflow/backend/internal/config"
	"openownership-workflow/backend/internal/database"
	"openownership-workflow/backend/internal/handlers"
	"openownership-workflow/backend/internal/repositories"
	"openownership-workflow/backend/internal/services"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel}))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	gormDB, err := database.OpenPostgres(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("connect database", "error", err)
		os.Exit(1)
	}
	sqlDB, err := gormDB.DB()
	if err != nil {
		logger.Error("open database handle", "error", err)
		os.Exit(1)
	}
	defer sqlDB.Close()
	if err := database.Migrate(gormDB); err != nil {
		logger.Error("run migrations", "error", err)
		os.Exit(1)
	}
	if err := database.SeedDemoData(gormDB); err != nil {
		logger.Error("seed demo data", "error", err)
		os.Exit(1)
	}

	redisClient := cache.OpenRedis(cfg.RedisURL)
	defer redisClient.Close()

	repo := repositories.New(gormDB)
	dashboardService := services.NewDashboardService(repo, redisClient)
	authService := services.NewAuthService(repo, cfg.JWTSecret)
	userService := services.NewUserService(repo)
	accessService := services.NewAccessService(repo)
	auditService := services.NewAuditService(repo)
	submissionService := services.NewSubmissionService(repo, dashboardService)

	router := handlers.NewRouter(handlers.Dependencies{
		Config:      cfg,
		Auth:        authService,
		Users:       userService,
		Access:      accessService,
		Audit:       auditService,
		Submissions: submissionService,
		Dashboard:   dashboardService,
		Logger:      logger,
	})

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("api listening", "port", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err)
			stop()
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown failed", "error", err)
	}
}
