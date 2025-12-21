// internal/http/handler/health.go
package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type HealthHandler struct {
	logger *zap.Logger
	db     *sql.DB
	redis  *redis.Client
	start  time.Time
}

func NewHealthHandler(logger *zap.Logger, db *sql.DB, redisClient *redis.Client) *HealthHandler {
	return &HealthHandler{
		logger: logger,
		db:     db,
		redis:  redisClient,
		start:  time.Now(),
	}
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqStart := time.Now()

	// Hard cap for the whole health request
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	checks := map[string]string{
		"db":    "ok",
		"redis": "ok",
	}

	dbOK := true
	redisOK := true

	// DB check with short timeout
	if h.db != nil {
		dbCtx, dbCancel := context.WithTimeout(ctx, 500*time.Millisecond)
		if err := h.db.PingContext(dbCtx); err != nil {
			checks["db"] = "down"
			dbOK = false
		}
		dbCancel()
	} else {
		checks["db"] = "down"
		dbOK = false
	}

	// Redis check with short timeout
	if h.redis != nil {
		rdCtx, rdCancel := context.WithTimeout(ctx, 500*time.Millisecond)
		if _, err := h.redis.Ping(rdCtx).Result(); err != nil {
			checks["redis"] = "down"
			redisOK = false
		}
		rdCancel()
	} else {
		checks["redis"] = "down"
		redisOK = false
	}

	httpStatus := http.StatusOK
	status := "ok"

	if !dbOK || !redisOK {
		httpStatus = http.StatusInternalServerError
		if !dbOK && !redisOK {
			status = "down"
		} else {
			status = "degraded"
		}
	}

	resp := map[string]any{
		"status": status,
		"uptime": time.Since(h.start).String(),
		"checks": checks,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	_ = json.NewEncoder(w).Encode(resp)

	h.logger.Info("health_check",
		zap.Duration("duration", time.Since(reqStart)),
		zap.String("status", status),
		zap.String("db", checks["db"]),
		zap.String("redis", checks["redis"]),
	)
}
