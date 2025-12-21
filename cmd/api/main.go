// cmd/api/main.go
package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/masoudblackhat007-tech/secure-todo/internal/config"
	"github.com/masoudblackhat007-tech/secure-todo/internal/http/handler"
	"github.com/masoudblackhat007-tech/secure-todo/internal/http/router"
	"github.com/masoudblackhat007-tech/secure-todo/internal/logging"
	"github.com/masoudblackhat007-tech/secure-todo/internal/store"
	"go.uber.org/zap"
)

func main() {
	// Root context for the process lifecycle
	rootCtx, rootCancel := context.WithCancel(context.Background())
	defer rootCancel()

	// OS signal handling for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Load config (no zap yet; fall back to stdlib log on failure)
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Init logger
	logger, err := logging.Init(cfg.Server.Env)
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	logger.Info("starting secure-todo api",
		zap.String("env", cfg.Server.Env),
		zap.String("port", cfg.Server.Port),
	)

	// Init DB + Redis
	var (
		db          *sql.DB
		redisClient = (*store.RedisClient)(nil)
	)

	db, err = store.NewDB(rootCtx, &cfg.Database, logger)
	if err != nil {
		logger.Fatal("failed to init db", zap.Error(err))
	}
	defer func() {
		_ = db.Close()
	}()

	redisClient, err = store.NewRedis(rootCtx, &cfg.Redis, logger)
	if err != nil {
		logger.Fatal("failed to init redis", zap.Error(err))
	}
	defer func() {
		_ = redisClient.Close()
	}()

	// Handlers + router
	healthHandler := handler.NewHealthHandler(logger, db, redisClient.Client)

	r := router.NewRouter(logger, cfg, healthHandler)

	// HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  time.Minute,
	}

	// Start server
	go func() {
		logger.Info("http server listening", zap.String("addr", srv.Addr))

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("http server failed", zap.Error(err))
		}
	}()

	// Wait for termination signal
	<-quit
	logger.Info("shutdown signal received")

	// Stop root context dependent routines (if any)
	rootCancel()

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", zap.Error(err))
	} else {
		logger.Info("graceful shutdown complete")
	}
}
