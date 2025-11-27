// internal/http/router/router.go
package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/masoudblackhat007-tech/secure-todo/internal/config"
	"go.uber.org/zap"
)

// responseRecorder wraps http.ResponseWriter to capture the status code.
type responseRecorder struct {
	http.ResponseWriter
	status int
}

func (rw *responseRecorder) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// NewRouter sets up the application's HTTP routes and middlewares.
func NewRouter(logger *zap.Logger, cfg *config.Config, healthHandler http.Handler) http.Handler {
	r := chi.NewRouter()

	// ---- Base middlewares ----
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	// Secure headers
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			// X-XSS-Protection故意 تنظیم نمی‌شود (deprecated در مرورگرهای جدید)
			next.ServeHTTP(w, req)
		})
	})

	// No-cache for dynamic responses
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Cache-Control", "no-store")
			next.ServeHTTP(w, req)
		})
	})

	// Minimal CORS for dev (no wildcards in production)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if cfg.Server.Env != "production" {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			}

			if req.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, req)
		})
	})

	// Request logging (بدون body / header)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			start := time.Now()

			rr := &responseRecorder{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			next.ServeHTTP(rr, req)

			duration := time.Since(start)

			logger.Info("http_request",
				zap.String("method", req.Method),
				zap.String("path", req.URL.Path),
				zap.Int("status", rr.status),
				zap.Duration("duration", duration),
			)
		})
	})

	// Recover from panics and return 500 instead of crashing the process.
	r.Use(middleware.Recoverer)

	// ---- Routes ----
	r.Get("/healthz", healthHandler.ServeHTTP)

	return r
}
