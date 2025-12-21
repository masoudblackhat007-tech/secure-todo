// internal/store/store.go
package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/masoudblackhat007-tech/secure-todo/internal/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type RedisClient struct {
	Client *redis.Client
}

func (rc *RedisClient) Close() error {
	if rc == nil || rc.Client == nil {
		return nil
	}
	return rc.Client.Close()
}

func NewDB(ctx context.Context, cfg *config.DatabaseConfig, logger *zap.Logger) (*sql.DB, error) {
	if cfg == nil {
		return nil, fmt.Errorf("database config is nil")
	}

	// NOTE: driver is "pgx" (registered by jackc/pgx stdlib import above)
	db, err := sql.Open("pgx", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("sql open failed: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("db ping failed: %w", err)
	}

	logger.Info("db initialized")
	return db, nil
}

func NewRedis(ctx context.Context, cfg *config.RedisConfig, logger *zap.Logger) (*RedisClient, error) {
	if cfg == nil {
		return nil, fmt.Errorf("redis config is nil")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if _, err := rdb.Ping(pingCtx).Result(); err != nil {
		_ = rdb.Close()
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	logger.Info("redis initialized")
	return &RedisClient{Client: rdb}, nil
}
