package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/masoudblackhat007-tech/secure-todo/internal/config"

	"go.uber.org/zap"
)

// NewRouter sets up the application's http routes.
func NewRouter(logger *zap.Logger, cfg *config.Config, healthHandler http.Handler) http.Handler {
	r := chi.NewRouter()

	// ---- Global middlewares ----
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	// Custom secure logging (no secrets allowed)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			logger.Info("incoming request",
				zap.String("method", req.Method),
				zap.String("path", req.URL.Path),
			)
			next.ServeHTTP(w, req)
		})
	})

	// ---- Routes ----
	r.Mount("/health", healthHandler)

	return r
}
