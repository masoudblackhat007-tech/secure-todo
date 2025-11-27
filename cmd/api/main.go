package main

import (
	"log"
	"net/http"

	"github.com/masoudblackhat007-tech/secure-todo/internal/config"
	"github.com/masoudblackhat007-tech/secure-todo/internal/http/router"
	"github.com/masoudblackhat007-tech/secure-todo/internal/logging"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// initialize logger (دو خروجی: logger و error)
	logger, err := logging.Init(cfg.Server.Env)
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}

	// handler موقت برای /health تا وقتی که خودت handler واقعی بنویسی
	healthHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	r := router.NewRouter(logger, cfg, healthHandler)

	addr := ":" + cfg.Server.Port
	logger.Info("starting secure-todo api", zap.String("addr", addr))

	if err := http.ListenAndServe(addr, r); err != nil {
		logger.Fatal("server exited", zap.Error(err))
	}
}
