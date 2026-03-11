package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yashashkumar/ai-issue-agent/internal/config"
	"github.com/yashashkumar/ai-issue-agent/internal/database"
	"github.com/yashashkumar/ai-issue-agent/internal/handler"
	"github.com/yashashkumar/ai-issue-agent/internal/repository"
	"github.com/yashashkumar/ai-issue-agent/internal/router"
	"github.com/yashashkumar/ai-issue-agent/internal/service"
	"github.com/yashashkumar/ai-issue-agent/internal/spawner"
)

func main() {
	// Initialize Config
	cfg := config.Load()

	// Initialize Logger
	var logHandler slog.Handler
	opts := &slog.HandlerOptions{
		Level: parseLogLevel(cfg.LogLevel),
	}
	if cfg.LogFormat == "json" {
		logHandler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		logHandler = slog.NewTextHandler(os.Stdout, opts)
	}
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize Database
	db, err := database.Connect(cfg.DatabasePath)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Migrate(ctx, logger); err != nil {
		logger.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	// Initialize Repositories
	projectRepo := repository.NewProjectRepo(db)
	agentRepo := repository.NewAgentRepo(db)
	worktreeRepo := repository.NewWorktreeRepo(db)

	// Initialize Spawner Engine
	spawnerEngine := spawner.NewSpawner(agentRepo, worktreeRepo, logger, cfg.MaxConcurrentAgents, cfg.SpawnerChannelBuffer)
	spawnerEngine.Start()

	// Initialize Services
	projectSvc := service.NewProjectService(projectRepo)
	agentSvc := service.NewAgentService(agentRepo)

	// Initialize Handlers
	projectHandler := handler.NewProjectHandler(projectSvc)
	agentHandler := handler.NewAgentHandler(agentSvc, projectSvc, spawnerEngine)
	spawnHandler := handler.NewSpawnHandler(spawnerEngine, projectSvc, cfg.DefaultWorkBaseDir)
	webhookHandler := handler.NewWebhookHandler(projectRepo, spawnerEngine)
	// Build Router
	deps := router.RouterDependencies{
		ProjectHandler: projectHandler,
		AgentHandler:   agentHandler,
		SpawnHandler:   spawnHandler,
		WebhookHandler: webhookHandler,
		Logger:         logger,
	}
	mux := router.NewRouter(deps)

	addr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// Graceful Shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		logger.Info("shutting down server gracefully...")
		cancel()

		// Stop taking new spawns and kill running ones
		spawnerEngine.Stop()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Error("server shutdown error", "error", err)
		}
	}()

	logger.Info("starting server", "addr", addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Error("server stopped gracefully", "err", err)
	}
}

func parseLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
